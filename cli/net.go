package main

//This file contains socket realted primitives used

type socket chan bool
type socketcollection struct {
	socketmap map[string]socket
}

func (s *socketcollection) close(name string) (closed bool) {

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

func (s *socketcollection) add(name string, conn closeable) {
	var closer = make(socket)
	s.socketmap[name] = closer
	go func(c closeable, cl *socket) {
		defer c.Close()
		<-closer //Waiting for closing signal
	}(conn, &closer)

}

func (s *socketcollection) have(name string) bool {
	for k, _ := range s.socketmap {
		if k == name {
			return true
		}
	}
	return false
}
