package logger

import (
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger(prefix string) *Logger {
	return &Logger{log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lshortfile)}
}

func (l *Logger) Info(message string) {
	l.logger.Printf("[INFO] %s", message)
}

func (l *Logger) Warn(message string) {
	l.logger.Printf("[WARN] %s", message)
}

func (l *Logger) Error(message string) {
	l.logger.Printf("[ERROR] %s", message)
}

func (l *Logger) Fatal(message string) {
	l.logger.Fatalf("[FATAL] %s", message)
}
