package main

import (
	"fmt"

	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-app/application/ctx"
	"github.com/deweppro/go-logger"
)

type (
	//Simple model
	Simple struct{}
	//Config model
	Config struct {
		Env string `yaml:"env"`
	}
)

//NewSimple init Simple
func NewSimple(_ Config) *Simple {
	fmt.Println("--> call NewSimple")
	return &Simple{}
}

//Up  method for start Simple in DI container
func (s *Simple) Up(_ ctx.Context) error {
	fmt.Println("--> call *Simple.Up")
	return nil
}

//Down  method for stop Simple in DI container
func (s *Simple) Down(_ ctx.Context) error {
	fmt.Println("--> call *Simple.Down")
	return nil
}

func main() {
	application.New().
		Logger(logger.Default()).
		ConfigFile(
			"./config.yaml",
			Config{},
		).
		Modules(
			NewSimple,
		).
		Run()
}
