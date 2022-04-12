package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync/atomic"

	"github.com/deweppro/go-errors"
)

//nolint: golint
const (
	ANSI_RESET  = "\u001B[0m"
	ANSI_BLACK  = "\u001B[30m"
	ANSI_RED    = "\u001B[31m"
	ANSI_GREEN  = "\u001B[32m"
	ANSI_YELLOW = "\u001B[33m"
	ANSI_BLUE   = "\u001B[34m"
	ANSI_PURPLE = "\u001B[35m"
	ANSI_CYAN   = "\u001B[36m"

	eof = "\n"
)

var (
	scan       *bufio.Scanner
	yesNo             = []string{"y", "n"}
	debugLevel uint32 = 0
)

func init() {
	scan = bufio.NewScanner(os.Stdin)
}

func output(msg string, vars []string, def string) {
	if len(def) > 0 {
		def = fmt.Sprintf(" [%s]", def)
	}
	v := ""
	if len(vars) > 0 {
		v = fmt.Sprintf(" (%s)", strings.Join(vars, "/"))
	}
	Infof("%s%s%s: ", msg, v, def)
}

//Input console input request
func Input(msg string, vars []string, def string) string {
	output(msg, vars, def)

	for {
		if scan.Scan() {
			r := scan.Text()
			if len(r) == 0 {
				return def
			}
			if len(vars) == 0 {
				return r
			}
			for _, v := range vars {
				if v == r {
					return r
				}
			}
			output("Bad answer! Try again", vars, def)
		}
	}
}

//InputBool console bool input request
func InputBool(msg string, def bool) bool {
	v := "n"
	if def {
		v = "y"
	}
	v = Input(msg, yesNo, v)
	return v == "y"
}

func color(c, msg string, args []interface{}) {
	fmt.Printf(c+msg+ANSI_RESET, args...)
}

func colorln(c, msg string, args []interface{}) {
	if !strings.HasSuffix(msg, eof) {
		msg += eof
	}
	color(c, msg, args)
}

//Infof console message writer for info level
func Infof(msg string, args ...interface{}) {
	colorln(ANSI_RESET, "[INF] "+msg, args)
}

//Warnf console message writer for warning level
func Warnf(msg string, args ...interface{}) {
	colorln(ANSI_YELLOW, "[WAR] "+msg, args)
}

//Errorf console message writer for error level
func Errorf(msg string, args ...interface{}) {
	colorln(ANSI_RED, "[ERR] "+msg, args)
}

//ShowDebug init show debug
func ShowDebug(ok bool) {
	var v uint32 = 0
	if ok {
		v = 1
	}
	atomic.StoreUint32(&debugLevel, v)
}

//Debugf console message writer for debug level
func Debugf(msg string, args ...interface{}) {
	if atomic.LoadUint32(&debugLevel) > 0 {
		colorln(ANSI_BLUE, "[DEB] "+msg, args)
	}
}

//FatalIfErr console message writer if err is not nil
func FatalIfErr(err error, msg string, args ...interface{}) {
	if err != nil {
		Fatalf(errors.WrapMessage(err, msg, args...).Error())
	}
}

//Fatalf console message writer with exit code 1
func Fatalf(msg string, args ...interface{}) {
	colorln(ANSI_RED, "[ERR] "+msg, args)
	os.Exit(1)
}
