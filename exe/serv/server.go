package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/devansh42/shree/remote"

	"github.com/devansh42/shree/exe"

	"golang.org/x/crypto/ssh"
)

type ppt struct {
	Caddr string
	Cport uint32
	Oaddr string
	Oport uint32
}

var sshListener net.Listener

func publicCallBackFunc(certc *ssh.CertChecker) func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	return func(connm ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
		p, err := certc.Authenticate(connm, key) //Authenticates Certificates
		if err != nil {
			log.Print("Error occured while Authenticating user\t", err)
			return nil, err
		}
		p.Extensions["permit-port-forwarding"] = "" //Permiting only portforwarding
		return p, err

	}
}

func initServer() {
	port := os.Getenv(SHREE_SSH_PORT)
	sshListener, err := net.Listen("tcp", exe.JoinHost("", port)) //Listen @ port 2200
	if err != nil {
		log.Fatal("Failed to Listen")
	}
	log.Printf("Starting server at port %d .....", port)

	serverConfig := new(ssh.ServerConfig)
	serverConfig.AddHostKey(getHostKey())
	certc := new(ssh.CertChecker)
	certc.IsUserAuthority = userAuthenticator
	serverConfig.PublicKeyCallback = publicCallBackFunc(certc)
	for {
		inconn, err := sshListener.Accept()
		if err != nil {
			log.Println("Failed to accept connection")
			continue
		}
		sconn, newch, newrq, err := ssh.NewServerConn(inconn, serverConfig)
		if err != nil {
			//handle server error
			log.Print(err)
		}
		go handleNewServerConn(sconn, newch, newrq)

	}
}

func rejectReply(ch *ssh.Request) {
	if ch.WantReply {
		ch.Reply(false, nil)
	}

}

func handleTCPFwdRequest(ch *ssh.Request, servconn *ssh.ServerConn) {
	var p struct {
		Address string
		Port    uint32
	}

	err := ssh.Unmarshal(ch.Payload, &p)
	if err != nil {
		rejectReply(ch)
		log.Print(err)
		return
	}
	listener, err := net.Listen("tcp", exe.JoinHost(p.Address, int(p.Port)))
	if err != nil {
		rejectReply(ch)
		log.Println("Couldn't start server on given port", err)
		return
	}

	if ch.WantReply {
		var xp struct {
			Port uint32 //Port on which connection is listening at remote side
		}

		_, pp, _ := net.SplitHostPort(listener.Addr().String())
		pi, _ := strconv.Atoi(pp)
		xp.Port = uint32(pi)
		b := ssh.Marshal(&xp)
		ch.Reply(true, b) //Replying ok
	}
	for {
		inconn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue //Couldn't continue
		}

		raddr := inconn.RemoteAddr().String()
		host, sport, _ := net.SplitHostPort(raddr)
		port, _ := strconv.Atoi(sport)
		pp := ppt{p.Address, p.Port, host, uint32(port)}
		b := ssh.Marshal(&pp)

		sch, rch, err := servconn.OpenChannel("forwarded-tcpip", b)
		if err != nil {
			//handle error
			log.Print("couldn't open channel ", err.Error())
		}
		go ssh.DiscardRequests(rch)
		go exe.HandleConnectionIO(inconn, sch)

	}

}

func globalRequestHandler(newrq <-chan *ssh.Request, servconn *ssh.ServerConn) {
	for ch := range newrq {
		switch ch.Type {
		case "tcpip-forward":
			go handleTCPFwdRequest(ch, servconn)
		default: //Rejecting all other requests
			if ch.WantReply {
				ch.Reply(false, nil)

			}
		}
	}
}

func handleNewServerConn(conn *ssh.ServerConn, newch <-chan ssh.NewChannel, newrq <-chan *ssh.Request) {

	go globalRequestHandler(newrq, conn)
	for ch := range newch {
		switch ch.ChannelType() {
		default: //Default behvaiour will to discard the request
			err := ch.Reject(ssh.Prohibited, "channel-type not supported")
			if err != nil {
				log.Println("Couldn't close the channel : ", err.Error())
			}

		case "session":
			_, req, err := ch.Accept()
			ssh.DiscardRequests(req)
			if err != nil {
				log.Print(err)
			}
		}

	}

}
func userAuthenticator(auth ssh.PublicKey) bool {
	certc := new(ssh.Certificate)
	getBackendClient().Call("Backend.GetCAUserPublicCertificate", new(remote.CertificateRequest), certc)

	return bytes.Equal(auth.Marshal(), certc.SignatureKey.Marshal())

}

func getHostKey() ssh.Signer {
	fname := os.Getenv(SHREE_SSH_PRIVATE_KEY)
	f, err := ioutil.ReadFile(fname) //host private key

	k, err := ssh.ParsePrivateKey(f)
	if err != nil {
		log.Fatal("Couldn't parse host private key may be it is broken ", err.Error())
	}
	return k
}
