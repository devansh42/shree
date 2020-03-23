package shree

import (
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
