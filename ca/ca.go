package main

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	parseFlag()
	initCA()
	StartServer() //Starting the server
}

func initCA() {
	onceLock = new(sync.Once)
	onceLock.Do(func() {
		getCAPrivateKey()
		getCAHostPubliKey()
		getCAUserPubliKey()
	}) //Loads certificates private key in memory
}

var onceLock *sync.Once
var hostPrivatekey, privateKeySigner ssh.Signer
var marshaledHostPublicKey, marshaledUserPublicKey []byte

//getCAUserPubliKey loads ca user public key
func getCAUserPubliKey() {
	p := os.Getenv(CAUSERPUBKEY)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal("Couldn't load user public key\t", err)
	}
	marshaledUserPublicKey = b
}

//getCAHostPubliKey loads ca host public key
func getCAHostPubliKey() {
	p := os.Getenv(CAHOSTPUBKEY)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal("Couldn't load host public key\t", err)
	}
	marshaledHostPublicKey = b
}

//Get Hosts private key
func getCAHostPrivateKey() {
	p := os.Getenv(CAHOSTPRIKEY)
	b, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal("Couldn't load host private key\t", err)
	}
	pr, err := ssh.ParsePrivateKey(b)
	if err != nil {
		log.Fatal("Couldn't parse host private key")
	}
	hostPrivatekey = pr
}

//getCAPrivateKey loads ca private key memory
//to be used with sync.Once
func getCAPrivateKey() {
	prpath := os.Getenv(CAPRIVATEFILE)
	b, err := ioutil.ReadFile(prpath)
	if err != nil {
		log.Println("Couldn't read ca private key", err.Error())
	}
	signer, err := ssh.ParsePrivateKey(b)
	if err != nil {
		log.Println("Couldn't parse private key ", err.Error())
	}
	privateKeySigner = signer
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}

}

//getCertificate signes the certificate with validity of 1 yr
//it only return non nil if any problem occured in signing process
func getCertificate(username string, tobesigned ssh.PublicKey, certType uint32) (*ssh.Certificate, error) {
	cert := new(ssh.Certificate)
	cert.Key = tobesigned
	cert.ValidPrincipals = []string{username} //Valid  principal is the username of the user
	now := time.Now()
	cert.Serial = uint64(now.Unix())
	cert.CertType = certType //Sets certificate type
	//Valid for a year
	cert.ValidBefore = uint64(now.Add(time.Hour * 24 * 365).Unix())
	//Permits only port forwarding
	cert.Extensions = map[string]string{"permit-port-forwarding": ""}
	var signer ssh.Signer
	if certType == ssh.HostCert {
		signer = hostPrivatekey
	} else {
		signer = privateKeySigner
	}
	err := cert.SignCert(rand.Reader, signer)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
