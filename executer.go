package shree

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh/terminal"

	cli "github.com/urfave/cli/v2"
)

//This file contains definitions of executor programs command wise

func signIn(c *cli.Context) error {
	var user user
	var print = fmt.Print
	var println = fmt.Println
	print(COLOR_BLUE)
	print("Hey! please enter your ")
	println("Email:\t")
	fmt.Scan(&user.Email)
	println("Password:\t")
	pass, _ := terminal.ReadPassword(1)
	println("Authenticating...")
	user.setPassword(pass)

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
