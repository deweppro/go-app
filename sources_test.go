/**
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	_, err = f.Write(data)
	require.NoError(t, err)
	filename := f.Name()
	require.NoError(t, f.Close())
	require.NoError(t, Sources(filename).Decode(&c1, &c2))
	require.Equal(t, `hello`, c1.LogFile)
	require.Equal(t, `dev`, c1.Env)
	require.Equal(t, `hello`, c2.LogFile)
	require.Equal(t, `dev`, c2.Env)
}

func TestUnit_JsonConfig(t *testing.T) {
	var (
		c1   = ConfigLogger{}
		c2   = ConfigLogger{}
		data = []byte(`{"env":"dev","log":"hello","level":5}`)
	)

	f, err := ioutil.TempFile(os.TempDir(), "temp-*-config.json")
	require.NoError(t, err)
	_, err = f.Write(data)
	require.NoError(t, err)
	filename := f.Name()
	require.NoError(t, f.Close())
	require.NoError(t, Sources(filename).Decode(&c1, &c2))
	require.Equal(t, `hello`, c1.LogFile)
	require.Equal(t, `dev`, c1.Env)
	require.Equal(t, `hello`, c2.LogFile)
	require.Equal(t, `dev`, c2.Env)
}

func TestUnit_TomlConfig(t *testing.T) {
	var (
		c1   = ConfigLogger{}
		c2   = ConfigLogger{}
		data = []byte(`
env = "dev"
log = "hello"
`)
	)

	f, err := ioutil.TempFile(os.TempDir(), "temp-*-config.toml")
	require.NoError(t, err)
	_, err = f.Write(data)
	require.NoError(t, err)
	filename := f.Name()
	require.NoError(t, f.Close())
	require.NoError(t, Sources(filename).Decode(&c1, &c2))
	require.Equal(t, `hello`, c1.LogFile)
	require.Equal(t, `dev`, c1.Env)
	require.Equal(t, `hello`, c2.LogFile)
	require.Equal(t, `dev`, c2.Env)
}
