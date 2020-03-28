package main

import (
	"testing"
	"time"
)

func TestLocalPortForwarding(t *testing.T) {

	forwardLocalPort("tcp", 3000, 6379)
	<-time.After(time.Second * 30)
}
