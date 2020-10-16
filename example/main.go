/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"fmt"

	"github.com/deweppro/go-app"
)

var _ app.ServiceInterface = (*Simple)(nil)

type (
	Simple       struct{}
	SimpleConfig struct {
		Env string `yaml:"env"`
	}
)

func NewSimple(_ *SimpleConfig) *Simple {
	fmt.Println("call NewSimple")
	return &Simple{}
}

func (s *Simple) Up() error {
	fmt.Println("call *Simple.Up")
	return nil
}

func (s *Simple) Down() error {
	fmt.Println("call *Simple.Down")
	return nil
}

func main() {
	app.New().
		ConfigFile(
			"./config.yaml",
			&SimpleConfig{},
		).
		Modules(
			NewSimple,
		).
		Run()
}
