package main

import (
	"fmt"

	"github.com/deweppro/go-app"
	"github.com/deweppro/go-logger"
)

var _ app.Servicer = (*Simple)(nil)

type (
	//Simple model
	Simple struct{}
	//SimpleConfig config model
	SimpleConfig struct {
		Env string `yaml:"env"`
	}
)

//NewSimple init Simple
func NewSimple(_ *SimpleConfig) *Simple {
	fmt.Println("call NewSimple")
	return &Simple{}
}

//Up  method for start Simple in DI container
func (s *Simple) Up() error {
	fmt.Println("call *Simple.Up")
	return nil
}

//Down  method for stop Simple in DI container
func (s *Simple) Down() error {
	fmt.Println("call *Simple.Down")
	return nil
}

func main() {
	app.New().
		Logger(logger.Default()).
		ConfigFile(
			"./config.yaml",
			&SimpleConfig{},
		).
		Modules(
			NewSimple,
		).
		Run()
}
