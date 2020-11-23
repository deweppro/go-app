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

func (a *application) Logger(log logger.Logger) *application {
	a.logger = log
	return a
}

func (a *application) Modules(modules ...interface{}) *application {
	a.modules.Add(modules...)
	return a
}

func (a *application) ConfigFile(filename string, configs ...interface{}) *application {
	a.cfile = filename
	a.configs.Add(configs...)
	return a
}

func (a *application) Run() {
	var err error
	if len(a.cfile) > 0 {
		// read config file
		if a.sources, err = NewSources(a.cfile); err != nil {
			panic(err)
		}

		// init logger
		config := &ConfigLogger{}
		if err = a.sources.Decode(config); err != nil {
			panic(err)
		}
		a.logout = NewLogger(config.LogFile)
		if a.logger == nil {
			a.logger = logger.Default()
		}
		a.logout.Handler(a.logger)
		a.modules.Add(func() logger.Logger { return a.logger }, config, ENV(config.Env))

		// decode all configs
		configs := a.configs.Get()
		if err = a.sources.Decode(configs...); err != nil {
			panic(err)
		}
		a.modules.Add(configs...)

		if len(config.PidFile) > 0 {
			if err = pid2File(config.PidFile); err != nil {
				panic(err)
			}
		}
	}

	a.modules.Add(a.force)
	a.launch()
}

func (a *application) launch() {
	var (
		err error
		ex  = 0
	)

	a.logger.Infof("app register components")
	if err = a.packages.Register(a.modules.Get()); err != nil {
		a.logger.Errorf("app register components: %s", err.Error())
		os.Exit(1)
	}

	a.logger.Infof("app build dependency")
	if err = a.packages.Build(); err != nil {
		a.logger.Errorf("app build dependency: %s", err.Error())
		os.Exit(1)
	}

	a.logger.Infof("app up dependency")
	if err = a.packages.Up(); err != nil {
		a.logger.Errorf("app up dependency: %s", err.Error())
		ex++
	}

	if err == nil {
		go OnSyscallStop(a.force.Close)
		<-a.force.C.Done()
	}

	a.logger.Infof("app down dependency")
	if err = a.packages.Down(); err != nil {
		a.logger.Errorf("app down dependency: %s", err.Error())
		ex++
	}

	if err = a.logout.Close(); err != nil {
		panic(err)
	}

	if ex > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
