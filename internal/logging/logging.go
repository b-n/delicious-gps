package logging

import (
	"log"
	"os"
)

var (
	outputDebug bool
	infoLogger  *log.Logger
	debugLogger *log.Logger
	errorLogger *log.Logger
)

func Check(e error) {
	if e != nil {
		errorLogger.Fatal(e)
		panic(e)
	}
}

func Init(debugMode bool) {
	outputDebug = debugMode
	timeFormat := log.Ldate | log.Ltime

	debugLogger = log.New(os.Stdout, "DEBUG: ", timeFormat)
	infoLogger = log.New(os.Stdout, "INFO: ", timeFormat)
	errorLogger = log.New(os.Stderr, "ERROR: ", timeFormat)
}

func Debug(s string) {
	Debugf("%s", s)
}

func Debugf(s string, v ...interface{}) {
	if outputDebug {
		debugLogger.Printf(s, v...)
	}
}

func Info(s string) {
	Infof("%s", s)
}

func Infof(s string, v ...interface{}) {
	infoLogger.Printf(s, v)
}

func Error(e error) {
	errorLogger.Fatal(e)
}
