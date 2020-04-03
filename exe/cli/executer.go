package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/devansh42/shree/remote"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

//This file contains definitions of executor programs command wise

func signIn(c *cli.Context) error {
	var user remote.User
	var print = fmt.Print
	var println = fmt.Println
	print(COLOR_BLUE)
	print("Hey!! please enter your ")
	print("Email:\t")
	resetConsoleColor()
	fmt.Scan(&user.Email)
	print(COLOR_BLUE)
	print("Username\t")
	resetConsoleColor()
	fmt.Scan(&user.Username)
	print(COLOR_BLUE)
	print("Password:\t")
	resetConsoleColor()
	pass, _ := terminal.ReadPassword(1) //Reading password from stdin
	user.Password = hash(pass)
	println("Authenticating...")
	cli := getBackendClient()
	if cli == nil {
		return errors.New("Couldn't reach to backend server")
	}

	isnew, err := authUser(&user, cli)
	if err != nil {
		print(COLOR_RED)
		print(err.Error())
		resetConsoleColor()
		return err //Exiting from current session
	}
	print(COLOR_GREEN)
	if isnew {
		println("New Account Created!!")
	} else {
		println("Authentication Completed !! \nHello, ", user.Username)
	}
	println("Initializing Access Key Managment....")
	kpass := askForPassword()
	manageCertificate(&user, kpass)
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

//listLocalTunnels  lists local tunnels
func listLocalTunnels(c *cli.Context) error {
	listConnectedLocalTunnel()
	return nil
}

//disconnectLocalTunnel, disconnects local tunnel on local machine previously initiated by shree
func disconnectLocalTunnel(c *cli.Context) error {
	port := c.Uint("port")
	if port == 0 {
		return errors.New("Please provide source port")
	}
	disconnectLocalyForwardedPort(int(port))
	return nil
}

//exposeRemoteTunnel, connects local port to local port
func exposeRemoteTunnel(c *cli.Context) error {
	expose := c.Uint("expose")
	println("Trying to Expose ", expose, " .......")
	//Exposes port
	//0 src ports specifies any port
	forwardRemotePort("tcp", 0, int(expose))
	return nil
}

//listRemoteTunnels, lists connected local remote tunnels
func listRemoteTunnels(c *cli.Context) error {
	listConnectedRemoteTunnel()
	return nil
}

//disconnectRemoteTunnel turns the remote tunnel off
func disconnectRemoteTunnel(c *cli.Context) error {
	port := c.Uint("dest")
	if port == 0 {
		return errors.New("Please provide destination port on local machine")
	}
	disconnectRemoteForwardedPort(sprint(port))
	return nil
}

//exitApp exits apps
func exitApp(c *cli.Context) error {
	exitGracefully()
	println("Exiting...")
	<-time.After(time.Second)
	os.Exit(0)
	return nil
}

func printHelp(c *cli.Context) error {

	return nil
}

//exitGracefully gracefully exits app after stopping all the Sockets
func exitGracefully() {
	for _, v := range remoteForwardedPorts { //Exiting remote forwardings
		v.Closer <- true
	}
	for _, v := range locallyForwardedPort { //Exiting local forwardings
		v.Closer <- true
	}
}
