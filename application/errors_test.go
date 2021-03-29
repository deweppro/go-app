package application_test

import (
	"errors"
	"testing"

	"github.com/deweppro/go-app/application"

	"github.com/stretchr/testify/require"
)

func TestWrapErrors(t *testing.T) {
	require.Equal(
		t,
		application.WrapErrors(nil, errors.New("Hello"), "test").Error(),
		"test: Hello",
	)

	require.Equal(
		t,
		application.WrapErrors(nil, nil, "test"),
		nil,
	)

	require.Equal(
		t,
		application.WrapErrors(errors.New("Hello"), errors.New("World"), "test").Error(),
		"test: World: Hello",
	)
}
