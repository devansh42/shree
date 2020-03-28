package main

import (
	"github.com/devansh42/shree/remote"
	"golang.org/x/crypto/ssh"
)

type CA struct{} //Just for naming purposes

func (c *CA) GetNewCertificate(req *remote.CertificateRequest, cert *ssh.Certificate) (err error) {
	//Gets an certificate with one year valdity
	cert, err = getCertificate(req.User.Username, req.PublicKey)
	return
}
