package main

//This file contains code for remote port forwarding

import (
	"fmt"
	"net"
	"os"

	"github.com/devansh42/shree/exe"

	"golang.org/x/crypto/ssh"
)

const (
	SSH_PORT           = "2200"
	SSH_HOST           = "ssh.bsnl.online"
	keysshclientsocket = "sshclient"
)

var (
	//Connection for long lived tcp connections
	sshClientConnection *ssh.Client
	//remoteForwardedPorts lists remote forwarded port
	remoteForwardedPorts []*exe.Forwardedport
)

//getClientSigner returns signed certificate for authentication
func getClientSigner() ssh.Signer {
	pass := askForPassword()
	havepub, havepr, havecert, pki := searchForPKICredentials(currentUser.Uid)
	if havepub == havepr == havecert == true {

	}
	cert, _, _, _, err := ssh.ParseAuthorizedKey(pki.cert)
	if err != nil {
		//handle
	}
	prv, err := ssh.ParsePrivateKeyWithPassphrase(pki.prv, pass)
	if err != nil {
		println("Error occured while processing Private Key : ", err.Error())
	}
	signer, err := ssh.NewCertSigner(cert.(*ssh.Certificate), prv)
	if err != nil {

	}
	return signer
}

func getHostCallBack() ssh.HostKeyCallback {
	cert := getServerCertificate()
	return ssh.FixedHostKey(cert)
}

//forwardRemotePort, forwards remote port src->dest
//it binds dest port on localhost with src port on remote machine
func forwardRemotePort(protocol string, src, dest int) {
	if sshClientConnection == nil {
		config := &ssh.ClientConfig{
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(getClientSigner())},
			HostKeyCallback: getHostCallBack(),
			User:            currentUser.Username}

		cli, err := ssh.Dial(protocol, exe.JoinHost(SSH_HOST, SSH_PORT), config)
		if err != nil {
			println("Couldn't establish connection to backend server due to\t", err.Error())
			println("Please try again or report if problem persists")
			return
		}

		sshClientConnection = cli
	}

	//opening socket on remote machine to listen
	//here we are using port no. 0, so let remote side can decide the port to be open
	listener, err := sshClientConnection.Listen(protocol, exe.JoinHost("0.0.0.0", 0))
	if err != nil {
		println("Couldn't forward port to remote machine\t", err.Error())
		return
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

	exe.HandleForwardedListener(rfp)

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