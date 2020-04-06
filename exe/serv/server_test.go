package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"testing"

	"github.com/devansh42/shree/exe"

	"golang.org/x/crypto/ssh"

	"github.com/devansh42/shree/remote"
)

//This file contains various server for testing purposes

func gethostCallback(t *testing.T) ssh.HostKeyCallback {
	f := "../../keys/ca_host_key.pub"
	b, err := ioutil.ReadFile(f)
	fatalErr(t, err)
	k, _, _, _, err := ssh.ParseAuthorizedKey(b)
	fatalErr(t, err)
	certchecker := new(ssh.CertChecker)
	certchecker.IsHostAuthority = func(auth ssh.PublicKey, address string) bool {
		b := ssh.MarshalAuthorizedKey(auth)
		c := ssh.MarshalAuthorizedKey(k)
		return bytes.Equal(b, c)
	}
	return certchecker.CheckHostKey
}

func getSigner(t *testing.T) ssh.AuthMethod {
	fp := "../../keys/id_user"
	b, err := ioutil.ReadFile(fp)
	s, err := ssh.ParsePrivateKey(b)
	fatalErr(t, err)
	fc := "../../keys/id_user-cert.pub"
	b, err = ioutil.ReadFile(fc)
	fatalErr(t, err)
	ck, _, _, _, err := ssh.ParseAuthorizedKey(b)
	fatalErr(t, err)
	cert := ck.(*ssh.Certificate)
	p, err := ssh.NewCertSigner(cert, s)
	fatalErr(t, err)
	return ssh.PublicKeys(p)
}

func getSSHClientConfig(t *testing.T) *ssh.ClientConfig {
	config := new(ssh.ClientConfig)
	config.User = "devansh42"

	config.HostKeyCallback = gethostCallback(t)

	config.Auth = []ssh.AuthMethod{getSigner(t)}
	return config
}

func setupTestEnvironment(t *testing.T) {
	os.Setenv(SHREE_SSH_PORT, sprint(7500))
	os.Setenv(SHREE_SSH_PRIVATE_KEY, "../../keys/id_host")
	os.Setenv(SHREE_SSH_PUBLIC_KEY, "../../keys/id_host.pub")
	os.Setenv(SHREE_BACKEND_ADDR, "localhost:6500") //Address of rpc server
	startTestRPCServer(t)
	initApp()
	go initServer() //Starts the actual ssh server

}

func startTestRPCServer(t *testing.T) {
	//Backend
	l, err := net.Listen("tcp", os.Getenv(SHREE_BACKEND_ADDR))
	fatalErr(t, err)
	rpc.Register(&Backend{}) //Registering backend
	go rpc.Accept(l)
	t.Log("Accepting backend connections at ", os.Getenv(SHREE_BACKEND_ADDR))
}

func fatalErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

type Backend struct{}

func (b *Backend) GetCAUserPublicCertificate(req *remote.CertificateRequest, resp *remote.CertificateResponse) error {
	f, err := ioutil.ReadFile("../../keys/ca_user_key.pub")
	if err != nil {
		log.Print(err)
		return err
	}
	k, _, _, _, err := ssh.ParseAuthorizedKey(f)
	if err != nil {
		log.Println(err)
		return err
	}
	resp.Bytes = ssh.MarshalAuthorizedKey(k)
	return nil
}

func (b *Backend) IssueHostCertificate(req *remote.HostCertificateRequest, resp *remote.CertificateResponse) (err error) {
	f, err := rF("../../keys/id_host-cert.pub")

	if err != nil {
		return err
	}

	resp.Bytes = f
	return
}

func TestSSHServer(t *testing.T) {
	setupTestEnvironment(t)
	//Let's make a pseudo tcp connections

	for i := 0; i < 15; i++ {
		exe.StartTestHttpServer(3000 + i)
		t.Log("Http server is listening at ", 3000+i)
	}
	cli, err := ssh.Dial("tcp", exe.JoinHost("localhost", os.Getenv(SHREE_SSH_PORT)), getSSHClientConfig(t))
	if err != nil {
		t.Fatal("Handshake failed due to ", err.Error())
	}
	var m = make(map[string]string)
	for i := 0; i < 15; i++ {
		l, err := cli.Listen("tcp", exe.JoinHost("0.0.0.0", 0))
		fatalErr(t, err)
		f := new(exe.Forwardedport)
		f.Listener = l
		f.Closer = make(exe.Closerch)
		f.DestPort = sprint(3000 + i)
		_, p, _ := net.SplitHostPort(l.Addr().String())
		f.SrcPort = p
		go exe.HandleForwardedListener(f) //Handles listening
		t.Log("port fwd established ", f.SrcPort, "\t->\t", f.DestPort)
		m[f.DestPort] = f.SrcPort
	}
	for k, v := range m {
		//ping remotely forwarde ports
		res, err := http.Get(sprint("http://localhost:", v, "/", k))
		if err != nil {
			t.Log("Couldn't make to port ", v)
			continue
		}
		t.Log("Success! Http request from port ", v, " with code ", res.StatusCode)
	}

}
