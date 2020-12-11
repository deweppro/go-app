/**
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"os"

	"github.com/deweppro/go-logger"
)

//Log model
type Log struct {
	file    *os.File
	handler logger.Logger
}

//NewLogger init logger
func NewLogger(filename string) *Log {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &Log{file: file}
}

//Handler set custom logger for application
func (l *Log) Handler(log logger.Logger) {
	l.handler = log
	l.handler.SetOutput(l.file)
}

//Close log file
func (l *Log) Close() error {
	return l.file.Close()
}
