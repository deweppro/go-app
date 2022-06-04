package application

import (
	"os"

	"github.com/deweppro/go-logger"
)

//Log model
type Log struct {
	file    *os.File
	handler logger.Logger
	conf    *BaseConfig
}

//NewLogger init logger
func NewLogger(conf *BaseConfig) *Log {
	file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &Log{file: file, conf: conf}
}

//Handler set custom logger for application
func (l *Log) Handler(log logger.Logger) {
	l.handler = log
	l.handler.SetOutput(l.file)
	l.handler.SetLevel(l.conf.Level)
}

//Close log file
func (l *Log) Close() error {
	if l.handler != nil {
		l.handler.Close()
	}
	return l.file.Close()
}
