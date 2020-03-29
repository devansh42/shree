package main

//This file contains code for remote port forwarding

import (
	"io"
	"log"
	"net"
	"strconv"

	"golang.org/x/crypto/ssh"
)

const (
	SSH_PORT           = "2200"
	SSH_HOST           = "ssh.bsnl.online"
	keysshclientsocket = "sshclient"
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
		//fmt.Print(config.ClientVersion)

		cli, err := ssh.Dial(protocol, net.JoinHostPort(SSH_HOST, SSH_PORT), config)
		if err != nil {
			println("Couldn't establish connection to backend server due to\t", err.Error())
			println("Please try again or report if problem persists")
		}
		socketCollection.add(keysshclientsocket, cli)

	}

	//	defer cli.Close()

	listener, err := cli.Listen(protocol, joinHost("0.0.0.0", src)) //opening socket on remote machine to listen
	handleErr(err)

	//	defer listener.Close()
	//go func() {
	defer cli.Close()
	defer listener.Close()
	for {
		oconn, err := net.Dial(protocol, joinHost("localhost", dest))
		handleErr(err)

		iconn, err := listener.Accept()
		handleErr(err)

		go handle(iconn, oconn)

	}
	//}()

}

func joinHost(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func handle(iconn, oconn net.Conn) {
	var chDone = make(chan bool)
	defer iconn.Close()
	go func() {
		_, err := io.Copy(iconn, oconn)
		if err != nil {
			log.Println(err)
		}

		chDone <- true
	}()
	go func() {
		_, err := io.Copy(oconn, iconn)
		if err != nil {
			log.Println(err)
		}
		chDone <- true

	}()
	<-chDone

}
