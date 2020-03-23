package shree

//This file contains code for remote port forwarding

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"golang.org/x/crypto/ssh"
)

const (
	SSH_PORT = 22
)

func readSSHKey() ssh.Signer {
	path := "/home/devansh42/.ssh/id_gs"
	b, _ := ioutil.ReadFile(path)
	prv, err := ssh.ParsePrivateKey(b)
	handleErr(err)

	return prv

}

//forwardRemotePort, forwards remote port src->dest
//it binds dest port on localhost with src port on remote machine
func forwardRemotePort(protocol string, src, dest int) {
	s := readSSHKey()

	fmt.Print(string(ssh.MarshalAuthorizedKey(s.PublicKey())))

	config := &ssh.ClientConfig{
		//	HostKeyAlgorithms: []string{"ecdsa-sha2-nistp256"},
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(s)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            "igniterab"}
	fmt.Print(config.ClientVersion)

	cli, err := ssh.Dial(protocol, joinHost("ec2.bsnl.online", SSH_PORT), config)
	handleErr(err)
	//	defer cli.Close()

	listener, err := cli.Listen(protocol, joinHost("localhost", src)) //opening socket on remote machine to listen
	handleErr(err)

	//	defer listener.Close()
	go func() {
		defer cli.Close()
		defer listener.Close()
		for {
			oconn, err := net.Dial(protocol, joinHost("localhost", dest))
			handleErr(err)

			iconn, err := listener.Accept()
			handleErr(err)

			go handle(iconn, oconn)

		}
	}()

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
