package shree

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func generateNewKeyPair() {
	pk, err := rsa.GenerateKey(rand.Reader, 4096)
	handleErr(err)
	y := x509.MarshalPKCS1PrivateKey(pk)
	x := pem.Block{Type: "RSA PRIVATE KEY", Headers: nil, Bytes: y}
	pb, err := ssh.NewPublicKey(&pk.PublicKey)
	handleErr(err)
	pub := ssh.MarshalAuthorizedKey(pb)
	var prb = pem.EncodeToMemory(&x)
	fmt.Println(string(pub))

	fmt.Println(string(prb))
	ioutil.WriteFile("/home/devansh42/.ssh/id_demo", prb, 0400)
	ioutil.WriteFile("/home/devansh42/.ssh/id_demo.pub", pub, 0400)

}
