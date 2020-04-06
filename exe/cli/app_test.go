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

const testingHttpServerPort = 9090
