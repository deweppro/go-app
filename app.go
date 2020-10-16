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

type (
	ENV string

	application struct {
		cfile    string
		configs  *modules
		modules  *modules
		sources  *sources
		packages *DI
		logout   *log
		logger   logger.Logger
		force    *ForceClose
	}
)

func New() *application {
	return &application{
		modules:  NewModules(),
		configs:  NewModules(),
		packages: NewDI(),
		force:    newForceClose(),
	}
}

func (_app *application) Logger(log logger.Logger) *application {
	_app.logger = log
	return _app
}

func (_app *application) Modules(modules ...interface{}) *application {
	_app.modules.Add(modules...)
	return _app
}

func (_app *application) ConfigFile(filename string, configs ...interface{}) *application {
	_app.cfile = filename
	_app.configs.Add(configs...)
	return _app
}

func (_app *application) Run() {
	var err error
	if len(_app.cfile) > 0 {
		// read config file
		if _app.sources, err = NewSources(_app.cfile); err != nil {
			panic(err)
		}

		// init logger
		config := &ConfigLogger{}
		if err := _app.sources.Decode(config); err != nil {
			panic(err)
		}
		_app.logout = NewLogger(config.LogFile)
		if _app.logger == nil {
			_app.logger = logger.Default()
		}
		_app.logout.Handler(_app.logger)
		_app.modules.Add(func() logger.Logger { return _app.logger }, config, ENV(config.Env))

		// decode all configs
		configs := _app.configs.Get()
		if err := _app.sources.Decode(configs...); err != nil {
			panic(err)
		}
		_app.modules.Add(configs...)

		if len(config.PidFile) > 0 {
			if err = pid2File(config.PidFile); err != nil {
				panic(err)
			}
		}
	}

	_app.modules.Add(_app.force)
	_app.launch()
}

func (_app *application) launch() {
	var (
		err error
		ex  = 0
	)

	_app.logger.Infof("app register components")
	if err = _app.packages.Register(_app.modules.Get()); err != nil {
		_app.logger.Errorf("app register components: %s", err.Error())
		os.Exit(1)
	}

	_app.logger.Infof("app build dependency")
	if err = _app.packages.Build(); err != nil {
		_app.logger.Errorf("app build dependency: %s", err.Error())
		os.Exit(1)
	}

	_app.logger.Infof("app up dependency")
	if err = _app.packages.Up(); err != nil {
		_app.logger.Errorf("app up dependency: %s", err.Error())
		ex++
	}

	if err == nil {
		go OnSyscallStop(_app.force.Close)
		<-_app.force.C.Done()
	}

	_app.logger.Infof("app down dependency")
	if err = _app.packages.Down(); err != nil {
		_app.logger.Errorf("app down dependency: %s", err.Error())
		ex++
	}

	if err = _app.logout.Close(); err != nil {
		panic(err)
	}

	if ex > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
