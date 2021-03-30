package application

import (
	"os"
	"os/signal"
	"syscall"
)

//OnSyscallStop calling a function if you send a system event stop
func OnSyscallStop(callFunc func()) {
	quit := make(chan os.Signal, 4)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	callFunc()
}

//OnSyscallUp calling a function if you send a system event SIGHUP
func OnSyscallUp(callFunc func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP)
	<-quit

	callFunc()
}
