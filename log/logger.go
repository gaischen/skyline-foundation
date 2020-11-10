package log

import (
	"fmt"
	"github.com/vanga-top/skyline-foundation/log/appender"
	"github.com/vanga-top/skyline-foundation/log/level"
	"sync"
	"time"
)

type Logger interface {
	Level() level.LogLevel
	Name() string //log name always log file's name
	SetLevel(logLevel level.LogLevel)
	GetLogger() Logger
	SetAppender(appender appender.Appender)

	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Fatal(args ...interface{})
}

var initedLogger map[string]Logger
var mutex sync.Mutex

func init() {
	initedLogger = make(map[string]Logger)
}

func NewLogger(logname string, logLevel level.LogLevel) Logger {
	if logger, ok := initedLogger[logname]; ok {
		return logger
	}
	defer mutex.Unlock()
	mutex.Lock()
	var log Logger = &logger{name: logname,}
	initedLogger[logname] = log //set to map
	return log
}

type logger struct {
	name  string
	level level.LogLevel
	appender.Appender
}

func (l *logger) Level() level.LogLevel {
	return l.level
}

func (l *logger) Name() string {
	return l.name
}

func (l *logger) SetLevel(logLevel level.LogLevel) {
	if logLevel < 0 {
		return
	}
	l.level = logLevel
}

func (l *logger) GetLogger() Logger {
	return l
}

func (l *logger) SetAppender(appender appender.Appender) {
	panic("implement me")
}

func (l *logger) Debug(args ...interface{}) {
	setPrefix(args...)
}

func (l *logger) Info(args ...interface{}) {
	setPrefix(args...)
}

func (l *logger) Error(args ...interface{}) {
	setPrefix(args...)
}

func (l *logger) Warn(args ...interface{}) {
	setPrefix(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	setPrefix(args...)
	panic("os out...")
}

func setPrefix(args ...interface{}) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timeStr, args)
}
