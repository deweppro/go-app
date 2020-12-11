/**
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnit_Modules(t *testing.T) {
	tmp1 := Modules{8, 9, "W"}
	tmp2 := Modules{18, 19, "aW", tmp1}
	main := Modules{1, 2, "qqq"}.Add(tmp2).Add(99)

	require.Equal(t, Modules{1, 2, "qqq", 18, 19, "aW", 8, 9, "W", 99}, main)
}
