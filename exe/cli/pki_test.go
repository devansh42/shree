package main

import (
	"io/ioutil"
	"testing"

	"github.com/devansh42/shree/remote"

	"golang.org/x/crypto/ssh"
)

func TestGenerateNewKeyPair(t *testing.T) {
	prv, pub := generateNewKeyPair([]byte("hello123"))
	t.Log("Private Key\n", string(prv))
	t.Log("Public Key\n", string(pub))
	//Lets decode private key
	_, err := ssh.ParsePrivateKeyWithPassphrase(prv, []byte("hello123"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Decoded Key Pair")

}

func TestWriteKeys(t *testing.T) {
	defer cleanup()
	initApp()
	for x := 0; x < 5; x++ {
		prv, pub := generateNewKeyPair([]byte("hello123"))
		writePairToDB(pub, prv, int64(x))
		//Lets read is back
		pu, err := localdb.Get([]byte(sprint(keypubkey, x)), nil)
		if err != nil {
			t.Fatal(err)
		}
		pr, err := localdb.Get([]byte(sprint(keyprvkey, x)), nil)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Private Key Reader\n", string(pr))
		t.Log("Public Key Reader\n", string(pu))
	}
}

//TestFetchAndPersistCAPublicHostCertificate fetches service from ca public
func TestFetchAndPersistCAPublicHostCertificate(t *testing.T) {
	initApp()
	defer cleanup()
	startDemoRPCServer(t)
	cert, err := fetchServerCertificateAndPersist()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Fetched Certificate\n", string(marshalauthkey(cert)))
	t.Log("Fetching from Localdb")
	b, _ := localdb.Get([]byte(keyservercertificate), nil)
	t.Log(string(b))
	//Deleting Certificates
	localdb.Delete([]byte(keyservercertificate), nil)
}

func TestGetCAPublicHostCertificateFromLocaldbrOrFetch(t *testing) {

}

func (b *Backend) GetCAPublicCertificate(req *remote.CertificateRequest, cert *remote.CertificateResponse) error {
	f, err := ioutil.ReadFile("./id_host-cert.pub")
	if err != nil {
		return err
	}
	cert.Bytes = f
	return nil
}
