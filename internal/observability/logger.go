package observability

import (
	"log"
	"os"
)

type Logger struct {
	std *log.Logger
}

func NewLogger() *Logger {
	return &Logger{std: log.New(os.Stdout, "", log.LstdFlags|log.LUTC)}
}

func (l *Logger) Info(msg string) {
	l.std.Println("INFO", msg)
}

func (l *Logger) Error(msg string) {
	l.std.Println("ERROR", msg)
}
