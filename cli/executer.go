package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/devansh42/shree/remote"

	"golang.org/x/crypto/ssh/terminal"

	cli "github.com/urfave/cli/v2"
)

//This file contains definitions of executor programs command wise

func signIn(c *cli.Context) error {
	var user remote.User
	var print = fmt.Print
	var println = fmt.Println
	print(COLOR_BLUE)
	print("Hey!! please enter your ")
	println("Email:\t")
	fmt.Scan(&user.Email)
	println("Username\t")
	fmt.Scan(&user.Username)
	println("Password:\t")
	pass, _ := terminal.ReadPassword(1)
	user.Password = hash(pass)
	println("Authenticating...")
	isnew, err := authUser(&user)
	if err != nil {
		print(COLOR_RED)
		print(err.Error())

	}
	print(COLOR_GREEN)
	if isnew {
		println("New Account Created!!")
	} else {
		println("Authentication Completed !! \nHello, ", user.Username)
	}
	return nil
}

func signOut(c *cli.Context) error {
	return nil
}

//connectLocalTunnel, connects 2 local ports on local machine
func connectLocalTunnel(c *cli.Context) error {
	src := c.Uint("src")
	dest := c.Uint("dest")
	if src == 0 || dest == 0 {
		log.Fatal("src/dest port not specified")
	}
	protocol := c.String("protocol")
	if protocol == "" {
		protocol = "tcp" //default protocol is tcp
	}
	forwardLocalPort(protocol, int(src), int(dest))
	return nil
}

//disconnectLocalTunnel, disconnects local tunnel on local machine previously initiated by shree
func disconnectLocalTunnel(c *cli.Context) error {
	port := c.Uint("src")
	if port == 0 {
		return errors.New("Please provide source port")
	}
	disconnectLocalyForwardedPort(int(port))
	return nil
}

//   ruchiyadav1821
// Labneighbours@netflix
