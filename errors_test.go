package app_test

import (
	"errors"
	"testing"

	"github.com/deweppro/go-app"
	"github.com/stretchr/testify/require"
)

func TestWrapErrors(t *testing.T) {
	require.Equal(
		t,
		app.WrapErrors(nil, errors.New("Hello"), "test").Error(),
		"test: Hello",
	)

	require.Equal(
		t,
		app.WrapErrors(nil, nil, "test"),
		nil,
	)

	require.Equal(
		t,
		app.WrapErrors(errors.New("Hello"), errors.New("World"), "test").Error(),
		"test: World: Hello",
	)
}
