package main

import (
	"fmt"
	"strings"

	"github.com/deweppro/go-app/console"
)

func main() {
	root := console.New("tool", "help tool")

	cmd := console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("simple", "first-level command")
		setter.Example("simple aa/bb/cc -a=hello -b=123 --cc=123.456 -e")

		setter.Flag(func(f console.FlagsSetter) {
			f.StringVar("a", "demo", "this is a string argument")
			f.IntVar("b", 1, "this is a int64 argument")
			f.FloatVar("cc", 1e-5, "this is a float64 argument")
			f.Bool("e", "this is a bool argument")
		})

		setter.Argument(1, func(s []string) ([]string, error) {
			if !strings.Contains(s[0], "/") {
				return nil, fmt.Errorf("argument must contain `/`")
			}
			return strings.Split(s[0], "/"), nil
		})

		setter.ExecFunc(func(args []string, a string, b int64, c float64, d bool) {
			fmt.Println(args, a, b, c, d)
		})
	})

	root.AddCommand(cmd)
	root.Exec()
}
