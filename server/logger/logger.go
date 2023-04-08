package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
)

func InitLogger() {
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error creating log file: ", err)
	}
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func getInfo() (file string, line int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	return
}

func LogInfo(v ...interface{}) {
	file, line := getInfo()
	InfoLogger.Printf("%s:%d - ", file, line)
	InfoLogger.Println(v...)
}

func LogError(v ...interface{}) {
	file, line := getInfo()
	ErrorLogger.Printf("%s:%d - ", file, line)
	ErrorLogger.Println(v...)
}

func LogWarn(v ...interface{}) {
	file, line := getInfo()
	WarnLogger.Printf("%s:%d - ", file, line)
	WarnLogger.Println(v...)
}