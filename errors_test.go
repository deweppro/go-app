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
