package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/devansh42/shree/remote"
	"github.com/go-redis/redis"
)

type Backend struct{}

func (b *Backend) Auth(user, response *remote.User) (err error) {
	db := redis.NewClient(&redis.Options{Network: "tcp", Addr: os.Getenv(REDIS_SERVER_ADDR)})
	defer db.Close()
	xkey := hash(sprint("u:", user.Username))
	re := db.Get(xkey)

	if re.Err() != nil {
		//username doesn't exsists
		//Let's create one
		uid := time.Now().Unix()
		r := db.HMSet(hash(sprint("u", uid)), map[string]interface{}{
			"u": user.Username,
			"e": user.Email,
			"p": hash(string(user.Password)),
		})
		_, err = r.Result()

		if err != nil {
			//handle it
			return errors.New("Internal server error")
		}

		//Lets set username bindings
		db.Set(xkey, uid, 0) //Setting uid binding

		response.Email = user.Email
		response.Username = user.Username
		response.Uid = uid
		response.IsNew = true
		return
	}

	uid, _ := re.Int64()

	key := hash(sprint("u", uid))
	resp := db.HGetAll(key)
	res, err := resp.Result()
	if err != nil {
		return errors.New("Internal server error")

	}
	if res["p"] != hash(string(user.Password)) {
		//username/password not matched
		return errors.New("401") //User  found but invalid credentials
	}

	//User Authenticated
	response.Email = res["e"]
	response.Username = res["u"]
	response.Uid = uid

	return
}

func (b *Backend) IssueHostCertificate(req *remote.HostCertificateRequest, resp *remote.CertificateResponse) error {
	err := getCAClient().Call("CA.IssueHostCertificate", req, resp)
	return err
}

func (b *Backend) IssueCertificate(req *remote.CertificateRequest, resp *remote.CertificateResponse) error {
	err := getCAClient().Call("CA.GetNewCertificate", req, resp)
	//So far we just relaying the request to the ca
	return err
}

func (b *Backend) GetCAPublicCertificate(req *remote.CertificateRequest, resp *remote.CertificateResponse) error {
	return getCAClient().Call("CA.GetCAHostPublicKey", req, resp)
	//So far we just relaying the request to the ca

}

func (b *Backend) GetCAUserPublicCertificate(req *remote.CertificateRequest, cert *remote.CertificateResponse) error {
	return getCAClient().Call("CA.GetCAUserPublicKey", req, cert)
}

func hash(s string) string {
	m := md5.New()
	return string(m.Sum([]byte(s)))
}

var sprint = fmt.Sprint
