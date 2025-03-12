package logutil

import (
	"GoCouchbase/config"
	"errors"
	"log"
	"os"
	"sync"
)

var (
	infoLogger *NdjsonLogger
	once       sync.Once
)

func initLogger() {
	once.Do(func() {
		logDir, err := getLogDir()
		if err != nil {
			log.Fatal(err)
		}
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		infoLogFilePath := logDir + "info.ndjson"
		debugLogFilePath := logDir + "debug.ndjson"
		errorLogFilePath := logDir + "error.ndjson"

		infoLogFile, err := os.OpenFile(infoLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		debugLogFile, err := os.OpenFile(debugLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		errorLogFile, err := os.OpenFile(errorLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}

		infoLogger = &NdjsonLogger{
			infoLogger:  log.New(infoLogFile, "", 0),
			debugLogger: log.New(debugLogFile, "", 0),
			errorLogger: log.New(errorLogFile, "", 0),
		}
	})
}

func getLogDir() (string, error) {
	switch config.GetConfig().RunEnv {
	case config.NoEnv:
		return "", errors.New("RUN_ENV is not set")
	case config.Dev, config.Prod:
		return "/app/couchbase/logs", nil
	case config.OnPrem:
		return "", errors.New("OnPrem is not supported")
	case config.Local:
		return "./logs/", nil
	default:
		return "", errors.New("invalid RUN_ENV")
	}
}

func GetLogger() *NdjsonLogger {
	initLogger()
	return infoLogger
}
