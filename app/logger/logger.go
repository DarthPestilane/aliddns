package logger

import (
	"log"
	"os"
)

var logEntry *Logger

type Logger struct {
	client *log.Logger
}

func (l *Logger) Info(msg string, ctx ...interface{}) {
	l.client.Printf("INFO %s %v\n", msg, ctx)
}

func (l *Logger) Error(msg string, ctx ...interface{}) {
	l.client.Printf("ERROR %s %v\n", msg, ctx)
}

func Register() {
	logEntry = &Logger{client: log.New(os.Stdout, "[ALIDDNS]", log.LstdFlags)}
}

func Provide() *Logger {
	return logEntry
}
