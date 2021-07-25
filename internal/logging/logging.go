package logging

import (
	"fmt"
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

	debugLogger = log.New(os.Stdout, "", timeFormat)
	infoLogger = log.New(os.Stdout, "", timeFormat)
	errorLogger = log.New(os.Stderr, "", timeFormat)
}

func Debug(s string) {
	Debugf("%s", s)
}

func Debugf(s string, v ...interface{}) {
	if outputDebug {
		debugLogger.Printf("[DEBUG]: %s", fmt.Sprintf(s, v...))
	}
}

func Info(s string) {
	Infof("%s", s)
}

func Infof(s string, v ...interface{}) {
	infoLogger.Printf("[INFO]: %s", fmt.Sprintf(s, v...))
}

func Error(e error) {
	errorLogger.Fatal(e)
}
