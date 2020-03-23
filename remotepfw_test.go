package shree

import (
	"testing"
	"time"
)

func TestRemotePortForwarding(t *testing.T) {
	forwardRemotePort("tcp", 5000, 8000)
	<-time.After(time.Minute * 10) //Waiting for two minutes
}
