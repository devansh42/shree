package shree

import (
	"net"
	"net/http"
	"testing"
)

//This file contains various server for testing purposes

func TestSimpleHttpServer(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from, SimpleHTTP Server"))
		w.WriteHeader(200)
	})
	http.ListenAndServe(":8000", nil)
}

func TestServerImplementation(t *testing.T) {
	initServer()
	//Initialize the server
}

func TestLocalVsRemoteAddr(t *testing.T) {
	c, _ := net.Dial("tcp", "www.google.com:80")
	t.Log(c.RemoteAddr().String())
	t.Log(c.LocalAddr().String())

}
