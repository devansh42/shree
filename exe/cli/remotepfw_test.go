package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/devansh42/shree/exe"
)

func TestRemotePortForwarding(t *testing.T) {
	os.Setenv(SSH_HOST, "localhost")
	os.Setenv(SSH_PORT, "9500")

	initApp()
	startDemoRPCServer(t)
	startDemoSSHServer(t)
	defer cleanup()

	for i := 0; i < 15; i++ {
		startTestHttpServer(3000 + i)

	}
	t.Log("All http server running")
	var mapping = make(map[int]string)
	//Forwarding remote port
	for i := 0; i < 15; i++ {
		src := forwardRemotePort("tcp", 3000+i)
		//Forwarding remote port
		t.Logf("Remote tunnel established for remote:%s \t->\t %d  ", src, 3000+i)
		mapping[3000+i] = src
	}

	//Pinging connections
	for k, v := range mapping {
		resp, err := http.Get(sprint("localhost:", v, "/", k))
		if err != nil {
			t.Log("Couldn't reach to server at port ", v, " mapped with port ", k, "on local host")
			continue
		}
		t.Log("Status Recived from ", v, " is ", resp.StatusCode)

	}

	//Listing the remote forwarded port
	listConnectedRemoteTunnel()
	//Closing remote tunnel

	for i := 0; i < 15; i++ {
		disconnectRemoteForwardedPort(sprint(i + 3000))
		t.Log("Disconnected remote tunnel for port ", 3000+i)
	}
}

type ppt struct {
	CAddr string
	CPort uint
	DAddr string
	DPort uint
}

type callbackfn func(conn ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error)

func callmaker(certcheker *ssh.CertChecker) callbackfn {
	return func(conn ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
		p, e := certcheker.Authenticate(conn, pubKey)
		return p, e
	}
}
func startDemoSSHServer(t *testing.T) {
	listener, err := net.Listen("tcp", exe.JoinHost(os.Getenv(SSH_HOST), os.Getenv(SSH_PORT)))
	fatalTestErr(t, err)
	defer listener.Close()
	t.Log("SSH is listening... ", os.Getenv(SSH_PORT))
	config := new(ssh.ServerConfig)
	certChecker := new(ssh.CertChecker)
	certChecker.IsUserAuthority = func(auth ssh.PublicKey) bool {
		fb, _ := ioutil.ReadFile("./ca_user_key.pub")
		cert := auth.(*ssh.Certificate)
		return bytes.Equal(fb, ssh.MarshalAuthorizedKey(cert.SignatureKey))
	}
	config.PublicKeyCallback = callmaker(certChecker)
	fb, _ := ioutil.ReadFile("./id_host")
	signer, _ := ssh.ParsePrivateKey(fb)
	config.AddHostKey(signer) //Private key of the server
	for {
		c, err := listener.Accept()
		fatalTestErr(t, err)
		sconn, newch, reqch, err := ssh.NewServerConn(c, config)
		fatalTestErr(t, err)

		go func(ch <-chan *ssh.Request) {
			//handling incomming requests

			for x := range ch {
				switch x.Type {
				case "tcpip-forward": //Handling tcp ip forwarding
					var p struct {
						Addr string
						Port uint
					}
					ssh.Unmarshal(x.Payload, &p)
					l, err := net.Listen("tcp", exe.JoinHost(p.Addr, p.Port))
					fatalTestErr(t, err)
					_, xport, _ := net.SplitHostPort(l.Addr().String())
					var xp struct {
						Port uint
					}
					iport, _ := strconv.Atoi(xport)
					xp.Port = uint(iport)
					x.Reply(true, ssh.Marshal(&xp))
					for {
						inc, err := l.Accept()
						fatalTestErr(t, err)
						raddr := inc.RemoteAddr().String()
						host, sport, _ := net.SplitHostPort(raddr)
						port, _ := strconv.Atoi(sport)
						pp := ppt{p.Addr, p.Port, host, uint(port)}
						b := ssh.Marshal(&pp)

						sch, rch, err := sconn.OpenChannel("forwarded-tcpip", b)
						if err != nil {
							//handle error
							log.Print("couldn't open channel ", err.Error())
						}
						go ssh.DiscardRequests(rch)
						go exe.HandleConnectionIO(inc, sch)
					}
				}
			}

		}(reqch)

		go func(ch <-chan ssh.NewChannel) { //Discarding some of channel requests
			for x := range ch {
				switch x.ChannelType() {
				case "session":
					_, y, err := x.Accept()
					fatalTestErr(t, err)
					go ssh.DiscardRequests(y)
				default:
					x.Reject(ssh.Prohibited, "don't have time s")

				}
			}
		}(newch)
	}
}

func fatalTestErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
