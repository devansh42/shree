package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/rpc"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/devansh42/shree/remote"
)

const keyuser = "user"

//currentUser, is the current user
var currentUser *remote.User

//getCurrentUserInfo, displays current user info
func getCurrentUserInfo() {
	defer resetConsoleColor()
	if currentUser != nil {
		print(COLOR_YELLOW)
		println("Username:\t", currentUser.Username)
		println("Email:\t", currentUser.Email)
		println("User Id:\t", currentUser.Uid)
		return
	}
	print(COLOR_RED)
	println("You are un-authenticated user, please authenticate yourself with \"sign me in\" ")

}

//authUser authenticates user it creates if doesn't exists
func authUser(user *remote.User, cli *rpc.Client) (bool, error) {
	var resp = new(remote.User)
	err := cli.Call("Backend.Auth", user, resp)
	if err != nil {
		switch err.Error() {
		case "401":
			return false, errors.New("Invalid Username/Password")
		default:
			return false, errors.New("Unknown Error:\t" + err.Error())
		}
	}
	currentUser = resp //Is the current logined user
	user.Uid = resp.Uid
	//Lets overwrite the user
	b, _ := json.Marshal(resp)
	localdb.Put([]byte(keyuser), b, nil)
	return resp.IsNew, nil
}

func manageCertificate(user *remote.User, pass []byte) {
	var pub, prv []byte
	dpub, err := validatelocalCertificate(user, pass)
	if err != nil {
		print(COLOR_RED)
		println(err.Error())
		resetConsoleColor()
		type certificateErrType interface {
			Broken() bool
		}
		switch err.(type) {
		//Checking for error behaviour this covers all the errors related to certificate

		case certificateErrType:
			println("Requesting new Certificate for you...")
			pub = marshalauthkey(dpub)
		default:
			println("Generating new credentials for you...")
			//now will generate new rsa key pair
			prv, pub = generateNewKeyPair(pass)
			writePairToDB(pub, prv, user.Uid) //Persisting to db

		}
		cert := new(remote.CertificateResponse)
		cli := getBackendClient()
		if cli == nil {
			println("Couldn't request  certificate authentication")
			return
		}
		err = cli.Call("Backend.IssueCertificate", &remote.CertificateRequest{User: *user, PublicKey: pub}, cert)
		if err != nil {
			print(COLOR_RED)
			println("Couldn't issue cetificate on your behalf due to\t", err.Error())
			resetConsoleColor()
		}
		marshaled := cert.Bytes

		print(COLOR_GREEN)
		println("Certicate Generated\nCertificate Finger Print")
		print(COLOR_YELLOW)
		certpub, _, _, _, _ := parseauthkey(cert.Bytes)
		fmt.Printf("%s\n", ssh.FingerprintLegacyMD5(certpub))
		resetConsoleColor()
		//Writting certificate to db
		localdb.Put([]byte(fmt.Sprint(keycertkey, user.Uid)), marshaled, nil)
	} else {
		print(COLOR_GREEN)
		println("Found Local Credential!!")
		resetConsoleColor()
	}
}

//validatelocalCertificate validates credentials locally,
// generally it checks expiry date & principals etc
func validatelocalCertificate(user *remote.User, pass []byte) (dpub ssh.PublicKey, err error) {

	hpr, hpub, hcert, pki := searchForPKICredentials(user.Uid)
	if hpr && hpub {
		//We have all the credentials

		//Let's Parse the private key
		signer, err := ssh.ParsePrivateKeyWithPassphrase(pki.prv, pass)
		if err != nil {
			err = errors.New("Couldn't decrypt your private due to\t" + err.Error())
			return nil, err
		}
		_, _, _, _, err = parseauthkey(pki.pub)
		if err != nil {
			err = errors.New("Couldn't parse your public due to\t" + err.Error())
			return nil, err
		}
		dpub := signer.PublicKey() //Public key corrosponding to private key
		if !bytes.Equal(marshalauthkey(dpub), pki.pub) {
			//Private key and public key pair match
			err = errors.New("Invalid Public Key")
			return nil, err
		}
		if hcert {
			certp, _, _, _, _ := parseauthkey(pki.cert)
			cert := certp.(*ssh.Certificate)
			vb := cert.ValidBefore
			now := time.Now().Unix()
			if vb < uint64(now) {
				//Not Valid
				//handle this
				err = certificateErr{Reason: "Certificate Expired!!", expired: true}
				return dpub, err
			}
			if !bytes.Equal(marshalauthkey(cert.Key), pki.pub) {
				//Certificate and key pair matched
				err = certificateErr{Reason: "Invalid Certificate", broken: true}
				return dpub, err
			}
		} else {
			//We have valid key pair but doesn't have certificate
			return dpub, certificateErr{broken: true, Reason: "Certificate not found"}
		}
	} else {
		err = errors.New("Broken Key Pair found")
	}
	return
}

type certificateErr struct {
	Reason  string
	broken  bool
	expired bool
}

func (c certificateErr) Error() string {
	return c.Reason
}
func (c certificateErr) Broken() bool {
	return c.broken
}
func (c certificateErr) Expired() bool {
	return c.expired
}
