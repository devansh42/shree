package main

import (
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
