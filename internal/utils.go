package internal

import (
	"io/ioutil"
	"math/rand"
	"strconv"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// PidFile write pid file
func PidFile(filename string) error {
	pid := strconv.Itoa(syscall.Getpid())
	return ioutil.WriteFile(filename, []byte(pid), 0755)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
