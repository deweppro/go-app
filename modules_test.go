package app_test

import (
	"testing"

	"github.com/deweppro/go-app"
	"github.com/stretchr/testify/require"
)

func TestUnit_Modules(t *testing.T) {
	tmp1 := app.Modules{8, 9, "W"}
	tmp2 := app.Modules{18, 19, "aW", tmp1}
	main := app.Modules{1, 2, "qqq"}.Add(tmp2).Add(99)

	require.Equal(t, app.Modules{1, 2, "qqq", 18, 19, "aW", 8, 9, "W", 99}, main)
}
