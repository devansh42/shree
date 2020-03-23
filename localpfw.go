package shree

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	cli "github.com/urfave/cli/v2"
)

const keyLocalpfw = "localpfw"

//forwardLocalPort, forwards local port src->dest
//This listens connection from src and relay them to dest
//protocol defines protocol to be used tcp or udp
func forwardLocalPort(protocol string, src, dest int) (err error) {
	relay, err := net.Listen(protocol, net.JoinHostPort("", strconv.Itoa(src)))
	if err != nil {
		log.Print("Couldn't connect ports due to : ", err.Error())
		return err
	}

	socketCollection.add(strings.Join([]string{"flp", strconv.Itoa(src)}, ":"), relay) //Adding listener

	handleConnetions := func(inc, outc net.Conn) {
		defer outc.Close()
		var cl = make(chan bool)
		go func() {
			_, err := io.Copy(outc, inc)
			if err != nil {
				//Somethings
				log.Print("Problem while reading from connection\t", err.Error())
			}
			cl <- true
		}()
		go func() {
			_, err := io.Copy(outc, inc)
			if err != nil {
				log.Print("Problem while reading from connection\t", err.Error())

				//Somethings
			}
		}()
		<-cl //waiting for end of conversation
	}

	go func() { //to handle connections
		defer relay.Close()

		for {
			inconn, err := relay.Accept()
			if err != nil {
				log.Print("Couldn't accept incoming connection due to : ", err.Error())
			}
			outconn, err := net.Dial(protocol, net.JoinHostPort("", strconv.Itoa(dest)))
			if err != nil {
				log.Print("Couldn't dial to dest port due to : ", err.Error())
			}

			//Below two lines relay the connection
			handleConnetions(inconn, outconn)
		}

	}()

	//Let's this action to the localdb for logging and ux purposes
	v, err := localdb.Get([]byte(keyLocalpfw), nil)
	var ar = make([][2]int, 1)

	if err != nil {
		//Key doesn't exists

		ar[0] = [2]int{src, dest}
	} else {
		json.Unmarshal(v, &ar) //Error check suppressed
		ar = append(ar, [2]int{src, dest})
	}
	b, _ := json.Marshal(ar)
	localdb.Put([]byte(keyLocalpfw), b, nil)

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
