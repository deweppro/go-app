package application

import (
	"os"

	"github.com/deweppro/go-app/application/ctx"
	"github.com/deweppro/go-app/application/dic"
	"github.com/deweppro/go-app/application/source"
	"github.com/deweppro/go-app/application/sys"
	"github.com/deweppro/go-app/console"
	"github.com/deweppro/go-app/internal"
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
		sources  source.Sources
		packages *dic.Dic
		logout   *Log
		logger   logger.Logger
		ctx      ctx.Context
	}
)

//New create application
func New() *App {
	return &App{
		modules:  Modules{},
		configs:  Modules{},
		packages: dic.New(),
		ctx:      ctx.New(),
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
		switch v := mod.(type) {
		case Modules:
			a.modules = a.modules.Add(v...)
		default:
			a.modules = a.modules.Add(v)
		}
	}

	return a
}

//ConfigFile set config file path and configs models
func (a *App) ConfigFile(filename string, configs ...interface{}) *App {
	a.cfile = filename
	for _, config := range configs {
		a.configs = a.configs.Add(config)
	}

	return a
}

//Run run application
func (a *App) Run() {
	var err error
	if len(a.cfile) == 0 {
		a.logout = NewLogger(&BaseConfig{
			Level:   4,
			LogFile: "/dev/stdout",
		})
		a.logger = logger.Default()
		a.logout.Handler(a.logger)
	}
	if len(a.cfile) > 0 {
		// read config file
		a.sources = source.Sources(a.cfile)

		// init logger
		config := &BaseConfig{}
		if err = a.sources.Decode(config); err != nil {
			console.FatalIfErr(err, "decode config file: %s", a.cfile)
		}
		a.logout = NewLogger(config)
		if a.logger == nil {
			a.logger = logger.Default()
		}
		a.logout.Handler(a.logger)
		a.modules = a.modules.Add(func() logger.Logger { return a.logger }, ENV(config.Env))

		// decode all configs
		var configs []interface{}
		configs, err = internal.TypingPtr(a.configs, func(i interface{}) error {
			return a.sources.Decode(i)
		})
		if err != nil {
			a.logger.WithFields(logger.Fields{
				"err": err.Error(),
			}).Fatalf("decode config file")
		}
		a.modules = a.modules.Add(configs...)

		if len(config.PidFile) > 0 {
			if err = internal.PidFile(config.PidFile); err != nil {
				a.logger.WithFields(logger.Fields{
					"err":  err.Error(),
					"file": config.PidFile,
				}).Fatalf("create pid file")
			}
		}
	}

	a.modules = a.modules.Add(a.ctx)
	a.launch()
}

func (a *App) launch() {
	var (
		err error
		ex  = 0
	)

	a.logger.Infof("app register components")
	if err = a.packages.Register(a.modules...); err != nil {
		a.logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Fatalf("app register components")
	}

	a.logger.Infof("app build dependency")
	if err = a.packages.Build(); err != nil {
		a.logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Fatalf("app build dependency")
	}

	a.logger.Infof("app up dependency")
	if err = a.packages.Up(a.ctx); err != nil {
		a.logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Errorf("app up dependency")
		ex++
	}

	if err == nil {
		go sys.OnSyscallStop(a.ctx.Close)
		<-a.ctx.Done()
	}

	a.logger.Infof("app down dependency")
	if err = a.packages.Down(a.ctx); err != nil {
		a.logger.WithFields(logger.Fields{
			"err": err.Error(),
		}).Errorf("app down dependency")
		ex++
	}

	if err = a.logout.Close(); err != nil {
		console.FatalIfErr(err, "close log file")
	}

	if ex > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}
