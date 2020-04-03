package main

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/devansh42/shree/exe"
)

func TestHomeDirectory(t *testing.T) {
	h, err := os.UserHomeDir()
	if err != nil {
		//handle error
	}
	t.Log(h)
}

const testingHttpServerPort = 9090

//starts testing http server on given port
func startTestHttpServer(port int) {

	http.HandleFunc("/"+sprint(port), func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world\n"))
		w.Write([]byte("Here is the remote addr\n"))
		w.Write([]byte(r.RemoteAddr))
		w.WriteHeader(200)
	})
	log.Println("Testing server is listening at ", port)
	go http.ListenAndServe(exe.JoinHost("", port), nil)

}
