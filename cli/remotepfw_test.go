package main

import (
	"net"
	"testing"
	"time"
)

func TestRemotePortForwarding(t *testing.T) {
	forwardRemotePort("tcp", 6000, 8080)
	<-time.After(time.Minute * 10) //Waiting for two minutes
}

func TestRemotePortForwarding1(t *testing.T) {
	forwardRemotePort("tcp", 7800, 8080)
	<-time.After(time.Minute * 10) //Waiting for two minutes
}

func TestPortZeroListen(t *testing.T) {

	x := 0
	for {
		listen, err := net.Listen("tcp", net.JoinHostPort("", "0"))

		if err != nil {
			t.Log("Couldn't listen ", err.Error())
			t.Log(x)
			return
		}
		defer listen.Close()
		t.Log(listen.Addr())
		x++
	}
	t.Log(x)
}
