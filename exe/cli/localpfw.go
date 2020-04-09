package main

import (
	"fmt"
	"log"
	"net"

	"github.com/devansh42/shree/exe"
)

const keyLocalpfw = "localpfw"

var (
	locallyForwardedPort []*exe.Forwardedport
)

//forwardLocalPort, forwards local port src->dest
//This listens connection from src and relay them to dest
//protocol defines protocol to be used tcp or udp
func forwardLocalPort(protocol string, src, dest int) (err error) {
	listener, err := net.Listen(protocol, exe.JoinHost("", src))
	if err != nil {
		log.Print("Couldn't connect ports due to : ", err.Error())
		return err
	}

	lfp := &exe.Forwardedport{fmt.Sprint(dest), fmt.Sprint(src), listener, make(exe.Closerch)}

	locallyForwardedPort = append(locallyForwardedPort, lfp)

	//Printing status
	print(COLOR_GREEN_UNDERLINED)
	println("Successfully Port Forwarding Established")
	println(listener.Addr().String(), "\t->\t", exe.JoinHost("", lfp.DestPort))
	resetConsoleColor()

	go exe.HandleForwardedListener(lfp) //Handling current Session

	return
}

//disconnectLocalyForwardedPort, disconnects localy forwarded port
func disconnectLocalyForwardedPort(src int) {

	for i, v := range locallyForwardedPort {
		if v.SrcPort == sprint(src) {
			v.Closer <- true //Sending closing signal
			print(COLOR_YELLOW)
			println("Local port forwarding disabled for port ", src)
			resetConsoleColor() //Reseting
			temp := locallyForwardedPort[i+1:]
			locallyForwardedPort = locallyForwardedPort[:i]
			locallyForwardedPort = append(locallyForwardedPort, temp...)
			return
		}
	}

	println(COLOR_CYAN)
	println("Couldn't found the desired port ", src)
	resetConsoleColor()
}

//This function lists local connected tunnels
func listConnectedLocalTunnel() {
	println("Local connected port(s) ", len(locallyForwardedPort), " found")
	print(COLOR_YELLOW)
	for i, v := range locallyForwardedPort {
		println(i+1, "\t", exe.JoinHost("", v.SrcPort), "\t->\t", exe.JoinHost("", v.DestPort))
	}
	resetConsoleColor()

}
