/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"os"

	"github.com/deweppro/go-logger"
)

type log struct {
	file    *os.File
	handler logger.Logger
}

func NewLogger(filename string) *log {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &log{file: file}
}

func (l *log) Handler(log logger.Logger) {
	l.handler = log
	l.handler.SetOutput(l.file)
}

func (l *log) Close() error {
	return l.file.Close()
}
