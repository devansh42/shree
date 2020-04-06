package main

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"
	"os"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/devansh42/shree/remote"
)

//File contains

//TestAuthSignupUser tests user signup process
func TestAuthSignupUser(t *testing.T) {
	initApp()
	defer cleanup()
	startDemoRPCServer(t) //Starting demo rpc server
	cli := getBackendClient()
	defer cli.Close()
	user := &remote.User{Email: "xyz", IsNew: true, Uid: 1, Username: "xyz", Password: hash([]byte("password"))}
	isnew, err := authUser(user, cli)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(currentUser, isnew) //As after login it sets current user
}

//TestManageCertificateWithCredentials tests the manage certificate function
//if user has credentials locally available but don't have certificate
func TestManageCertificateWithCredentials(t *testing.T) {
	initApp()
	defer cleanup()
	defer cleanupDemoCredentials()
	startDemoRPCServer(t) //Starting demo server

	bpass := []byte("hello1234")

	user := new(remote.User)
	user.Uid = 1
	user.Username = "devansh42"

	bprv, bpub := generateNewKeyPair(bpass)
	writePairToDB(bpub, bprv, 1)
	manageCertificate(user, bpass)

}

//TestManageCertificateWithoutCredentials tests the manage certificate function
//if user havn't credentials locally available
func TestManageCertificateWithoutCredentials(t *testing.T) {
	initApp()
	defer cleanup()
	defer cleanupDemoCredentials()
	startDemoRPCServer(t) //Starting demo server

	bpass := []byte("hello1234")

	user := new(remote.User)
	user.Username = "devansh42"
	user.Uid = 1
	generateAndPersistCredentialsForTest(user, bpass, t)
	manageCertificate(user, bpass)

}

//TestManageCertificateFoundCredentials tests conditions when found valid credentials
func TestManageCertificateFoundCredentials(t *testing.T) {
	initApp()
	defer cleanup()
	defer cleanupDemoCredentials()
	startDemoRPCServer(t) //Starting demo server

	bpass := []byte("hello1234")

	user := new(remote.User)
	user.Username = "devansh42"
	user.Uid = 1

	// bprv, bpub := generateNewKeyPair(bpass)
	// writePairToDB(bpub, bprv, 1)
	manageCertificate(user, bpass)

}

//Export for backend related server
type Backend struct{}

//Auth pretends the auth method for the service
func (b *Backend) Auth(user, resp *remote.User) error {
	log.Print("User recived ", user)
	//x := remote.User{user.Email, user.IsNew, user.Password, user.Uid, user.Username}
	resp.Email = user.Email
	resp.IsNew = user.IsNew
	resp.Password = user.Password
	resp.Uid = user.Uid
	resp.Username = user.Username

	if user.IsNew { //For sake of testing is isnew is true it means user wants to test signup

	} else {
	}

	return nil
}

//IssueCertificate issues new certificate related to user request, this implementation is specific to this implementation
func (b *Backend) IssueCertificate(req *remote.CertificateRequest, resp *remote.CertificateResponse) error {
	cert := new(ssh.Certificate)
	cert.ValidPrincipals = []string{req.User.Username}
	cert.CertType = ssh.UserCert
	cert.ValidBefore = uint64(time.Now().Add(time.Minute * 60 * 24 * 365).Unix())
	prvb, err := ioutil.ReadFile("../../keys/ca_user_key")
	if err != nil {
		return err
	}

	prv, err := ssh.ParsePrivateKey(prvb)
	if err != nil {
		return err
	}
	pub, _, _, _, _ := parseauthkey(req.PublicKey)
	cert.Key = pub
	err = cert.SignCert(rand.Reader, prv)
	if err != nil {
		return err
	}
	resp.Bytes = marshalauthkey(cert)

	return nil
}

func startDemoRPCServer(t *testing.T) {

	os.Setenv(SHREE_BACKEND_ADDR, ":6500") //Starting demo rpc server, inorder to work with getBackendClient Method
	listener, err := net.Listen("tcp", ":6500")
	if err != nil {
		t.Fatal("Couldn't start server at given port")
	}
	rpc.Register(&Backend{})
	t.Log("RPC server listening on :6500")
	go rpc.Accept(listener)
}

//generateAndPersistCredentialsForTest generates credentials for testing
//and issue a new certificate and persist it
//grpc server should already be started
func generateAndPersistCredentialsForTest(user *remote.User, bpass []byte, t *testing.T) {
	cert := new(remote.CertificateResponse)
	bprv, bpub := generateNewKeyPair(bpass)
	t.Log("New Key Pair Generated")
	writePairToDB(bpub, bprv, 1)
	t.Log("Key Pair written to the localdb")
	cli := getBackendClient()
	err := cli.Call("Backend.IssueCertificate", &remote.CertificateRequest{User: *user, PublicKey: bpub}, cert)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("New certificate generated")
	localdb.Put([]byte(sprint(keycertkey, 1)), cert.Bytes, nil)
	t.Log("Certificate written to db")
}

//cleanupDemoCredentials, cleanups demo credentials
func cleanupDemoCredentials() {

	localdb.Delete([]byte(sprint(keycertkey, 1)), nil)
	localdb.Delete([]byte(sprint(keyprvkey, 1)), nil)
	localdb.Delete([]byte(sprint(keypubkey, 1)), nil)

}
