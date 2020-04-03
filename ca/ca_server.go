package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

const ServerPort = 8082

//StartServer starts rpc server
func StartServer() {
	listener, err := net.Listen("tcp", net.JoinHostPort("", fmt.Sprint(ServerPort)))
	if err != nil {
		log.Fatal("Couldn't start ca server due to\t", err.Error())
	}
	log.Print("Listening at ", net.JoinHostPort("", fmt.Sprint(ServerPort)))
	server := rpc.NewServer()
	server.Register(&CA{}) //Registering service

	server.Accept(listener)

}
