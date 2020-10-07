/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
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
