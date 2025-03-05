package utils

import (
	"GoCouchbase/config"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	once        sync.Once
)

func initLogger() {
	once.Do(func() {
		cfg := config.GetConfig()

		infoLogPath, err := getLogPath(cfg.RunEnv, "info")
		if err != nil {
			log.Fatal(err)
		}
		errLogPath, err := getLogPath(cfg.RunEnv, "error")
		if err != nil {
			log.Fatal(err)
		}
		infoLogFile, err := os.OpenFile(infoLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("error opening info log file: %v", err)
		}
		errLogFile, err := os.OpenFile(errLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("error opening error log file: %v", err)
		}
		infoMultiWriter := io.MultiWriter(os.Stdout, infoLogFile)
		errorMultiWriter := io.MultiWriter(os.Stderr, errLogFile)

		infoLogger = log.New(infoMultiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		errorLogger = log.New(errorMultiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

		log.SetOutput(infoMultiWriter)
	})
}

func getLogPath(runEnv config.RunEnv, loggerType string) (string, error) {
	var logFileName string
	switch loggerType {
	case "info":
		logFileName = "app.log"
	case "error":
		logFileName = "err.log"
	default:
		return "", errors.New("invalid logger type")
	}

	var logDir string
	switch runEnv {
	case config.NoEnv:
		log.Fatal("RUN_ENV is not set")
	case config.Dev, config.Prod:
		if _, err := os.Stat("/app/couchbase/logs"); os.IsNotExist(err) {
			err = os.MkdirAll("/app/couchbase/logs", 0755)
			if err != nil {
				log.Fatalf("Failed to create log directory: %v", err)
			}
		}
		logDir = "/app/couchbase/logs"
	case config.OnPrem:
		logDir = "/home/ubuntu"
	case config.Local:
		logDir = "."
	default:
		log.Fatalf("RUN_ENV is invalid")
	}

	logPath := fmt.Sprintf("%s/%s", logDir, logFileName)
	return logPath, nil
}

func GetInfoLogger() *log.Logger {
	initLogger()
	return infoLogger
}

func GetErrorLogger() *log.Logger {
	initLogger()
	return errorLogger
}
