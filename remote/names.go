package remote

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
	PublicKey []byte
}

type CertificateResponse struct {
	Bytes []byte
}
