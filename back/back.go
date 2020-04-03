package main

import (
	"crypto/md5"
	"database/sql"
	"errors"

	"golang.org/x/crypto/ssh"

	"github.com/devansh42/shree/remote"
	_ "github.com/go-sql-driver/mysql"
)

type Backend struct{}

const databaseurl = ""

func (b *Backend) Auth(user, response *remote.User) (err error) {
	db, err := sql.Open("mysql", databaseurl)
	if err != nil {
		return
	}
	pstmt, err := db.Prepare("select uid,email from users where username=?  limit 1")
	if err != nil {
		return
	}
	rs, err := pstmt.Query(user.Username, hashPasswd(user.Password))
	if err != nil {
		return
	}
	if rs.Next() {
		//Have User
		rs.Scan(&user.Uid, &user.Email)
		pstmt, err = db.Prepare("select uid from users where username=? and password=? limit 1")
		rs, err = pstmt.Query(user.Username, hashPasswd(user.Password))
		if err != nil {

		}
		if !rs.Next() {
			//Auth failed
			return errors.New("401") //User not found invalid credentials
		}
		response = user
		return
	} else {
		//Lets make a new user
		pstmt, err = db.Prepare("insert into users(username,password,email)values(?,?,?)")
		res, err := pstmt.Exec(user.Username, hashPasswd(user.Password), user.Email)
		if err != nil {
			//Handle error
		}
		id, err := res.LastInsertId()
		if err != nil {

		}
		user.IsNew = true
		user.Uid = id
		response = user //New User created
	}
	return
}

func (b *Backend) IssueCertificate(req *remote.CertificateRequest, resp *ssh.Certificate) error {
	err := getCAClient().Call("CA.GetNewCertificate", req, resp)
	//So far we just relaying the request to the ca
	return err
}

func (b *Backend) GetCAPublicCertificate(req *remote.CertificateRequest, resp *ssh.Certificate) error {
	return getCAClient().Call("CA.GetCAPublicCertificate", req, resp)
	//So far we just relaying the request to the ca

}

func (b *Backend) GetCAUserPublicCertificate(req *remote.CertificateRequest, cert *ssh.Certificate) error {
	return getCAClient().Call("CA.GetCAUserPublicCertificate", req, cert)
}

func hashPasswd(b []byte) []byte {
	m := md5.New()

	return m.Sum(b)
}
