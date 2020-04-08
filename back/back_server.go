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

const BACKEND_SERVER_ADDR = "BACKEND_SERVER_PORT"
const CA_SERVER_ADDR = "CA_SERVER_ADDR"
const REDIS_SERVER_ADDR = "REDIS_SERVER_ADDR"

func parseFlag() {
	b := flag.String("baddr", "", "(required)  Addrs to start backend server")
	c := flag.String("caddr", "", "(required)  Addrs to CA  server")
	r := flag.String("raddr", "", "(required)  Addrs to the redis server")
	l := flag.String("logdir", "App Directory", "  Directory for server logs")
	flag.Parse()
	for !flag.Parsed() {
		//Waiting for arguments to be passed
	}

	ar := []string{*b, *c, *r}
	//Let's validate network address string
	for _, v := range ar {
		_, _, err := net.SplitHostPort(v)
		if err != nil {
			log.Fatal("Couldn't parse address ", v, " due to ", err.Error())
		}
	}

	os.Setenv(BACKEND_SERVER_ADDR, *b)
	os.Setenv(CA_SERVER_ADDR, *c)
	os.Setenv(REDIS_SERVER_ADDR, *r)
	//Lets try to open dir
	homeDir, _ := os.UserHomeDir()
	fname := strings.Join([]string{homeDir, ".shree", "back"}, string(os.PathSeparator))

	wr := exe.SetLogFile(*l, fname)
	//Logs are written to stdout and the log file
	log.SetOutput(wr)
}

func StartServer() {
	parseFlag()

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

func main() {
	StartServer() //Starting the server
}

const caPort = 8082

func getCAClient() *rpc.Client {
	cli, _ := rpc.Dial("tcp", os.Getenv(CA_SERVER_ADDR))
	return cli
}
