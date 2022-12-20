package console

import "strings"

type (
	//ValidFunc validate argument interface
	ValidFunc func([]string) ([]string, error)
	//Argument model
	Argument struct {
		ValidFunc ValidFunc
	}
)

// NewArgument constructor
func NewArgument() *Argument {
	return &Argument{}
}

type (
	//Args list model
	Args struct {
		list []Arg
		next []string
	}
	//Arg model
	Arg struct {
		Key   string
		Value string
	}
	//ArgGetter argument getter interface
	ArgGetter interface {
		Has(name string) bool
		Get(name string) *string
	}
)

// NewArgs constructor
func NewArgs() *Args {
	return &Args{
		list: make([]Arg, 0),
		next: make([]string, 0),
	}
}

func (a *Args) Has(name string) bool {
	for _, v := range a.list {
		if v.Key == name {
			return true
		}
	}
	return false
}

func (a *Args) Get(name string) *string {
	for _, v := range a.list {
		if v.Key == name {
			return &v.Value
		}
	}
	return nil
}

func (a *Args) Next() []string {
	return a.next
}

func (a *Args) Parse(list []string) *Args {
	for i := 0; i < len(list); i++ {
		// args
		if strings.HasPrefix(list[i], "-") {
			arg := Arg{}
			v := strings.TrimLeft(list[i], "-")
			vs := strings.SplitN(v, "=", 2)
			if len(vs) == 2 {
				arg.Key, arg.Value = vs[0], vs[1]
				a.list = append(a.list, arg)
				continue
			}

			if i+1 < len(list) && !strings.HasPrefix(list[i+1], "-") {
				arg.Key, arg.Value = vs[0], list[i+1]
				a.list = append(a.list, arg)
				i++
				continue
			}

			arg.Key = vs[0]
			a.list = append(a.list, arg)
			continue
		}
		//commands
		a.next = append(a.next, list[i])
	}

	return a
}
