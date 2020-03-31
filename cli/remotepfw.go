package main

//This file contains code for remote port forwarding

import (
	"fmt"
	"io"
	"net"

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
	remoteForwardedPorts []*Forwardedport
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
		println("Error occured while processing Private Ket : ", err.Error())
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
	if !socketCollection.have(keysshclientsocket) {
		config := &ssh.ClientConfig{
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(getClientSigner())},
			HostKeyCallback: getHostCallBack(),
			User:            currentUser.Username}

		cli, err := ssh.Dial(protocol, net.JoinHostPort(SSH_HOST, SSH_PORT), config)
		if err != nil {
			println("Couldn't establish connection to backend server due to\t", err.Error())
			println("Please try again or report if problem persists")
			return
		}

		socketCollection.add(keysshclientsocket, cli)
		sshClientConnection = cli
	}

	//opening socket on remote machine to listen
	//here we are using port no. 0, so let remote side can decide the port to be open
	listener, err := sshClientConnection.Listen(protocol, joinHost("0.0.0.0", 0))
	if err != nil {
		println("Couldn't forward port to remote machine\t", err.Error())
		return
	}
	//It means we have successfully established the connection, lets examine which port assigned to us
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	rfp := &Forwardedport{fmt.Sprint(dest), fmt.Sprint(port), listener}
	//Register forwarded ports
	remoteForwardedPorts = append(remoteForwardedPorts, rfp)
	//Remote new generater remote listener
	print(COLOR_GREEN_UNDERLINED)
	println("Successfully Port Forwarding Established")
	println(listener.Addr().String(), "\t->\t", joinHost("", rfp.DestPort))
	resetConsoleColor()

	handleForwardedListener(rfp)

}

func handleForwardedListener(conn *Forwardedport) {
	listener := conn.Listener
	//Running forever
	for {
		//Making a connection for relaying data to local port
		relayConn, err := net.Dial("tcp", joinHost("", conn.DestPort))
		if err != nil {
			//Couldn't handle this connection
			continue
		}
		acceptedConn, err := listener.Accept() //Accepting connection at remote port
		if err != nil {
			//handle it
			continue
		}
		go handleConnectionIO(acceptedConn, relayConn)

	}
}

//handleConnectionIO handle i/o b/w relayed connection and accepted connections
func handleConnectionIO(acceptedConn, relayConn net.Conn) {
	defer acceptedConn.Close()
	defer relayConn.Close()
	closer := make(chan bool)
	go func() {
		io.Copy(relayConn, acceptedConn)
		closer <- true
	}()

	go func() {
		io.Copy(acceptedConn, relayConn)
		closer <- true
	}()
	<-closer //Whenever it hears a signal it closes both sides of remote connection

}
