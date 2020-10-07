/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"os"
	"os/signal"
	"syscall"
)

func OnSyscallStop(callFunc func()) {
	quit := make(chan os.Signal, 4)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	callFunc()
}

func OnSyscallUp(callFunc func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP)
	<-quit

	callFunc()
}
