package infrastructure

import (
	"log"
	"os"
)

// Logger represents a realization of Logger interface
type Logger struct {
	// info logger
	info *log.Logger
	// error logger
	error *log.Logger
}

// NewLogger returns new logger with two embedded loggers for different levels.
func NewLogger() Logger {
	return Logger{
		info:  log.New(os.Stdout, "[INFO]: ", log.LstdFlags),
		error: log.New(os.Stdout, "[ERROR]: ", log.LstdFlags),
	}
}

// Info prints message, specified by format and any variables, via the info logger.
func (l Logger) Info(format string, v ...any) {
	l.info.Printf(format, v...)
}

// Info prints message, specified by format and any variables, via the error logger.
func (l Logger) Error(format string, v ...any) {
	l.error.Printf(format, v...)
}
