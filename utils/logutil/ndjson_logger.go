package logutil

import (
	"encoding/json"
	"fmt"
	"log"
)

type NdjsonLogger struct {
	infoLogger  *log.Logger
	debugLogger *log.Logger
	errorLogger *log.Logger
}

func (n *NdjsonLogger) Print(lvl LogLevel, includeTrace bool, v ...interface{}) {
	n.Println(lvl, includeTrace, v...)
}

func (n *NdjsonLogger) Println(lvl LogLevel, includeTrace bool, v ...interface{}) {
	le := newLogEntry(lvl, fmt.Sprintf("%v\n", v))
	if includeTrace {
		le.TraceBack = getTraceBack()
	}
	l, err := json.Marshal(le)
	if err != nil {
		log.Fatal(err)
		return
	}
	writeLogByLevel(n, lvl, l)
}

func (n *NdjsonLogger) Printf(lvl LogLevel, includeTrace bool, format string, v ...interface{}) {
	le := newLogEntry(lvl, fmt.Sprintf(format, v))
	if includeTrace {
		le.TraceBack = getTraceBack()
	}
	l, err := json.Marshal(le)
	if err != nil {
		log.Fatal(err)
	}
	writeLogByLevel(n, lvl, l)
}

func (n *NdjsonLogger) Fatal(v ...interface{}) {
	n.Println(Error, false, v...)
	log.Fatal(v)
}

func (n *NdjsonLogger) Fatalln(v ...interface{}) {
	n.Println(Error, false, v...)
	log.Fatalln(v)
}

func (n *NdjsonLogger) Fatalf(format string, v ...interface{}) {
	n.Printf(Error, false, format, v...)
	log.Fatalf(format, v)
}

func writeLogByLevel(n *NdjsonLogger, lvl LogLevel, b []byte) {
	switch lvl {
	case Info:
		n.infoLogger.Println(string(b))
	case Debug:
		n.debugLogger.Println(string(b))
	case Error:
		n.errorLogger.Println(string(b))
	default:
		log.Fatalf("unsupported log level: %v", lvl)
	}
}
