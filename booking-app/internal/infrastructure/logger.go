package infrastructure

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
}

func NewLogger() Logger {
	return Logger{
		info:  log.New(os.Stdout, "[INFO]: ", log.LstdFlags),
		error: log.New(os.Stdout, "[ERROR]: ", log.LstdFlags),
	}
}

func (l Logger) Info(format string, v ...any) {
	l.info.Printf(format, v...)
}

func (l Logger) Error(format string, v ...any) {
	l.error.Printf(format, v...)
}
