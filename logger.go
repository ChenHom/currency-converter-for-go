package main

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

// initLogger initializes the logger with different log levels
func initLogger() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// logInfo logs informational messages
func logInfo(message string) {
	infoLogger.Println(message)
}

// logWarn logs warning messages
func logWarn(message string) {
	warnLogger.Println(message)
}

// logError logs error messages
func logError(message string) {
	errorLogger.Println(message)
}
