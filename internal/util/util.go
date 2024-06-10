package util

import (
	"fmt"
	"io"
	"log"
	"os"
)

var logger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)
	logger = log.New(multiWriter, "", log.Lshortfile|log.LstdFlags)
}

func LogError(err error) {
	if err != nil {
		logger.Output(2, err.Error())
	}
}

func HandlePanic() {
	if r := recover(); r != nil {
		LogError(fmt.Errorf("panic recovered: %v", r))
	}
}
