package logutil

import (
	"GoCouchbase/config"
	"time"
)

type LogEntry struct {
	Env       string    `json:"env"`
	Timestamp string    `json:"timestamp"`
	Level     string    `json:"lvl"`
	Message   string    `json:"msg"`
	Service   string    `json:"svc"`
	TraceBack *[]string `json:"tb,omitempty"`
}

func newLogEntry(lvl LogLevel, msg string) *LogEntry {
	cfg := config.GetConfig()
	return &LogEntry{
		Env:       cfg.RunEnv.String(),
		Timestamp: time.Now().Format("2006-01-02T15:04:05"),
		Level:     lvl.String(),
		Message:   msg,
		Service:   "couchbase_bridge",
		TraceBack: nil,
	}
}
