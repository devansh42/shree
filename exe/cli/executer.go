package main

import (
	"errors"
	"fmt"
	"log"
	"net"
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
	println("This feature is under development, just use sign me in, to change user account")
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
	expose := c.Uint("dest")
	println("Trying to Expose ", expose, " .......")
	//Exposes port
	//0 src ports specifies any port
	bpass := askForPassword()
	forwardRemotePort("tcp", int(expose), bpass)
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
	<-time.After(time.Millisecond * 200) //200ms
	os.Exit(0)
	return nil
}

func whoAmI(c *cli.Context) error {
	getCurrentUserInfo()
	return nil
}

func getProps(c *cli.Context) (err error) {
	for _, v := range c.FlagNames() {
		switch v {
		case "backend":
			x := c.Bool(v)
			if x {
				println("Backend Addr\t", os.Getenv(SHREE_BACKEND_ADDR))
			}
		case "remote":
			x := c.Bool(v)
			if x {
				println("SSH Addr\t", os.Getenv(SHREE_SSH_ADDR))

			}
		default:
			err = errors.New("Couldn't found properties")

		}
	}
	return
}

func setProps(c *cli.Context) (err error) {
	for _, v := range c.FlagNames() {
		switch v {
		case "backend":
			x := c.String(v)
			if x != "" {
				_, _, err = net.SplitHostPort(x)
				if err != nil {
					return
				}
				os.Setenv(SHREE_BACKEND_ADDR, x)
				setProp(SHREE_BACKEND_ADDR, []byte(x))
			}
		case "remote":
			x := c.String(v)
			if x != "" {
				_, _, err = net.SplitHostPort(x)
				if err != nil {
					return
				}
				os.Setenv(SHREE_SSH_ADDR, x)
				setProp(SHREE_SSH_ADDR, []byte(x))

			}
		default:
			err = errors.New("Couldn't found properties")
		}
	}
	return
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
