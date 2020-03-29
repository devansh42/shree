package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/devansh42/shree/remote"

	"golang.org/x/crypto/ssh/terminal"

	"golang.org/x/crypto/ssh"
)

const (
	keyprvkey            = "prvkey"
	keypubkey            = "pubkey"
	keycertkey           = "certkey"
	hostCertificateURL   = ""
	keyservercertificate = "cert_server"
)

type pkicredentials struct {
	prv, pub, cert []byte
}

var (
	marshalauthkey = ssh.MarshalAuthorizedKey
	parseauthkey   = ssh.ParseAuthorizedKey

	//decryptPrivateKey decrypts private key encrypted with password
	decryptPrivateKey = ssh.ParsePrivateKeyWithPassphrase
)

//generateNewKeyPair generates new rsa priavte key and ssh key pair
//This pair can be saved into db
func generateNewKeyPair(passphrase []byte) (rsaprv []byte, sshpub []byte) {
	pk, err := rsa.GenerateKey(rand.Reader, 4096)

	y := x509.MarshalPKCS1PrivateKey(pk)
	pb, err := ssh.NewPublicKey(&pk.PublicKey)
	handleErr(err)
	x, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", y, passphrase, x509.PEMCipherAES256)
	sshpub = ssh.MarshalAuthorizedKey(pb)
	var prb = pem.EncodeToMemory(x)
	return prb, sshpub

}

//Searches for certificate,pub key and pr key for given user
func searchForPKICredentials(uid int64) (havepub, havepr, havecert bool, pkic pkicredentials) {
	pr, err := localdb.Get([]byte(fmt.Sprint(keyprvkey, uid)), nil)
	havepr = err == nil
	pu, err := localdb.Get([]byte(fmt.Sprint(keypubkey, uid)), nil)
	havepub = err == nil
	cert, err := localdb.Get([]byte(fmt.Sprint(keycertkey, uid)), nil)
	havecert = err == nil
	pkic = pkicredentials{pr, pu, cert}
	return
}

//writePairToDB writes key pair to localdb
func writePairToDB(pub, prv []byte, uid int64) {
	localdb.Put([]byte(fmt.Sprint(keypubkey, uid)), pub, nil)
	localdb.Put([]byte(fmt.Sprint(keyprvkey, uid)), prv, nil)
}

//askForPassword for password and returns it
func askForPassword() []byte {
	println("\nEnter password to continue\t")
	b, err := terminal.ReadPassword(1)
	if err != nil {
		//handle this error
	}
	return b
}

//fetchServerCertificateAndPersist, fetches certificate from default certificate repo
func fetchServerCertificateAndPersist() (cert ssh.PublicKey, err error) {
	println("Fetching Server Certifcate....")
	cert = new(ssh.Certificate)
	err = getBackendClient().Call("Backend.GetCAPublicCertificate", new(remote.CertificateRequest), cert)
	if err != nil {
		return nil, err
	}
	println("Cetificate Fetched\nFingerprint\n", ssh.FingerprintLegacyMD5(cert))
	localdb.Put([]byte(keyservercertificate), marshalauthkey(cert), nil)
	return cert, nil
}

func getServerCertificate() (cert ssh.PublicKey) {
	//Let's search in localdb
	bc, err := localdb.Get([]byte(keyservercertificate), nil)
	if err != nil {
		cert, err = fetchServerCertificateAndPersist()

	} else {
		cert, _, _, _, err = parseauthkey(bc)
		if err != nil { //Couldn't parse
			cert, err = fetchServerCertificateAndPersist()

		}
	}
	return
}

//hash makes md5 hash
func hash(b []byte) []byte {
	m := md5.New()
	return m.Sum(b)
}
