package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
)

const BACKEND_SERVER_ADDR = "BACKEND_SERVER_PORT"
const CA_SERVER_ADDR = "CA_SERVER_ADDR"
const REDIS_SERVER_ADDR = "REDIS_SERVER_ADDR"

func StartServer() {
	addr := os.Getenv(BACKEND_SERVER_ADDR)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Couldn't start tcp server due to\t", err.Error())
	}
	log.Println("Listening at ", addr)
	server := rpc.NewServer()
	server.Register(&Backend{}) //Registering services
	server.Accept(listener)

}

func main() {}

const caPort = 8082

func getCAClient() *rpc.Client {
	cli, _ := rpc.Dial("tcp", os.Getenv(CA_SERVER_ADDR))
	return cli
}
