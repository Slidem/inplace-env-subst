package main

import (
	"log"
	"os"
)

type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(v ...interface{})
}


// Returns the appropriate logger
// if env variable DEBUG is true,additional info will be printed to the console
// if env variable is not set or is false, no output will be printed to the console
func GetLogger() Logger {

	debugEnabled := os.Getenv("DEBUG")
	if debugEnabled == "true" {
		return log.New(os.Stdout, "", log.LstdFlags)
	} else if debugEnabled == "" || debugEnabled == "false" {
		return noplog
	} else {
		log.Panicf("Invalid value for DEBUG environment variable. Expected true / false but got " + debugEnabled)
	}

	return nil
}

// NOOP Logger --------------------------------------------

var noplog = &NopLogger{
	log.New(NullWriter(1), "", log.LstdFlags),
}

type NullWriter int

// Write implements the io.Write interface but is a noop.
func (NullWriter) Write([]byte) (int, error) { return 0, nil }

// NopLogger is a noop logger for passing to grpclog to minimize spew.
type NopLogger struct {
	*log.Logger
}

func (l *NopLogger) Fatal(args ...interface{}) {}

func (l *NopLogger) Fatalf(format string, args ...interface{}) {}

func (l *NopLogger) Fatalln(args ...interface{}) {}

func (l *NopLogger) Print(args ...interface{}) {}

func (l *NopLogger) Printf(format string, args ...interface{}) {}

func (l *NopLogger) Println(v ...interface{}) {}
