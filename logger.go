package gdit

import (
	"log"
	"os"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

type standardLogger struct{}

func (l *standardLogger) Info(args ...interface{}) {
	log.SetOutput(os.Stdout)
	log.Println("[GDIT-INFO]", args)
}

func (l *standardLogger) Warn(args ...interface{}) {
	log.SetOutput(os.Stdout)
	log.Println("[GDIT-WARN]", args)
}

func (l *standardLogger) Error(args ...interface{}) {
	log.SetOutput(os.Stderr)
	log.Println("[GDIT-ERROR]", args)
}
