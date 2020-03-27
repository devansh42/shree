package shree

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

//This file implements basic server for backend

func initServer() {
	listener, err := net.Listen("tcp", joinHost("", 8000)) //Listen @ port 2200
	if err != nil {
		log.Fatal("Failed to Listen")
	}
	log.Print("Starting server at port ", 8000)
	serverconfig := new(ssh.ServerConfig)
	serverconfig.AddHostKey(getHostKey())
	//serverconfig.NoClientAuth = true
	certc := new(ssh.CertChecker)
	certc.IsUserAuthority = userAuthenticator
	serverconfig.PasswordCallback = func(conm ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		return new(ssh.Permissions), nil
	}
	// serverconfig.PublicKeyCallback = func(connm ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	// 	// 	p, err := certc.Authenticate(connm, key) //Authenticates Certificates
	// 	// 	if err != nil {

	// 	// 		log.Print("auth erro", err)
	// 	// 		return nil, err
	// 	// 	}

	// 	return p, err
	// }
	for {
		inconn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection")
			continue
		}
		sconn, newch, newrq, err := ssh.NewServerConn(inconn, serverconfig)
		if err != nil {
			//handle server error
			log.Print(err)
		}
		go handleServerConn(sconn, newch, newrq)

	}
}

func globalRequestHandler(newrq <-chan *ssh.Request, servconn *ssh.ServerConn) {
	rejectReply := func(ch *ssh.Request) {
		if ch.WantReply {
			ch.Reply(false, nil)
		}
	}
	for ch := range newrq {
		switch ch.Type {
		case "tcpip-forward":
			var p struct {
				Address string
				Port    uint32
			}

			err := ssh.Unmarshal(ch.Payload, &p)
			if err != nil {
				rejectReply(ch)
				log.Print(err)
				continue
			}
			//This doesn't covers the case when client forwards '0' as port no
			listener, err := net.Listen("tcp", joinHost(p.Address, int(p.Port)))
			if err != nil {
				rejectReply(ch)
				log.Println("Couldn't start server on given port", err)
				continue
			}

			if ch.WantReply {
				var xp struct {
					Port uint32
				}
				xp.Port = p.Port
				b := ssh.Marshal(&xp) //Replying ok
				ch.Reply(true, b)
			}
			for {
				inconn, err := listener.Accept()
				if err != nil {
					log.Print(err)
					continue //Couldn't continue
				}
				type ppt struct {
					Caddr string
					Cport uint32
					Oaddr string
					Oport uint32
				}

				raddr := inconn.RemoteAddr().String()
				host, sport, _ := net.SplitHostPort(raddr)
				port, _ := strconv.Atoi(sport)
				pp := ppt{p.Address, p.Port, host, uint32(port)}
				b := ssh.Marshal(&pp)
				go func(b []byte, inconn net.Conn, servconn *ssh.ServerConn) {

					sch, rch, err := servconn.OpenChannel("forwarded-tcpip", b)
					if err != nil {

						//handle error
						log.Print("couldn't open channel ", err.Error())
					}
					go ssh.DiscardRequests(rch)

					handleServerConnIO(sch, inconn)

				}(b, inconn, servconn)
			}

		default: //Rejecting all other requests
			if ch.WantReply {
				ch.Reply(false, nil)

			}
		}
	}
}
func handleServerConnIO(sch ssh.Channel, inconn net.Conn) {
	var x = make(chan bool)
	defer inconn.Close()
	defer sch.Close()
	go func() {

		_, err := io.Copy(sch, inconn)
		handleErr(err)
		x <- true
	}()
	go func() {
		_, err := io.Copy(inconn, sch)
		handleErr(err)
		x <- true
	}()
	<-x
}

func handleServerConn(conn *ssh.ServerConn, newch <-chan ssh.NewChannel, newrq <-chan *ssh.Request) {

	go globalRequestHandler(newrq, conn) //Discarding all the out band requests

	for ch := range newch {
		log.Print(ch.ChannelType())
		switch ch.ChannelType() {
		//	ch.
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
			//sch.
			// case "tcpip-forward":
			// 	//This enables remote port forwarding
			// 	//	ch.

		}

	}

}
func userAuthenticator(auth ssh.PublicKey) bool {
	//IsUserAuthority
	hm, _ := os.UserHomeDir()
	f, err := ioutil.ReadFile(hm + "/.ssh/ca_user_key.pub")
	handleErr(err)
	pub, err := ssh.ParsePublicKey(f)
	handleErr(err)
	pm, am := pub.Marshal(), auth.Marshal()
	for i := 0; i < len(pm); i++ {
		if pm[i] != am[i] {
			return false
		}
	}
	return true

}

func getHostKey() ssh.Signer {
	hm, _ := os.UserHomeDir()
	f, err := ioutil.ReadFile(hm + "/.ssh/id_host") //host private key
	handleErr(err)
	k, err := ssh.ParsePrivateKey(f)
	handleErr(err)
	return k
}
