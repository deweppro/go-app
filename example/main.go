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

package main

import (
	"fmt"

	"github.com/deweppro/go-app"
)

var _ app.ServiceInterface = (*Simple)(nil)

type Simple struct{}

func NewSimple(_ *app.ConfigLogger) *Simple {
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
	app.
		New("config.yaml").
		ConfigModels(&app.ConfigLogger{}).
		Modules(NewSimple).
		PidFile("/tmp/app.pid").
		Run()
}
