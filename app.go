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
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type application struct {
	config  *sources
	logger  *logger
	modules *modules
	dep     *DI
}

func New(fileconfig string) *application {
	src, err := newSources(fileconfig)
	if err != nil {
		panic(err.Error())
	}
	log := &ConfigLogger{}
	if err = src.YAML(log); err != nil {
		panic(err.Error())
	}
	logger, err := newLogger(log)
	if err != nil {
		panic(err.Error())
	}
	return &application{
		config:  src,
		logger:  logger,
		modules: newModules(),
		dep:     NewDI(),
	}
}

func (a *application) ConfigModels(configs ...interface{}) *application {
	if err := a.config.YAML(configs...); err != nil {
		logrus.Fatal(err)
	}
	a.modules.Add(configs...)
	return a
}

func (a *application) Modules(modules ...interface{}) *application {
	a.modules.Add(modules...)
	return a
}

func (a *application) PidFile(pid string) *application {
	err := ioutil.WriteFile(pid, []byte(fmt.Sprintf("%d", syscall.Getpid())), 0755)
	if err != nil {
		logrus.Fatal(err)
	}
	return a
}

func (a *application) Run() {
	var (
		fc  = newForceClose()
		err error
		ex  = 0
	)

	logrus.Info("app register components")
	if err = a.dep.Register(a.modules.Add(fc).Get()); err != nil {
		logrus.Fatal(errors.Wrap(err, "[app register components]"))
	}

	logrus.Info("app build dependency")
	if err = a.dep.Build(); err != nil {
		logrus.Fatal(errors.Wrap(err, "[app build dependency]"))
	}

	logrus.Info("app up dependency")
	if err = a.dep.Up(); err != nil {
		logrus.Error(errors.Wrap(err, "[app up dependency]"))
		ex++
	}

	if err == nil {
		go OnSyscallStop(fc.Close)
		<-fc.C.Done()
	}

	if err = a.dep.Down(); err != nil {
		logrus.Error(errors.Wrap(err, "[app down dependency]"))
		ex++
	}
	logrus.Info("app down dependency")

	if err := a.logger.Close(); err != nil {
		logrus.Fatal(errors.Wrap(err, "[app logger down]"))
	}

	if ex > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
