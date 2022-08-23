package console

import (
	"fmt"
	"reflect"
)

type Command struct {
	root     bool
	name     string
	desc     string
	examples []string
	flags    *Flags
	args     *Argument
	execute  interface{}

	next []CommandGetter
}

type CommandGetter interface {
	Next(string) CommandGetter
	List() []CommandGetter
	Validate() error
	Is(string) bool
	Name() string
	Description() string
	Examples() []string
	ArgCall(d []string) ([]string, error)
	Flags() FlagsGetter
	Call() interface{}
	AddCommand(...CommandGetter)
	AsRoot() CommandGetter
	IsRoot() bool
}

type CommandSetter interface {
	Setup(string, string)
	Example(string)
	Flag(cb func(FlagsSetter))
	ArgumentFunc(call ValidFunc)
	ExecFunc(interface{})
	AddCommand(...CommandGetter)
}

func NewCommand(cb func(CommandSetter)) CommandGetter {
	cmd := &Command{
		next:     make([]CommandGetter, 0),
		flags:    NewFlags(),
		args:     NewArgument(),
		examples: make([]string, 0),
	}
	cb(cmd)
	return cmd
}

func (c *Command) Setup(name, description string) {
	c.name, c.desc = name, description
}

func (c *Command) AsRoot() CommandGetter {
	c.root = true
	c.name = ""
	return c
}

func (c *Command) IsRoot() bool {
	return c.root
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Description() string {
	return c.desc
}

func (c *Command) Examples() []string {
	return c.examples
}

func (c *Command) Example(s string) {
	c.examples = append(c.examples, s)
}

func (c *Command) Flag(cb func(FlagsSetter)) {
	cb(c.flags)
}

func (c *Command) Flags() FlagsGetter {
	return c.flags
}

func (c *Command) ArgumentFunc(call ValidFunc) {
	c.args.ValidFunc = call
}

func (c *Command) ArgCall(d []string) ([]string, error) {
	if c.args.ValidFunc == nil {
		return d, nil
	}
	return c.args.ValidFunc(d)
}

func (c *Command) ExecFunc(i interface{}) {
	c.execute = i
}

func (c *Command) Next(cmd string) CommandGetter {
	for _, getter := range c.next {
		if getter.Is(cmd) {
			return getter
		}
	}
	return nil
}

func (c *Command) List() []CommandGetter {
	return c.next
}

func (c *Command) Validate() error {
	if len(c.name) == 0 && !c.IsRoot() {
		return fmt.Errorf("command name is empty. use Setup(name, description)")
	}
	if reflect.ValueOf(c.execute).Kind() != reflect.Func {
		return fmt.Errorf("command [%s] ExecFunc: is not a func", c.name)
	}
	count := c.flags.Count() + 1
	if reflect.ValueOf(c.execute).Type().NumIn() != count {
		return fmt.Errorf("command [%s] Flags: fewer arguments declared than expected in ExecFunc", c.name)
	}
	return nil
}

func (c *Command) Call() interface{} {
	return c.execute
}

func (c *Command) Is(s string) bool {
	return c.name == s
}

func (c *Command) AddCommand(getter ...CommandGetter) {
	for _, v := range getter {
		if err := v.Validate(); err != nil {
			Fatalf(err.Error())
		}
		c.next = append(c.next, v)
	}
}
