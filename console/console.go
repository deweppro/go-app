package console

import (
	"fmt"
	"os"
	"reflect"
)

const helpArg = "help"

type Console struct {
	name        string
	description string
	next        []CommandGetter
}

func New(name, description string) *Console {
	return &Console{
		name:        name,
		description: description,
		next:        make([]CommandGetter, 0),
	}
}

func (c *Console) recover() {
	if d := recover(); d != nil {
		Fatalf("%+v", d)
	}
}

func (c *Console) AddCommand(getter ...CommandGetter) {
	defer c.recover()
	for _, v := range getter {
		if err := v.Validate(); err != nil {
			Fatalf(err.Error())
		}
		c.next = append(c.next, v)
	}
}

func (c *Console) Exec() {
	//defer c.recover()
	if err := c.validate(); err != nil {
		Fatalf(err.Error())
	}

	args := NewArgs().Parse(os.Args[1:])
	next, cmd, cur, h := c.build(args)
	if h {
		help(c.name, c.description, next, cmd, cur)
		return
	}
	c.run(cmd, args.Next()[len(cur):], args)
}

func (c *Console) validate() error {
	if len(c.name) == 0 {
		return fmt.Errorf("command name is empty")
	}
	return nil
}

func (c *Console) build(args *Args) (next []CommandGetter, command CommandGetter, cur []string, help bool) {
	var (
		i   int
		cmd string
	)
	next = c.next
	for i, cmd = range args.Next() {
		for _, command = range next {
			if !command.Is(cmd) {
				continue
			}
			next = command.Next()
			if len(next) > 0 {
				goto NEXT
			}
			goto END
		}

		if args.Has(helpArg) {
			cur, help = args.Next()[:i], true
			return
		} else {
			Fatalf("command not found")
		}

	NEXT:
		continue
	END:
		break
	}

	if len(args.Next()) > 0 {
		cur = args.Next()[:i+1]
	}
	if args.Has(helpArg) {
		help = true
		return
	}

	if command == nil {
		Fatalf("command not found")
	}
	return
}

func (c *Console) run(command CommandGetter, a []string, args *Args) {
	rv := make([]reflect.Value, 0)

	if command.ArgCount() > 0 {
		if len(a) < command.ArgCount() {
			Fatalf("command \"%s\" arguments must be - %d", command.ArgCount())
		}
		if val, err := command.ArgCall(a[:command.ArgCount()]); err != nil {
			Fatalf("command \"%s\" validate arguments: %s", command.Name(), err.Error())
		} else {
			rv = append(rv, reflect.ValueOf(val))
		}
	} else {
		rv = append(rv, reflect.ValueOf(nil))
	}

	err := command.Flags().Call(args, func(i interface{}) {
		rv = append(rv, reflect.ValueOf(i))
	})
	if err != nil {
		Fatalf("command \"%s\" validate flags: %s", command.Name(), err.Error())
	}

	if reflect.ValueOf(command.Call()).Type().NumIn() != len(rv) {
		Fatalf("command \"%s\" Flags: fewer arguments declared than expected in ExecFunc", command.Name())
	}

	reflect.ValueOf(command.Call()).Call(rv)
}
