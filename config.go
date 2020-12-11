/**
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

//go:generate easyjson

//ConfigLogger logger config model
//easyjson:json
type ConfigLogger struct {
	Env     string `yaml:"env" json:"env" toml:"env"`
	LogFile string `yaml:"log" json:"log" toml:"log"`
	PidFile string `yaml:"pid" json:"pid" toml:"pid"`
}
