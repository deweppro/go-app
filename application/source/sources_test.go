package source_test

import (
	"os"
	"testing"

	"github.com/deweppro/go-app/application"
	"github.com/deweppro/go-app/application/source"
	"github.com/stretchr/testify/require"
)

func TestUnit_YamlConfig(t *testing.T) {
	var (
		c1   = application.BaseConfig{}
		c2   = application.BaseConfig{}
		data = []byte(`
env: dev
log: hello
level: 5
`)
	)

	f, err := os.CreateTemp(os.TempDir(), "temp-*-config.yaml")
	require.NoError(t, err)
	_, err = f.Write(data)
	require.NoError(t, err)
	filename := f.Name()
	require.NoError(t, f.Close())
	require.NoError(t, source.Sources(filename).Decode(&c1, &c2))
	require.Equal(t, `hello`, c1.LogFile)
	require.Equal(t, `dev`, c1.Env)
	require.Equal(t, `hello`, c2.LogFile)
	require.Equal(t, `dev`, c2.Env)
}
