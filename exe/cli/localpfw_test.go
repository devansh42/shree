package main

import (
	"net/http"
	"testing"
)

//TestLocalPortForwarding tests all operations on local port forwarding scenenarios
func TestLocalPortForwarding(t *testing.T) {

	initApp()
	for i := 9000; i < 9015; i++ {
		startTestHttpServer(i)

	}
	t.Log("Started Test servers")
	for i := 3000; i < 3015; i++ {
		err := forwardLocalPort("tcp", i, 9000+(i-3000))
		if err != nil {
			t.Error(err)
			continue
		}
		t.Log("Local Portfwded ", i, "\t->\t", 9000+(i-3000))

	}
	t.Log("Started Local portfwded")
	//lets starts for http responses
	for i := 3000; i < 3015; i++ {
		t.Log("Hearing from server at ", i)
		r, err := http.Get("http://localhost:" + sprint(i) + "/" + sprint(9000+(i-3000)))
		if err != nil {
			t.Error(err)
			continue
		}

		t.Log(r.StatusCode)

	}

	t.Log("Now listing local ports")
	listConnectedLocalTunnel()
	t.Log("Disconnecting locally forwarded ports")
	for i := 3000; i < 3015; i++ {
		disconnectLocalyForwardedPort(i)
		t.Log("Disconnected local tunnel at port ", i)
	}
	listConnectedLocalTunnel()
}
