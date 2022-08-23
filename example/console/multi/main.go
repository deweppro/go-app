package main

import (
	"fmt"
	"strings"

	"github.com/deweppro/go-app/console"
)

func main() {
	console.ShowDebug(true)

	app := console.New("tool", "help tool")

	cmd := console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("a", "command a")
		setter.ExecFunc(func(args []string) {
			fmt.Println("a", args)
		})

		setter.AddCommand(console.NewCommand(func(setter console.CommandSetter) {
			setter.Setup("b", "command b")
			setter.ExecFunc(func(args []string) {
				fmt.Println("b", args)
			})
		}))

	})

	root := console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("root", "command root")
		setter.Flag(func(setter console.FlagsSetter) {
			setter.Bool("aaa", "bool a")
		})
		setter.ArgumentFunc(func(s []string) ([]string, error) {
			return []string{strings.Join(s, "-")}, nil
		})
		setter.ExecFunc(func(args []string, a bool) {
			fmt.Println("root", args, a)
		})
	})

	app.RootCommand(root)
	app.AddCommand(cmd)
	app.Exec()
}
