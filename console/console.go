package console

import (
	"os"
	"reflect"
)

const helpArg = "help"

type Console struct {
	name        string
	description string
	root        CommandGetter
}

func New(name, description string) *Console {
	return &Console{
		name:        name,
		description: description,
		root:        NewCommand(func(_ CommandSetter) {}).AsRoot(),
	}
}

func (c *Console) recover() {
	if d := recover(); d != nil {
		Fatalf("%+v", d)
	}
}

func (c *Console) AddCommand(getter ...CommandGetter) {
	defer c.recover()

	c.root.AddCommand(getter...)
}

func (c *Console) RootCommand(getter CommandGetter) {
	defer c.recover()

	next := c.root.List()
	c.root = getter.AsRoot()
	if err := c.root.Validate(); err != nil {
		Fatalf(err.Error())
	}
	c.root.AddCommand(next...)
}

func (c *Console) Exec() {
	defer c.recover()

	args := NewArgs().Parse(os.Args[1:])
	cmd, cur, h := c.build(args)
	if h {
		help(c.name, c.description, cmd, cur)
		return
	}
	c.run(cmd, args.Next()[len(cur):], args)
}

func (c *Console) build(args *Args) (command CommandGetter, cur []string, help bool) {
	var (
		i   int
		cmd string
	)
	for i, cmd = range args.Next() {
		if i == 0 {
			if nc := c.root.Next(cmd); nc != nil {
				command = nc
				continue
			}
			command = c.root
			break
		} else {
			if nc := command.Next(cmd); nc != nil {
				command = nc
				continue
			}
			break
		}
	}

	if len(args.Next()) > 0 {
		cur = args.Next()[:i]
	} else {
		command = c.root
	}

	if args.Has(helpArg) {
		help = true
	}

	return
}

func (c *Console) run(command CommandGetter, a []string, args *Args) {
	rv := make([]reflect.Value, 0)

	if command == nil || command.Call() == nil {
		Fatalf("command not found")
	}

	val, err := command.ArgCall(a)
	if err != nil {
		Fatalf("command \"%s\" validate arguments: %s", command.Name(), err.Error())
	}
	rv = append(rv, reflect.ValueOf(val))

	err = command.Flags().Call(args, func(i interface{}) {
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
