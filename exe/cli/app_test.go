package main

import (
	"os"
	"testing"
)

func TestHomeDirectory(t *testing.T) {
	h, err := os.UserHomeDir()
	if err != nil {
		//handle error
	}
	t.Log(h)
}

func TestDirExistance(t *testing.T) {
	h, err := os.UserHomeDir()
	f := sprint(h, "/", "demodir")
	_, err = os.Stat(f)
	if err != nil {
		t.Log(err, os.IsNotExist(err))
	}

}

const testingHttpServerPort = 9090
