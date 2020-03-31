package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	cli "github.com/urfave/cli/v2"
)

const keyLocalpfw = "localpfw"

var (
	locallyForwardedPort []*Forwardedport
)

//forwardLocalPort, forwards local port src->dest
//This listens connection from src and relay them to dest
//protocol defines protocol to be used tcp or udp
func forwardLocalPort(protocol string, src, dest int) (err error) {
	listener, err := net.Listen(protocol, joinHost("", src))
	if err != nil {
		log.Print("Couldn't connect ports due to : ", err.Error())
		return err
	}

	lfp := &Forwardedport{fmt.Sprint(dest), fmt.Sprint(src), listener}

	locallyForwardedPort = append(locallyForwardedPort, lfp)

	//Printing status
	print(COLOR_GREEN_UNDERLINED)
	println("Successfully Port Forwarding Established")
	println(listener.Addr().String(), "\t->\t", joinHost("", lfp.DestPort))
	resetConsoleColor()

	handleForwardedListener(lfp)

	return
}
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

//disconnectLocalyForwardedPort, disconnects localy forwarded port
func disconnectLocalyForwardedPort(src int) {
	closed := socketCollection.close(strings.Join([]string{"flp", strconv.Itoa(src)}, ":"))
	if closed {
		//Success
		fmt.Print(COLOR_GREEN)
		fmt.Print("Connection closed from port ", src)
		resetConsoleColor()
		//Removing data from local db
		b, err := localdb.Get([]byte(keyLocalpfw), nil)
		if err != nil {
			//key doesn't exists unlikly to happen

		}
		var ar [][2]int
		json.Unmarshal(b, &ar)
		if len(ar) > 1 {
			var bar [][2]int
			for _, v := range ar {
				s := v[0]
				if s == src {
					//Need to be removed
					continue
				}
				bar = append(bar, v)
			}
			b, _ := json.Marshal(bar)
			localdb.Put([]byte(keyLocalpfw), b, nil) //Putting back to db
		} else {
			localdb.Delete([]byte(keyLocalpfw), nil)

		}

	} else {
		fmt.Print(COLOR_RED)
		fmt.Print("Couldn't close connection from port, either it belongs to another process or its already closed")
		resetConsoleColor()
	}
}

//This function lists local connected tunnels
func listConnectedLocalTunnel(c *cli.Context) error {
	fmt.Println("Local connected ports")
	b, err := localdb.Get([]byte("localpfw"), nil)
	if err != nil {
		//This means no ports are connected
		fmt.Println("No port(s) connections available")
	} else {
		var ar [][2]int
		json.Unmarshal(b, &ar)
		for i, v := range ar {
			fmt.Print(COLOR_GREEN)
			fmt.Printf("\n%d\t%d\t->%d", i+1, v[0], v[1])
		}
		fmt.Print(COLOR_RESET)
	}
	return nil
}
