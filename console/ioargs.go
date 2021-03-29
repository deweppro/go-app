package console

import "strings"

type (
	ioArgs []ioArg
	ioArg  struct {
		Key   string
		Value string
	}
	IOArgsGetter interface {
		Has(name string) bool
		Get(name string) *string
	}
)

func (f ioArgs) Has(name string) bool {
	for _, v := range f {
		if v.Key == name {
			return true
		}
	}
	return false
}

func (f ioArgs) Get(name string) *string {
	for _, v := range f {
		if v.Key == name {
			return &v.Value
		}
	}
	return nil
}

func ioArgsParse(list []string) (next []string, args ioArgs) {
	for i := 0; i < len(list); i++ {
		// args
		if strings.HasPrefix(list[i], "-") {
			arg := ioArg{}
			v := strings.TrimLeft(list[i], "-")
			vs := strings.SplitN(v, "=", 2)
			if len(vs) == 2 {
				arg.Key, arg.Value = vs[0], vs[1]
				args = append(args, arg)
				continue
			}

			if i+1 < len(list) && !strings.HasPrefix(list[i+1], "-") {
				arg.Key, arg.Value = vs[0], list[i+1]
				args = append(args, arg)
				i++
				continue
			}

			arg.Key = vs[0]
			args = append(args, arg)
			continue
		}
		//commands
		next = append(next, list[i])
	}

	return
}
