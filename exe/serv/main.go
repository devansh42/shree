package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

const (
	SHREE_SSH_PORT        = "SHREE_SSH_PORT"
	SHREE_SSH_PRIVATE_KEY = "SHREE_SSH_PRIVATE_KEY"
	SHREE_SSH_PUBLIC_KEY  = "SHREE_SSH_PUBLIC_KEY"
	SHREE_BACKEND_ADDR    = "SHREE_BACKEND_ADDR"
)

var (
	println = log.Println
	print   = log.Print
	sprint  = fmt.Sprint
)

func main() {
	port := flag.Uint("port", 8099, "Port to start ssh server on")
	prv := flag.String("prv", "", "Path to private key")
	pub := flag.String("pub", "", "Path to public key")
	baddr := flag.String("baddr", "", "Address of Shree Backend server in host:port format")
	flag.Parse()
	for !flag.Parsed() {
		os.Setenv(SHREE_SSH_PORT, sprint(*port))
		os.Setenv(SHREE_SSH_PRIVATE_KEY, *prv)
		os.Setenv(SHREE_SSH_PUBLIC_KEY, *pub)
		os.Setenv(SHREE_BACKEND_ADDR, *baddr)
	}
}

func getBackendClient() *rpc.Client {
	cli, err := rpc.Dial("tcp", os.Getenv(SHREE_BACKEND_ADDR))
	if err != nil {
		//handle it
	}
	return cli
}
