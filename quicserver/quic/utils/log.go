package utils

// LogLevel of quic-go
type LogLevel uint8

const (
	// LogLevelNothing disables
	LogLevelNothing LogLevel = iota
	// LogLevelError enables err logs
	LogLevelError
	// LogLevelInfo enables info logs (e.g. packets)
	LogLevelInfo
	// LogLevelDebug enables debug logs (e.g. packet contents)
	LogLevelDebug
)

type Logger interface {
	SetLogLevel(LogLevel)
	SetLogTimeFormat(format string)
	WithPrefix(prefix string) Logger
	Debug() bool

	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

var DefaultLogger Logger
