package utils

import (
	"fmt"
	"time"
)

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

// Log is a global logger instance
var Log = New()

var (
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	reset  = "\033[0m"
)

func (l *Logger) Info(msg string) {
	fmt.Printf("%s[%s] INFO%s %s\n", green, timestamp(), reset, msg)
}

func (l *Logger) Warn(msg string) {
	fmt.Printf("%s[%s] WARN%s %s\n", yellow, timestamp(), reset, msg)
}

func (l *Logger) Error(msg string) {
	fmt.Printf("%s[%s] ERROR%s %s\n", red, timestamp(), reset, msg)
}

func (l *Logger) Debug(msg string) {
	fmt.Printf("%s[%s] DEBUG%s %s\n", blue, timestamp(), reset, msg)
}

func timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
