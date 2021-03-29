package application

import (
	"os"

	"github.com/deweppro/go-logger"
)

type (
	//ENV type for enviremants (prod, dev, stage, etc)
	ENV string

	//App application model
	App struct {
		cfile    string
		configs  Modules
		modules  Modules
		sources  Sources
		packages *DI
		logout   *Log
		logger   logger.Logger
		force    *ForceClose
	}
)

//New create application
func New() *App {
	return &App{
		modules:  Modules{},
		configs:  Modules{},
		packages: NewDI(),
		force:    newForceClose(),
	}
}

//Logger setup logger
func (a *App) Logger(log logger.Logger) *App {
	a.logger = log
	return a
}

//Modules append object to modules list
func (a *App) Modules(modules ...interface{}) *App {
	for _, mod := range modules {
		switch mod.(type) {
		case Modules:
			a.modules = a.modules.Add(mod.(Modules)...)
		default:
			a.modules = a.modules.Add(mod)
		}
	}

	return a
}

//ConfigFile set config file path and configs models
func (a *App) ConfigFile(filename string, configs ...interface{}) *App {
	a.cfile = filename
	a.configs = a.configs.Add(configs...)
	return a
}

//Run run application
func (a *App) Run() {
	var err error
	if len(a.cfile) > 0 {
		// read config file
		a.sources = Sources(a.cfile)

		// init logger
		config := &BaseConfig{}
		if err = a.sources.Decode(config); err != nil {
			panic(err)
		}
		a.logout = NewLogger(config)
		if a.logger == nil {
			a.logger = logger.Default()
		}
		a.logout.Handler(a.logger)
		a.modules = a.modules.Add(func() logger.Logger { return a.logger }, ENV(config.Env))

		// decode all configs
		if err = a.sources.Decode(a.configs...); err != nil {
			panic(err)
		}
		a.modules = a.modules.Add(a.configs...)

		if len(config.PidFile) > 0 {
			if err = pid2File(config.PidFile); err != nil {
				panic(err)
			}
		}
	}

	a.modules = a.modules.Add(a.force)
	a.launch()
}

func (a *App) launch() {
	var (
		err error
		ex  = 0
	)

	a.logger.Infof("app register components")
	if err = a.packages.Register(a.modules...); err != nil {
		a.logger.Fatalf("app register components: %s", err.Error())
	}

	a.logger.Infof("app build dependency")
	if err = a.packages.Build(); err != nil {
		a.logger.Fatalf("app build dependency: %s", err.Error())
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
