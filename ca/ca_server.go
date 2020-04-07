package main

import (
	"flag"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	"github.com/devansh42/shree/exe"
)

const (
	CAPRIVATEFILE = "SHREE_CAPRIVATE"
	CAUSERPUBKEY  = "SHREE_CAUSERPUBKEY"
	CAHOSTPRIKEY  = "SHREE_CAHOSTPRIVATE"
	CAHOSTPUBKEY  = "SHREE_CAHOSTPUBKEY"
	CA_ADDR       = "SHREE_CA_ADDR"
)

func parseFlag() {
	uprv := flag.String("uprv", "", "(required) CA user private key")
	upub := flag.String("upub", "", "(required) CA user public key")
	hprv := flag.String("hprv", "", "(required) CA host private key")
	hpub := flag.String("hpub", "", "(required) CA host public key")
	a := flag.String("addr", "", "(required) CA addr")
	l := flag.String("logdir", "App Dir", "CA log directory")
	flag.Parse()

	for !flag.Parsed() {
		//Waiting for parsing
	}

	ar := []string{*uprv, *upub, *hprv, *hpub}
	for _, v := range ar {
		_, err := os.Stat(v)
		if err != nil {
			log.Fatal("Couldn't open file", v, " due to ", err.Error())
		}

	}
	_, _, err := net.SplitHostPort(*a)
	if err != nil {
		log.Fatal("Couldn't parse nwtwork address due to ", err.Error())
	}

	os.Setenv(CA_ADDR, *a)
	os.Setenv(CAPRIVATEFILE, *uprv)
	os.Setenv(CAHOSTPRIKEY, *hprv)
	os.Setenv(CAUSERPUBKEY, *upub)
	os.Setenv(CAHOSTPUBKEY, *hpub)

	u, _ := os.UserHomeDir()
	var ps = string(os.PathSeparator)
	appDir := strings.Join([]string{u, ".shree", "ca", "logs.log"}, ps)
	wr := exe.SetLogFile(*l, appDir)
	log.SetOutput(wr)
}

//StartServer starts rpc server
func StartServer() {
	addr := os.Getenv(CA_ADDR)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Couldn't start ca server due to\t", err.Error())
	}
	log.Print("Listening at ", addr)
	server := rpc.NewServer()
	server.Register(&CA{}) //Registering service

	server.Accept(listener)

}
