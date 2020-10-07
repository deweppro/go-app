/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapErrors(t *testing.T) {
	require.Equal(
		t,
		WrapErrors(nil, errors.New("Hello"), "test").Error(),
		"test: Hello",
	)

	require.Equal(
		t,
		WrapErrors(nil, nil, "test"),
		nil,
	)

	require.Equal(
		t,
		WrapErrors(errors.New("Hello"), errors.New("World"), "test").Error(),
		"test: World: Hello",
	)
}
