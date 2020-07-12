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

package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	file  *os.File
	debug bool
}

func MustNew(cfg *ConfigLog) *Logger {
	file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	l := &Logger{
		file:  file,
		debug: cfg.Env == "dev",
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(l)
	logrus.SetLevel(logrus.DebugLevel)

	return l
}

func (l *Logger) Write(p []byte) (n int, err error) {
	if l.debug {
		n, err = os.Stdout.Write(p)
	}

	return l.file.Write(p)
}

func (l *Logger) Down() error {
	return l.file.Close()
}
