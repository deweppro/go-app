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

package app

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_YamlConfig(t *testing.T) {
	var (
		c1   = ConfigLogger{}
		c2   = ConfigLogger{}
		data = []byte(`
env: dev
log: hello
level: 5
`)
	)

	f, err := ioutil.TempFile(os.TempDir(), "temp-*-config.yaml")
	assert.NoError(t, err)
	_, err = f.Write(data)
	assert.NoError(t, err)
	filename := f.Name()
	assert.NoError(t, f.Close())
	conf, err := newSources(filename)
	assert.NoError(t, err)
	assert.NoError(t, conf.YAML(&c1, &c2))
	assert.Equal(t, `hello`, c1.LogFile)
	assert.Equal(t, `dev`, c1.Env)
	assert.Equal(t, `hello`, c2.LogFile)
	assert.Equal(t, `dev`, c2.Env)
}

func TestUnit_JsonConfig(t *testing.T) {
	var (
		c1   = ConfigLogger{}
		c2   = ConfigLogger{}
		data = []byte(`{"env":"dev","log":"hello","level":5}`)
	)

	f, err := ioutil.TempFile(os.TempDir(), "temp-*-config.json")
	assert.NoError(t, err)
	_, err = f.Write(data)
	assert.NoError(t, err)
	filename := f.Name()
	assert.NoError(t, f.Close())
	conf, err := newSources(filename)
	assert.NoError(t, err)
	assert.NoError(t, conf.JSON(&c1, &c2))
	assert.Equal(t, `hello`, c1.LogFile)
	assert.Equal(t, `dev`, c1.Env)
	assert.Equal(t, `hello`, c2.LogFile)
	assert.Equal(t, `dev`, c2.Env)
}
