package application

import (
	"os"

	"github.com/deweppro/go-app/application/ctx"
	"github.com/deweppro/go-app/application/dic"
	"github.com/deweppro/go-app/application/source"
	"github.com/deweppro/go-app/console"
	"github.com/deweppro/go-app/internal"
	"github.com/deweppro/go-logger"
	"github.com/deweppro/go-utils/syscall"
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

// New create application
func New() *App {
	return &App{
		modules:  Modules{},
		configs:  Modules{},
		packages: dic.New(),
		ctx:      ctx.New(),
	}
}

// Logger setup logger
func (a *App) Logger(log logger.Logger) *App {
	a.logger = log
	return a
}

// Modules append object to modules list
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

// ConfigFile set config file path and configs models
func (a *App) ConfigFile(filename string, configs ...interface{}) *App {
	a.cfile = filename
	for _, config := range configs {
		a.configs = a.configs.Add(config)
	}

	return a
}

// Run run application
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
	result := a.steps(
		[]step{
			{
				Message: "app register components",
				Call:    func() error { return a.packages.Register(a.modules...) },
			},
			{
				Message: "app build dependency",
				Call:    func() error { return a.packages.Build() },
			},
			{
				Message: "app up dependency",
				Call:    func() error { return a.packages.Up(a.ctx) },
			},
		},
		func(er bool) {
			if er {
				a.ctx.Close()
				return
			}
			go syscall.OnStop(a.ctx.Close)
			<-a.ctx.Done()
		},
		[]step{
			{
				Message: "app down dependency",
				Call:    func() error { return a.packages.Down(a.ctx) },
			},
		},
	)
	console.FatalIfErr(a.logout.Close(), "close log file")
	if result {
		os.Exit(1)
	}
	os.Exit(0)
}

type step struct {
	Call    func() error
	Message string
}

func (a *App) steps(up []step, wait func(bool), down []step) bool {
	var erc int

	for _, s := range up {
		a.logger.Infof(s.Message)
		if err := s.Call(); err != nil {
			a.logger.WithFields(logger.Fields{
				"err": err.Error(),
			}).Errorf(s.Message)
			erc++
			break
		}
	}

	wait(erc > 0)

	for _, s := range down {
		a.logger.Infof(s.Message)
		if err := s.Call(); err != nil {
			a.logger.WithFields(logger.Fields{
				"err": err.Error(),
			}).Errorf(s.Message)
			erc++
		}
	}

	return erc > 0
}
