package internal

import (
	"io/ioutil"
	"strconv"
	"syscall"
)

//PidFile write pid file
func PidFile(filename string) error {
	pid := strconv.Itoa(syscall.Getpid())
	return ioutil.WriteFile(filename, []byte(pid), 0755)
}
