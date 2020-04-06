package main

//This file contains code for remote port forwarding

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/devansh42/shree/exe"

	"golang.org/x/crypto/ssh"
)

const (
	SSH_PORT           = "SSH_PORT" //For Environment variables
	SSH_HOST           = "SSH_HOST" //For Environment variables
	keysshclientsocket = "sshclient"
)

var (
	//Connection for long lived tcp connections
	sshClientConnection *ssh.Client
	//remoteForwardedPorts lists remote forwarded port
	remoteForwardedPorts []*exe.Forwardedport
)

//getClientSigner returns signed certificate for authentication
func getClientSigner(bpass []byte) (ssh.Signer, error) {
	if currentUser == nil {
		//No logined user
		return nil, errors.New("It seems you are no longer been authenticted. Please authenticate yourself")
	}
	pass := bpass

	havepub, havepr, havecert, pki := searchForPKICredentials(currentUser.Uid)
	if havepub && havepr && havecert {
		///don't know what to do
	} else {
		return nil, errors.New("Invalid or Broken Credentials found, re-authenticate yourself.")
	}

	cert, _, _, _, err := ssh.ParseAuthorizedKey(pki.cert)
	if err != nil {
		return nil, errors.New("Broken Certificate found, please re-authenticate")
		//handle
	}
	prv, err := ssh.ParsePrivateKeyWithPassphrase(pki.prv, pass)
	if err != nil {
		return nil, errors.New("Error occured while processing Private Key : " + err.Error())
	}
	signer, err := ssh.NewCertSigner(cert.(*ssh.Certificate), prv)
	if err != nil {
		return nil, errors.New("Couldn't sign certificate with private key, need re-authentication " + err.Error())
	}
	return signer, nil
}

func getHostCallBack() ssh.HostKeyCallback {
	certb := getServerCertificate()
	certchecker := new(ssh.CertChecker)
	certchecker.IsHostAuthority = func(auth ssh.PublicKey, address string) bool {
		b := marshalauthkey(certb)
		c := marshalauthkey(auth)
		return bytes.Equal(b, c)
	}
	return certchecker.CheckHostKey
}

//forwardRemotePort, forwards remote port src->dest
//it binds dest port on localhost with src port on remote machine
func forwardRemotePort(protocol string, dest int, bpass []byte) string {
	signer, err := getClientSigner(bpass)
	if err != nil {
		print(COLOR_RED)
		println("Couldn't establish remote tunnel:\n", err.Error())
		resetConsoleColor()
		return ""
	}
	if sshClientConnection == nil {
		config := &ssh.ClientConfig{
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: getHostCallBack(),
			User:            currentUser.Username}

		cli, err := ssh.Dial(protocol, exe.JoinHost(os.Getenv(SSH_HOST), os.Getenv(SSH_PORT)), config)
		if err != nil {
			println("Couldn't establish connection to backend server due to\t", err.Error())
			println("Please try again or report if problem persists")
			return ""
		}

		sshClientConnection = cli
	}
	log.Print("New Request")
	//opening socket on remote machine to listen
	//here we are using port no. 0, so let remote side can decide the port to be open
	listener, err := sshClientConnection.Listen(protocol, exe.JoinHost("0.0.0.0", 0))
	if err != nil {
		println("Couldn't forward port to remote machine\t", err.Error())
		return ""
	}

	//It means we have successfully established the connection, lets examine which port assigned to us
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	rfp := &exe.Forwardedport{fmt.Sprint(dest), fmt.Sprint(port), listener, make(exe.Closerch)}
	//Register forwarded ports
	remoteForwardedPorts = append(remoteForwardedPorts, rfp)
	//Remote new generater remote listener
	print(COLOR_GREEN_UNDERLINED)
	println("Successfully Port Forwarding Established")
	println(listener.Addr().String(), "\t->\t", exe.JoinHost("", rfp.DestPort))
	resetConsoleColor()

	go exe.HandleForwardedListener(rfp)

	return port
}

//disconnectRemoteForwardedPort disconnects remotely forwarded port
func disconnectRemoteForwardedPort(dest string) {
	for i, v := range remoteForwardedPorts {
		if v.DestPort == dest {
			v.Closer <- true
			print(COLOR_YELLOW)
			println("Port forwarding disconnected for port ", dest)
			resetConsoleColor()
			//updating remote forwarded port list
			temp := remoteForwardedPorts[i+1:]
			remoteForwardedPorts = remoteForwardedPorts[:i]
			remoteForwardedPorts = append(remoteForwardedPorts, temp...)

			return
		}
	}
	print(COLOR_CYAN)
	println("Couldn't find any port binding for destination port ", dest)
	resetConsoleColor()
}

//This function lists remote connected tunnels
func listConnectedRemoteTunnel() {
	remotehost := os.Getenv("SHREE_REMOTE_HOST")
	println("Remote connected ports")
	print(COLOR_YELLOW)

	for i, v := range remoteForwardedPorts {
		println(i+1, "\t", exe.JoinHost(remotehost, v.SrcPort), "\t->\t", exe.JoinHost("", v.DestPort))
	}
	resetConsoleColor()

}
