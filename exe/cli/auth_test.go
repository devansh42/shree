package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"testing"

	"github.com/devansh42/shree/remote"
)

//File contains

//TestAuthSignupUser tests user signup process
func TestAuthSignupUser(t *testing.T) {
	initApp()
	defer cleanup()
	startDemoRPCServer(t) //Starting demo rpc server
	cli := getBackendClient()
	defer cli.Close()
	user := &remote.User{Email: "xyz", IsNew: true, Uid: 1, Username: "xyz", Password: hash([]byte("password"))}
	isnew, err := authUser(user, cli)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(currentUser, isnew) //As after login it sets current user
}

//Export for backend related server
type Backend struct{}

//Auth pretends the auth method for the service
func (b *Backend) Auth(user, resp *remote.User) error {
	log.Print("User recived ", user)
	//x := remote.User{user.Email, user.IsNew, user.Password, user.Uid, user.Username}
	resp.Email = user.Email
	resp.IsNew = user.IsNew
	resp.Password = user.Password
	resp.Uid = user.Uid
	resp.Username = user.Username

	if user.IsNew { //For sake of testing is isnew is true it means user wants to test signup

	} else {
	}

	return nil
}

func startDemoRPCServer(t *testing.T) {

	os.Setenv(SHREE_BACKEND_ADDR, ":6500") //Starting demo rpc server, inorder to work with getBackendClient Method
	listener, err := net.Listen("tcp", ":6500")
	if err != nil {
		t.Fatal("Couldn't start server at given port")
	}
	rpc.Register(&Backend{})
	t.Log("RPC server listening on :6500")
	go rpc.Accept(listener)
}
