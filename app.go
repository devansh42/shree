package shree

import (
	"os"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type peer struct {
	pid uint64
}

const shreePort = 1234

func newpeer() *peer {
	p := new(peer)
	p.pid = uint64(time.Now().Unix()) //Timestamp
	return p
}

var currentPeer *peer
var localdb *leveldb.DB
var socketCollection *socketcollection

//This method performs app initialization
func initApp() {
	//Making app folder
	appdir := getAppDir()

	_, err := os.Stat(appdir)
	if os.IsNotExist(err) { //Make one if doesn't exists
		os.Mkdir(appdir, 0600)
	}

	//Opening App folder
	db, err := leveldb.OpenFile(getAppFile("state"), nil) //Opening state database
	if err != nil {

	}
	localdb = db
	socketCollection = new(socketcollection) //Making new socker collector
	currentPeer = newpeer()
}

func getAppFile(fs ...string) string {
	f := append([]string{getAppDir()}, fs...)
	return strings.Join(f, string(os.PathSeparator))
}
func getAppDir() string {

	home, _ := os.UserHomeDir()
	return home + string(os.PathSeparator) + ".shree"
}
