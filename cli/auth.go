package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/devansh42/shree/remote"
)

const keyuser = "user"

//authUser authenticates user it creates if doesn't exists
func authUser(user *remote.User) (bool, error) {
	cli := getBackendClient()
	resp := new(remote.User)
	err := cli.Call("Backend.Auth", user, resp)
	if err != nil {
		switch err.Error() {
		case "401":
			return false, errors.New("Invalid Username/Password")
		default:
			return false, errors.New("Unknown Error:\t" + err.Error())
		}
	}
	//Lets overwrite the user
	b, _ := json.Marshal(resp)
	localdb.Put([]byte(keyuser), b, nil)
	return resp.IsNew, nil
}

func manageCertificate(user *remote.User) {
	println("Please enter your Password for initalizing Access Key Managment")
	pass, _ := terminal.ReadPassword(1)
	pubkey, err := validatelocalCertificate(user, pass)
	if err != nil {
		print(COLOR_RED)
		println(err.Error())
		resetConsoleColor()
		switch err.(type) {
		//Checking for error behaviour this covers all the errors related to certificate

		case interface {
			Broken() bool
		}:
			println("Requesting new Certificate for you...")

		default:
			println("Generating new credentials for you...")
			//now will generate new rsa key pair
			prv, pub := generateNewKeyPair(pass)
			writePairToDB(pub, prv, user.Uid) //Persisting to db
			pubkey, _, _, _, err = parseauthkey(pub)
			if err != nil {
				println("Broken public key\t", err.Error())
				println("Please try again or report the problem it if problem persists.")
				return
			}
		}
		cert := new(ssh.Certificate)
		err = getBackendClient().Call("Backend.IssueCertificate", &remote.CertificateRequest{User: *user, PublicKey: pubkey}, cert)
		if err != nil {
			print(COLOR_RED)
			println("Couldn't issue cetificate on your behalf due to\t", err.Error())
			resetConsoleColor()
		}
		marshaled := ssh.MarshalAuthorizedKey(cert)
		print(COLOR_GREEN)
		println("Certicate Generated\nHash")
		print(COLOR_YELLOW)
		fmt.Printf("\n%x", hash(marshaled))
		resetConsoleColor()
		//Writting certificate to db
		localdb.Put([]byte(fmt.Sprint(keycertkey, user.Uid)), marshaled, nil)
	}
}

func validatelocalCertificate(user *remote.User, pass []byte) (dpub ssh.PublicKey, err error) {

	hpr, hpub, hcert, pki := searchForPKICredentials(user.Uid)
	if hpr == hpub == hcert == true {
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

		dpub, _ = ssh.NewPublicKey(signer)
		if !bytes.Equal(marshalauthkey(dpub), pki.pub) {
			//Private key and public key pair match
			err = errors.New("Invalid Public Key")
			return nil, err
		}

		certp, _, _, _, _ := parseauthkey(pki.cert)
		cert := certp.(*ssh.Certificate)
		vb := cert.ValidBefore
		now := time.Now().Unix()
		if vb < uint64(now) {
			//Not Valid
			//handle this
			err = certificateErr{Reason: "Certificate Expired!!", expired: true}
			return nil, err
		}
		if !bytes.Equal(marshalauthkey(cert.Key), pki.pub) {
			//Certificate and key pair matched
			err = certificateErr{Reason: "Invalid Certificate", broken: true}
			return nil, err
		}
	} else {
		err = errors.New("Broken Key Pair not found")
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
