package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/devansh42/shree/remote"
)

const homedir = "/home/devansh42"

func TestGetCertificate(t *testing.T) {
	initTestEnviroment()
	ca := new(CA)
	b, _ := ioutil.ReadFile("../../keys/id_user.pub")

	user := new(remote.User)
	user.Uid = 1
	user.Username = "devansh42"
	cert := new(remote.CertificateResponse)
	err := ca.GetNewCertificate(&remote.CertificateRequest{*user, b}, cert)
	if err != nil {
		t.Error("Couldn't issue new certificate")
	}
	t.Log("Here is the certificate\t", string(b))
}

func TestGetHostPublicKey(t *testing.T) {
	initTestEnviroment()
	ca := new(CA)
	resp := new(remote.CertificateResponse)
	err := ca.GetCAHostPublicKey(nil, resp)
	if err != nil {
		t.Error("Couldn't load host public key")
	}

	t.Log("Here is the host public key\t", string(resp.Bytes))
}

func TestGetUserPublicKey(t *testing.T) {
	initTestEnviroment()
	ca := new(CA)
	resp := new(remote.CertificateResponse)
	err := ca.GetCAUserPublicKey(nil, resp)
	if err != nil {
		t.Error("Couldn't load user public key")
	}

	t.Log("Here is the user public key\t", string(resp.Bytes))
}

func initTestEnviroment() {
	os.Setenv(CAUSERPUBKEY, "../../keys/ca_user_key.pub")
	os.Setenv(CAHOSTPUBKEY, "../../keys/ca_host_key.pub")
	os.Setenv(CAPRIVATEFILE, "../../keys/ca_user_key")

	initCA()

}
