package main

import (
	"log"
	"os"
)

type Logger struct {
	client *log.Logger
}

func NewLogger() *Logger {
	client := log.New(os.Stdout, "[ALIDDNS] ", log.LstdFlags)
	return &Logger{
		client: client,
	}
}

func (l *Logger) Info(msg string, ctx ...interface{}) {
	l.client.Printf("INFO %s %v\n", msg, ctx)
}

func (l *Logger) Error(msg string, ctx ...interface{}) {
	l.client.Printf("ERROR %s %v\n", msg, ctx)
}
