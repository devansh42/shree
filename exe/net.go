package exe

import (
	"fmt"
	"io"
	"net"
)

type Closerch chan bool

type Forwardedport struct {
	DestPort string
	SrcPort  string
	Listener net.Listener
	Closer   Closerch
}
type LocalForwardedport struct {
	Forwardedport
}

type RemoteForwardedport struct {
	Forwardedport
}

//This file contains socket realted primitives used
type socket chan bool
type SocketCollection struct {
	socketmap map[string]socket
}

func (s *SocketCollection) Close(name string) (closed bool) {

	for k, v := range s.socketmap {
		if k == name {
			v <- true
			close(v)
			closed = true
			delete(s.socketmap, k)
		}
	}

	return

}

type closeable interface {
	Close() error
}

func (s *SocketCollection) Add(name string, conn closeable) {
	var closer = make(socket)
	s.socketmap[name] = closer
	go func(c closeable, cl *socket) {
		defer c.Close()
		<-closer //Waiting for closing signal
	}(conn, &closer)

}

func (s *SocketCollection) Have(name string) bool {
	for k, _ := range s.socketmap {
		if k == name {
			return true
		}
	}
	return false
}

func JoinHost(host string, port interface{}) string {
	return net.JoinHostPort(host, fmt.Sprint(port))
}

//HandleForwardedListener handles forwarded listener and manages its life cycle also
//It recieves an closing singal to break the infinte loop of connection acceptance
func HandleForwardedListener(conn *Forwardedport) {
	listener := conn.Listener
	defer listener.Close()

	var closeMe = new(bool)

	go func(closeMe *bool) {
		<-conn.Closer //Waiting for close signal
		*closeMe = true
	}(closeMe)

	//Running forever
	for {
		if *closeMe {
			break
		}

		//Making a connection for relaying data to local port
		relayConn, err := net.Dial("tcp", JoinHost("", conn.DestPort))
		if err != nil {
			//Couldn't handle this connection
			continue
		}
		acceptedConn, err := listener.Accept() //Accepting connection at remote port
		if err != nil {
			//handle it
			continue
		}
		go HandleConnectionIO(acceptedConn, relayConn)

	}
}

//handleConnectionIO handle i/o b/w relayed connection and accepted connections
func HandleConnectionIO(acceptedConn, relayConn io.ReadWriteCloser) {
	defer acceptedConn.Close()
	defer relayConn.Close()
	closer := make(chan bool)
	go func() {
		io.Copy(relayConn, acceptedConn)
		closer <- true
	}()

	go func() {
		io.Copy(acceptedConn, relayConn)
		closer <- true
	}()
	<-closer //Whenever it hears a signal it closes both sides of remote connection

}
