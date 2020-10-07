/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
