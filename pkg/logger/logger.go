package logger

import (
	"log"
	"os"
)

type Logger struct {
	Error *log.Logger
	Info  *log.Logger
}

func (l *Logger) New() *Logger {
	return &Logger{
		Error: log.New(os.Stderr, "ERROR : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		Info:  log.New(os.Stdout, "INFO : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}
}
