package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"

	"github.com/devansh42/shree/remote"

	"github.com/syndtr/goleveldb/leveldb"
)

type peer struct {
	pid uint64
}

const (
	shreePort     = 1234
	SHREE_CMD_TAG = "shree> "
)

func cleanup() {
	defer resetConsoleColor()
	defer func() {
		if localdb != nil {
			localdb.Close()
		}
	}()

}

func main() {
	defer cleanup()

	initApp()
	welComeMsg()
	app := getCliApp()
	for {
		print(SHREE_CMD_TAG)
		reader := bufio.NewReader(os.Stdin)
		line, _, _ := reader.ReadLine()
		s := strings.Split(string(line), " ")
		var ss []string
		ss = append(ss, "shree")
		ss = append(ss, s...)
		println(string(line))
		app.Run(ss)
	}
}

func newpeer() *peer {
	p := new(peer)
	p.pid = uint64(time.Now().Unix()) //Timestamp
	return p
}

var currentPeer *peer
var localdb *leveldb.DB

//This method performs app initialization
//It open database and setups application specific content
func initApp() {
	//Making app folder
	appdir := getAppDir()

	_, err := os.Stat(appdir)
	if os.IsNotExist(err) { //Make one if doesn't exists
		os.Mkdir(appdir, 0700)
	}

	//Opening App folder
	db, err := leveldb.OpenFile(getAppFile("state"), nil) //Opening state database
	if err != nil {

	}
	localdb = db
	currentPeer = newpeer()
	currentUser = new(remote.User)
}

func getAppFile(fs ...string) string {
	f := append([]string{getAppDir()}, fs...)
	return strings.Join(f, string(os.PathSeparator))
}
func getAppDir() string {

	home, _ := os.UserHomeDir()
	return home + string(os.PathSeparator) + ".shree"
}

const SHREE_BACKEND_ADDR = "SHREE_BACKEND_ADDR"

func getBackendClient() *rpc.Client {

	cli, err := rpc.Dial("tcp", os.Getenv(SHREE_BACKEND_ADDR))
	if err != nil {
		//handle error
		println("Couldn't reach to backend server\nTry again later")
		return nil
	}
	return cli
}

var (
	print   = fmt.Print
	println = fmt.Println
	sprint  = fmt.Sprint
)
