/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"io/ioutil"
	"strconv"
	"syscall"
)

func pid2File(filename string) error {
	pid := strconv.Itoa(syscall.Getpid())
	return ioutil.WriteFile(filename, []byte(pid), 0755)
}
