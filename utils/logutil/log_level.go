package logutil

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR"}[l]
}
