package logger

import (
	"fmt"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	logFile, err := os.OpenFile("./server.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Create log file failed...")
		os.Exit(1)
	}
	Logger = log.New(logFile, "", log.LstdFlags)
}
