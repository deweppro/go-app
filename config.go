/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

//go:generate easyjson

//easyjson:json
type ConfigLogger struct {
	Env     string `yaml:"env" json:"env"`
	LogFile string `yaml:"log" json:"log"`
}
