package application

import (
	"io/ioutil"
	"strconv"
	"syscall"
)

func pid2File(filename string) error {
	pid := strconv.Itoa(syscall.Getpid())
	return ioutil.WriteFile(filename, []byte(pid), 0755)
}
