package main

import (
	"fmt"

	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-logger"
)

//nolint: golint
type (
	Test0 struct{}
	Test1 struct{}
	Test2 struct{}

	Config struct {
		Env   string `yaml:"env"`
		Level string `yaml:"level"`
	}

	Params struct {
		Test1  *Test1
		Config Config
	}
)

func (s *Test2) Up() error {
	fmt.Println("--> call *Test2.Up")
	return nil
}

func (s *Test2) Down() error {
	fmt.Println("--> call *Test2.Down")
	return nil
}

func NewTest0(p Params) *Test0 {
	fmt.Println("--> call NewTest0")
	fmt.Println("--> Params.Config.Env=" + p.Config.Env)
	return &Test0{}
}

func NewTest2(_ *Test0) *Test2 {
	fmt.Println("--> call NewTest2")
	return &Test2{}
}

func main() {
	application.New().
		Logger(logger.Default()).
		ConfigFile(
			"./config.yaml",
			Config{},
		).
		Modules(
			&Test1{},
			NewTest0,
			NewTest2,
		).
		Run()
}
