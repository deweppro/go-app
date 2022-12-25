package internal

import (
	"math/rand"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/deweppro/go-utils/random"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// PidFile write pid file
func PidFile(filename string) error {
	pid := strconv.Itoa(syscall.Getpid())
	return os.WriteFile(filename, []byte(pid), 0755)
}

func RandString(n int) string {
	return string(random.Bytes(n))
}
