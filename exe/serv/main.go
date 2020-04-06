package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"time"

	"github.com/devansh42/shree/remote"
	"golang.org/x/crypto/ssh"
)

const (
	SHREE_SSH_PORT        = "SHREE_SSH_PORT"
	SHREE_SSH_PRIVATE_KEY = "SHREE_SSH_PRIVATE_KEY"
	SHREE_SSH_PUBLIC_KEY  = "SHREE_SSH_PUBLIC_KEY"
	SHREE_BACKEND_ADDR    = "SHREE_BACKEND_ADDR"
	SHREE_HOST_PRINCIPAL  = "SHREE_HOST_PRINCIPAL"
)

var (
	println = log.Println
	print   = log.Print
	sprint  = fmt.Sprint
)

func main() {
	appdir := getAppDir()
	port := flag.Uint("port", 8099, "Port to start ssh server on")
	prv := flag.String("prv", sprint(appdir, ps, "id_host"), "Path to private key")
	pub := flag.String("pub", sprint(appdir, ps, "id_host.pub"), "Path to public key")
	baddr := flag.String("baddr", "", "Address of Shree Backend server in host:port format")
	host := flag.String("host", "", "Host address of this instance to included in certificate principal")
	flag.Parse()
	for !flag.Parsed() {
		//Waiting for command line argument passing
	}
	os.Setenv(SHREE_SSH_PORT, sprint(*port))
	os.Setenv(SHREE_SSH_PRIVATE_KEY, *prv)
	os.Setenv(SHREE_SSH_PUBLIC_KEY, *pub)
	os.Setenv(SHREE_BACKEND_ADDR, *baddr)
	os.Setenv(SHREE_HOST_PRINCIPAL, *host)
	if *baddr == "" || *host == "" {
		log.Fatal("Required args not found, baddr or host")
	}
	initApp()
	initServer() //Starts the testing ssh server
}

func getBackendClient() *rpc.Client {
	cli, err := rpc.Dial("tcp", os.Getenv(SHREE_BACKEND_ADDR))
	if err != nil {
		//handle it
	}
	return cli
}

func initApp() {
	log.Println("Initializing App")
	dirName := getAppDir()
	_, err := os.Stat(dirName)
	if os.IsNotExist(err) {
		os.Mkdir(dirName, 0700)
	}
	checkForCredentials()
}

func getAppDir() string {
	h, _ := os.UserHomeDir()

	dirName := sprint(h, ps, ".shree")
	return dirName
}

//checkForCertificate, checks for certificate and request if not found
func checkForCredentials() {
	log.Print("Checking for Credentials...")
	prvn := os.Getenv(SHREE_SSH_PRIVATE_KEY)
	pubn := os.Getenv(SHREE_SSH_PUBLIC_KEY)
	prvb, err := rF(prvn)
	if err != nil {
		//file not found
		prvb, _ = generateCrendentials()
		//Lets write it to fs
		ioutil.WriteFile(prvn, prvb, 0400)
	}
	s, err := ssh.ParsePrivateKey(prvb)
	if err != nil {
		//Couldn't parse
		log.Print("Couldn't parse private due to : ", err.Error())

		prvb, _ = generateCrendentials()
		log.Println("New Credentials generated")
		//Lets write it to fs
		ioutil.WriteFile(prvn, prvb, 0400)

	}
	pubkey := s.PublicKey()
	if err != nil {
		log.Fatal("Couldn't derive new public key ", err.Error())
	}
	marshledPubkey := ssh.MarshalAuthorizedKey(pubkey)
	err = ioutil.WriteFile(pubn, marshledPubkey, 0400)
	if err != nil {
		//Couln't generate
	}
	certpath := sprint(getAppDir(), ps, "id_host-cert.pub")
	cb, err := rF(certpath)
	if err != nil {
		//Don't have certifiate`
		//Let's check for public key
		log.Println("Fetching Certificate..")
		b := fetchCertificate(marshledPubkey)
		log.Print("Certificate Fetched")
		hostCertifiate = getCertificateFromBytes(b, marshledPubkey)
	} else {

		hostCertifiate = getCertificateFromBytes(cb, marshledPubkey)
	}
	log.Print("Credential checked")
}

func getCertificateFromBytes(b, mpubkey []byte) *ssh.Certificate {
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey(b)
	if err != nil {

		log.Fatal("Couldn't parse the pubkey : ", err.Error())
	}

	x := pubkey.(*ssh.Certificate)
	// log.Print(string(ssh.MarshalAuthorizedKey(x.Key)))
	// log.Print(string(mpubkey))

	if !bytes.Equal(mpubkey, ssh.MarshalAuthorizedKey(x.Key)) {
		log.Fatal("Invalid certificate")
	}
	if x.ValidBefore < uint64(time.Now().Unix()) {
		log.Fatal("Expired Certificate")
	}
	return x
}

//fetchCertificate, requests a certificate from ca
//pub is the public key of this server
func fetchCertificate(pub []byte) (certificateBytes []byte) {

	cli := getBackendClient()
	prin := ""
	resp := new(remote.CertificateResponse)
	err := cli.Call("Backend.IssueHostCertificate", &remote.HostCertificateRequest{PublicKey: pub, Principal: prin}, resp)
	if err != nil {
		log.Fatal("Couldn't fetch the certificate : ", err.Error())
	}
	fn := sprint(getAppDir(), ps, "id_host-cert.pub")
	err = ioutil.WriteFile(fn, resp.Bytes, 0400) //Readonly file
	if err != nil {
		log.Fatal("Couldn't write certificate to fs : ", err.Error())
	}
	return resp.Bytes
}

//generateCredentials, These are generated only ones for a given instance
func generateCrendentials() (prv []byte, pub []byte) {
	log.Print("Generating new credentials....")

	pk, err := rsa.GenerateKey(rand.Reader, 4096)

	y := x509.MarshalPKCS1PrivateKey(pk)
	pb, err := ssh.NewPublicKey(&pk.PublicKey)
	if err != nil {
		//handle error
	}
	pbl := &pem.Block{Bytes: y, Type: "RSA PRIVATE KEY"}
	sshpub := ssh.MarshalAuthorizedKey(pb)
	var prb = pem.EncodeToMemory(pbl)
	log.Println("New Credential generated")
	return prb, sshpub
}

var ps = string(os.PathSeparator)
var hostCertifiate *ssh.Certificate
var rF = ioutil.ReadFile
