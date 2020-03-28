package main

import (
	"fmt"
	"net"
	"net/rpc"
)

const ServerPort = 8089

func StartServer() {
	listener, err := net.Listen("tcp", net.JoinHostPort("", fmt.Sprint(ServerPort)))
	panicErr(err)
	server := rpc.NewServer()
	server.RegisterName("Backend", Backend{}) //Registering services
	server.Accept(listener)

}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {}

const caPort = 8082

func getCAClient() *rpc.Client {
	cli, _ := rpc.Dial("tcp", net.JoinHostPort("", fmt.Sprint(caPort)))
	return cli
}
