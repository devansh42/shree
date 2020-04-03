package main

import (
	"fmt"
	"net"
	"net/rpc"
)

const ServerPort = 8082

//StartServer starts rpc server
func StartServer() {
	listener, err := net.Listen("tcp", net.JoinHostPort("", fmt.Sprint(ServerPort)))
	panicErr(err)
	server := rpc.NewServer()
	server.Register(&CA{}) //Registering service
	server.Accept(listener)

}
