package main

import (
	"io/ioutil"

	"github.com/devansh42/shree/remote"
	"golang.org/x/crypto/ssh"
)

const (
	causerpublickeyfile = "dededed"

	//Contains file uri to the ca host public key
	cahostpublickeyfile = "dedede"
)

type CA struct{} //Just for naming purposes

//GetNewCertificate, gets an certificate with one year valdity
func (c *CA) GetNewCertificate(req *remote.CertificateRequest, cert *ssh.Certificate) (err error) {
	cert, err = getCertificate(req.User.Username, req.PublicKey)
	return
}

//GetCAPublicCertificate, returns public host key of the ca
func (c *CA) GetCAPublicCertificate(req *remote.CertificateRequest, cert *ssh.Certificate) (err error) {
	b, _ := ioutil.ReadFile(cahostpublickeyfile)
	pub, _, _, _, _ := ssh.ParseAuthorizedKey(b)
	cert = pub.(*ssh.Certificate)
	return nil
}

//GetCAUserPublicCertificate, returns public user key of the ca
func (c *CA) GetCAUserPublicCertificate(req *remote.CertificateRequest, cert *ssh.Certificate) (err error) {
	b, _ := ioutil.ReadFile(cahostpublickeyfile)
	pub, _, _, _, _ := ssh.ParseAuthorizedKey(b)
	cert = pub.(*ssh.Certificate)
	return nil
}
