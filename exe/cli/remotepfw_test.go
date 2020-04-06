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
	"github.com/devansh42/shree/remote"
)

func initTestEnvironment(t *testing.T) {
	os.Setenv(SSH_HOST, "localhost")
	os.Setenv(SSH_PORT, "9500")
	initApp()
	startDemoRPCServer(t)    //Test server for fetching required certificates
	go startDemoSSHServer(t) //Test ssh server

}

func startTestHttpServerForPortRange(t *testing.T) {

	for i := 0; i < 15; i++ {
		exe.StartTestHttpServer(3000 + i)

	}
	t.Log("All http server running")
}

func diconnectRemoteTunnelRange(t *testing.T) {
	for i := 0; i < 15; i++ {
		disconnectRemoteForwardedPort(sprint(i + 3000))
		t.Log("Disconnected remote tunnel for port ", 3000+i)
	}
}

func remotePortFwdForTest(t *testing.T) {
	var mapping = make(map[int]string)
	//Forwarding remote port
	for i := 0; i < 15; i++ {
		src := forwardRemotePort("tcp", 3000+i, testpasswd)
		//Forwarding remote port
		t.Logf("Remote tunnel established for remote:%s \t->\t %d  ", src, 3000+i)
		mapping[3000+i] = src
	}

	//Pinging connections
	for k, v := range mapping {
		x := sprint("http://localhost:", v, "/", k)
		t.Log(x)
		resp, err := http.Get(x)
		if err != nil {
			t.Log("Couldn't reach to server at port ", v, " mapped with port ", k, "on local host")
			continue
		}
		t.Log("Status Recived from ", v, " is ", resp.StatusCode)

	}
}

//TestRemotePortForwardingWithoutCredentials tests the situation when remote port fwd doesn't works
func TestRemotePortForwardingWithoutCredentials(t *testing.T) {
	initTestEnvironment(t)
	defer cleanup()
	//Starting test servers
	startTestHttpServerForPortRange(t)

	//Forwarding port range
	remotePortFwdForTest(t)

	//Listing the remote forwarded port
	listConnectedRemoteTunnel()
	//Closing remote tunnel

}

//TestRemotePortForwardingWithCredentials test the situation,
// when remote port doesn't work due to invalid or broken credentials
func TestRemotePortForwardingWithCredentials(t *testing.T) {
	initTestEnvironment(t)
	defer cleanup()
	defer cleanupDemoCredentials()

	//Creating fake user
	currentUser = new(remote.User)
	currentUser.Username = "devansh42"
	currentUser.Uid = 1

	//generating demo credentials
	generateAndPersistCredentialsForTest(currentUser, testpasswd, t)

	//Starting test servers
	startTestHttpServerForPortRange(t)

	//Forwarding port range
	remotePortFwdForTest(t)

	//Listing the remote forwarded port
	listConnectedRemoteTunnel()
	//Closing remote tunnel

}

type ppt struct {
	CAddr string
	CPort uint32
	DAddr string
	DPort uint32
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

		fb, _ := ioutil.ReadFile("../../keys/ca_user_key.pub")

		pk, _, _, _, _ := parseauthkey(fb)

		o := bytes.Equal(marshalauthkey(pk), marshalauthkey(auth))

		return o
	}
	config.PublicKeyCallback = callmaker(certChecker)
	fb, _ := ioutil.ReadFile("../../keys/id_host")
	signer, _ := ssh.ParsePrivateKey(fb)
	bcert, _ := ioutil.ReadFile("../../keys/id_host-cert.pub")
	pert, _, _, _, err := parseauthkey(bcert)
	fatalTestErr(t, err)
	realsigner, err := ssh.NewCertSigner(pert.(*ssh.Certificate), signer)
	fatalTestErr(t, err)
	config.AddHostKey(realsigner) //Private key for ssh server

	for {
		c, err := listener.Accept()
		fatalTestErr(t, err)
		t.Log("New Connection")
		go func(c net.Conn) {
			sconn, newch, reqch, err := ssh.NewServerConn(c, config)
			if err != nil {
				t.Fatal("From Server ", err)
			}

			go func(ch <-chan *ssh.Request) {
				//handling incomming requests

				for x := range ch {
					switch x.Type {
					case "tcpip-forward": //Handling tcp ip forwarding
						var p struct {
							Addr string
							Port uint32
						}
						ssh.Unmarshal(x.Payload, &p)
						t.Log(p)
						l, err := net.Listen("tcp", exe.JoinHost(p.Addr, p.Port))
						fatalTestErr(t, err)
						_, xport, _ := net.SplitHostPort(l.Addr().String())
						iport, _ := strconv.Atoi(xport)

						if x.WantReply {
							var xp struct {
								Port uint32
							}
							xp.Port = uint32(iport)

							b := ssh.Marshal(&xp)
							x.Reply(true, b)
							if err := ssh.Unmarshal(b, &xp); err != nil {
								t.Log("From server", err)
							}
						}
						go func(l net.Listener) {
							for {
								inc, err := l.Accept()

								fatalTestErr(t, err)
								raddr := inc.RemoteAddr().String()
								host, sport, _ := net.SplitHostPort(raddr)
								port, _ := strconv.Atoi(sport)
								pp := ppt{p.Addr, uint32(iport), host, uint32(port)}
								b := ssh.Marshal(&pp)

								sch, rch, err := sconn.OpenChannel("forwarded-tcpip", b)
								if err != nil {
									//handle error
									log.Print("couldn't open channel ", err.Error())
								}
								go ssh.DiscardRequests(rch)
								go exe.HandleConnectionIO(inc, sch)
							}
						}(l)
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
		}(c)
	} //For ends
}

func fatalTestErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

var testpasswd = []byte("hello1234")
