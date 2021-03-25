package app_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/deweppro/go-app"
	"github.com/stretchr/testify/require"
)

func TestUnit_YamlConfig(t *testing.T) {
	var (
		c1   = app.BaseConfig{}
		c2   = app.BaseConfig{}
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
	require.NoError(t, app.Sources(filename).Decode(&c1, &c2))
	require.Equal(t, `hello`, c1.LogFile)
	require.Equal(t, `dev`, c1.Env)
	require.Equal(t, `hello`, c2.LogFile)
	require.Equal(t, `dev`, c2.Env)
}

func TestUnit_JsonConfig(t *testing.T) {
	var (
		c1   = app.BaseConfig{}
		c2   = app.BaseConfig{}
		data = []byte(`{"env":"dev","log":"hello","level":5}`)
	)

	f, err := ioutil.TempFile(os.TempDir(), "temp-*-config.json")
	require.NoError(t, err)
	_, err = f.Write(data)
	require.NoError(t, err)
	filename := f.Name()
	require.NoError(t, f.Close())
	require.NoError(t, app.Sources(filename).Decode(&c1, &c2))
	require.Equal(t, `hello`, c1.LogFile)
	require.Equal(t, `dev`, c1.Env)
	require.Equal(t, `hello`, c2.LogFile)
	require.Equal(t, `dev`, c2.Env)
}

func TestUnit_TomlConfig(t *testing.T) {
	var (
		c1   = app.BaseConfig{}
		c2   = app.BaseConfig{}
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
	require.NoError(t, app.Sources(filename).Decode(&c1, &c2))
	require.Equal(t, `hello`, c1.LogFile)
	require.Equal(t, `dev`, c1.Env)
	require.Equal(t, `hello`, c2.LogFile)
	require.Equal(t, `dev`, c2.Env)
}
