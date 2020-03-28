package main

import (
	"io/ioutil"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
)

const homedir = "/home/devansh42"

func TestGetCertificate(t *testing.T) {

	//Setting environmental variable for testing purposes
	os.Setenv(CAPRIVATEFILE, homedir+"/.ssh/ca_user_key")
	initCA()
	key := getTestPublicKey()
	certificate, err := getCertificate("devansh42", key)
	panicErr(err)
	b := ssh.MarshalAuthorizedKey(certificate)
	t.Log(string(b))
	ioutil.WriteFile(homedir+"/.ssh/demo-cert.pub", b, 0600)

}

//gets public key from local host
func getTestPublicKey() ssh.PublicKey {
	b, _ := ioutil.ReadFile(homedir + "/.ssh/id_demo.pub")
	key, _, _, _, err := ssh.ParseAuthorizedKey(b)
	//ssh.Parse
	panicErr(err)
	return key
}
