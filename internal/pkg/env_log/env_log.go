package env_log

import (
	"io"
	"log"
	"os"
)

var logFile *os.File
var err error
var mw io.Writer

func EnvLog() error {
	if logFile != nil {
		return nil
	}
	logFilePath := os.Args[0] + ".log"
	logFile, err = os.Create(logFilePath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	mw = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(mw)
	log.Printf("log at %v", logFilePath)
	return err
}

func GetFile() io.Writer {
	return mw
}
