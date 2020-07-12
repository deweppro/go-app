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

	"github.com/deweppro/go-app/pkg/event"
	"github.com/deweppro/go-app/pkg/filedecoder"
	"github.com/deweppro/go-app/pkg/initializer"
	"github.com/deweppro/go-app/pkg/logger"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type App struct {
	cfgdata []byte
	logger  *logger.Logger
	modules *Interfaces
	dep     *initializer.Dependencies
}

func New(fileconfig string) *App {
	data := filedecoder.ReadFile(fileconfig)
	logcfg := &logger.ConfigLog{}
	if err := filedecoder.Yaml(data, []interface{}{logcfg}); err != nil {
		logrus.Fatal(err)
	}

	return &App{
		cfgdata: data,
		logger:  logger.MustNew(logcfg),
		modules: NewInterfaces(),
		dep:     initializer.New(),
	}
}

func (app *App) ConfigModels(configs ...interface{}) *App {
	if err := filedecoder.Yaml(app.cfgdata, configs); err != nil {
		logrus.Fatal(err)
	}
	app.modules.Add(configs...)
	return app
}

func (app *App) Modules(modules ...interface{}) *App {
	app.modules.Add(modules...)
	return app
}

func (app *App) PidFile(pid string) *App {
	err := ioutil.WriteFile(pid, []byte(fmt.Sprintf("%d", syscall.Getpid())), 0755)
	if err != nil {
		logrus.Fatal(err)
	}

	return app
}

func (app *App) Run() {
	fc := NewForceClose()

	logrus.Info("app register components")
	if err := app.dep.Register(app.modules.Add(fc).Get()); err != nil {
		logrus.Fatal(errors.Wrap(err, "[app register components]"))
	}

	logrus.Info("app build dependency")
	if err := app.dep.Build(); err != nil {
		logrus.Fatal(errors.Wrap(err, "[app build dependency]"))
	}

	logrus.Info("app up dependency")
	errup := app.dep.Up()
	if errup != nil {
		logrus.Error(errors.Wrap(errup, "[app up dependency]"))
	}

	if errup == nil {
		go event.OnSyscallStop(fc.Close)
		<-fc.C.Done()
	}

	errdown := app.dep.Down()
	if errdown != nil {
		logrus.Error(errors.Wrap(errdown, "[app down dependency]"))
	}
	logrus.Info("app down dependency")

	if err := app.logger.Down(); err != nil {
		logrus.Fatal(errors.Wrap(err, "[app logger down]"))
	}

	if errup != nil || errdown != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
