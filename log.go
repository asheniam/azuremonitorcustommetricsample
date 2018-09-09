package main

import (
	"fmt"
)

var (
	stdLogger = newLogger()
)

type logger struct {
}

func newLogger() *logger {
	return &logger{}
}

func (l *logger) logerror(args ...interface{}) {
	fmt.Println(args)
}

func (l *logger) logerrorf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (l *logger) loginfof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func loginfof(format string, args ...interface{}) {
	stdLogger.loginfof(format, args...)
}

func logerrorf(format string, args ...interface{}) {
	stdLogger.loginfof(format, args...)
}

func logerror(args ...interface{}) {
	stdLogger.logerror(args...)
}
