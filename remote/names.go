package remote

import (
	"golang.org/x/crypto/ssh"
)

//User represents a shree user
type User struct {
	Email    string `json:email`
	Password []byte `json:password`
	Uid      int64  `json:uid`
	Username string `json:username`
	IsNew    bool   `json:isnew`
}

//CertificateRequest contains requests for certificate generation
type CertificateRequest struct {
	User      User
	PublicKey ssh.PublicKey
}

type CertificateResponse struct {
	Bytes []byte
}
