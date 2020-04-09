package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"

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
	loadProps()
	app := getCliApp()
	for {
		print(COLOR_YELLOW)
		print(SHREE_CMD_TAG)
		resetConsoleColor()
		reader := bufio.NewReader(os.Stdin)
		line, _, _ := reader.ReadLine()
		s := strings.Split(string(line), " ")
		var ss []string
		ss = append(ss, "shree")
		ss = append(ss, s...)

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
		print(COLOR_RED)
		println("Couldn't open local database, It might happen if another instance of shree is already running or something, please remote this problem")
		println("Reason: ", err.Error())
		resetConsoleColor()
		os.Exit(0) //Exiting from current session
	}
	localdb = db
	currentPeer = newpeer()
	//currentUser = new(remote.User)

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

//setProp, sets props the cli and persists it the db
func setProp(name string, value []byte) bool {
	err := localdb.Put([]byte(sprint("prop:", name)), value, nil)
	if err != nil {
		return false
	}
	return true
}

//getProp, gets props from cli
func getProp(name string) []byte {
	k := []byte(sprint("prop:", name))
	ok, err := localdb.Has(k, nil)
	if !ok || err != nil {
		return nil
	}
	v, _ := localdb.Get(k, nil)
	return v
}

//loadProps loads props and sets environment variable
func loadProps() {
	for _, v := range cliProps {
		vv := getProp(v)
		if vv != nil {
			os.Setenv(v, string(vv))
		}
	}
}

var cliProps = []string{
	SHREE_SSH_ADDR,
	SHREE_BACKEND_ADDR,
}

var (
	print   = fmt.Print
	println = fmt.Println
	sprint  = fmt.Sprint
)
