package main

import (
	"github.com/devansh42/shree/remote"
	"golang.org/x/crypto/ssh"
)

const (
	CAUSERPUBKEY = "SHREE_CAUSERPUBKEY"
	CAHOSTPUBKEY = "SHREE_CAHOSTPUBKEY"
)

type CA struct{} //Just for naming purposes

//GetNewCertificate, gets an certificate with one year valdity
func (c *CA) GetNewCertificate(req *remote.CertificateRequest, resp *remote.CertificateResponse) (err error) {
	pubkey, _, _, _, err := ssh.ParseAuthorizedKey(req.PublicKey)
	if err != nil {
		return err
	}
	cert, err := getCertificate(req.User.Username, pubkey)
	if err != nil {
		return err
	}
	resp.Bytes = ssh.MarshalAuthorizedKey(cert)
	return
}

//GetCAUserPublicKey, returns public user key of the ca
func (c *CA) GetCAUserPublicKey(req *remote.CertificateRequest, cert *remote.CertificateResponse) (err error) {
	cert.Bytes = marshaledUserPublicKey
	return nil
}

//GetCAHostPublicKey, returns public host key of the ca
func (c *CA) GetCAHostPublicKey(req *remote.CertificateRequest, cert *remote.CertificateResponse) (err error) {
	cert.Bytes = marshaledHostPublicKey
	return nil
}
