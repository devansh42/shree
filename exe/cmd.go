package exe

import (
	"fmt"
	"io"
	"log"
	"os"
)

var sprint = fmt.Sprint
var ps = string(os.PathSeparator)

func SetLogFile(logDir, appDir string) io.Writer {
	var file io.Writer
	file, err := os.OpenFile(sprint(logDir, ps, "logs.log"), os.O_APPEND, 0600)
	if err != nil {
		log.Print("Couldn't Found logfile, writting logs to current file")
		fname := appDir
		_, err = os.Stat(fname)
		if err != nil {
			os.MkdirAll(fname, 0700)
			//making  dir if not exists
		}
		fname = sprint(fname, ps, "logs.log")
		file, err = os.OpenFile(fname, os.O_APPEND, 0600)
	}
	//Logs are written to stdout and the log file
	wr := io.MultiWriter(os.Stdout, file)
	return wr
}
