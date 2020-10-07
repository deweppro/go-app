/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

type logger struct {
	file  *os.File
	debug bool
}

func newLogger(cfg *ConfigLogger) (*logger, error) {
	file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	l := &logger{
		file:  file,
		debug: cfg.Env == "dev",
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(l)
	if l.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return l, nil
}

func (l *logger) Write(p []byte) (n int, err error) {
	if l.debug {
		n, err = os.Stdout.Write(p)
	}
	return l.file.Write(p)
}

func (l *logger) Close() error {
	return l.file.Close()
}
