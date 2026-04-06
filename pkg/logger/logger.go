package logger

import (
	"log"
	"os"
)

type Logger struct {
	environment string
	logger      *log.Logger
}

func NewLogger(environment string) *Logger {
	return &Logger{
		environment: environment,
		logger:      log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Printf("[INFO] %s %v\n", message, args)
}

func (l *Logger) Error(message string, err error) {
	l.logger.Printf("[ERROR] %s: %v\n", message, err)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Printf("[WARN] %s %v\n", message, args)
}

func (l *Logger) Debug(message string, args ...interface{}) {
	if l.environment == "development" {
		l.logger.Printf("[DEBUG] %s %v\n", message, args)
	}
}

func (l *Logger) Fatal(message string, err error) {
	l.logger.Fatalf("[FATAL] %s: %v\n", message, err)
}
