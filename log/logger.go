package log

import (
	"fmt"
	"github.com/skyline/skyline-foundation/log/appender"
	"github.com/skyline/skyline-foundation/log/level"
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

func init() {
	initedLogger = make(map[string]Logger)
}

func NewLogger(logname string, logLevel level.LogLevel) Logger {
	if logger, ok := initedLogger[logname]; ok {
		return logger
	}
	var log Logger = &logger{name: logname,}
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
	fmt.Println(args)
}

func (l *logger) Info(args ...interface{}) {
	fmt.Println(args)
}

func (l *logger) Error(args ...interface{}) {
	fmt.Println(args)
}

func (l *logger) Warn(args ...interface{}) {
	fmt.Println(args)
}

func (l *logger) Fatal(args ...interface{}) {
	fmt.Println(args)
}
