package main

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	CAPRIVATEFILE = "CAPRIVATE"
)

func main() {
	//Main function to launch

}

func initCA() {
	onceLock = new(sync.Once)
	onceLock.Do(getCAPrivateKey) //Loads certificates private key in memory
}

var onceLock *sync.Once
var privateKeySigner ssh.Signer

//getCAPrivateKey loads ca private key memory
//to be used with sync.Once
func getCAPrivateKey() {
	prpath := os.Getenv(CAPRIVATEFILE)
	b, err := ioutil.ReadFile(prpath)
	panicErr(err)
	signer, err := ssh.ParsePrivateKey(b)
	panicErr(err)
	privateKeySigner = signer
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}

}

//getCertificate signes the certificate with validity of 1 yr
//it only return non nil if any problem occured in signing process
func getCertificate(username string, tobesigned ssh.PublicKey) (*ssh.Certificate, error) {
	cert := new(ssh.Certificate)
	cert.Key = tobesigned
	cert.ValidPrincipals = []string{username} //Valid  principal is the username of the user
	now := time.Now()
	cert.Serial = uint64(now.Unix())
	cert.CertType = ssh.UserCert //Sets certificate type
	//Valid for a year
	cert.ValidBefore = uint64(now.Add(time.Hour * 24 * 365).Unix())
	//Permits only port forwarding
	cert.Extensions = map[string]string{"permit-port-forwarding": ""}
	err := cert.SignCert(rand.Reader, privateKeySigner)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
