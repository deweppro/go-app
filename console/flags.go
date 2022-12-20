package console

import (
	"fmt"
	"strconv"
)

type (
	//Flags model
	Flags struct {
		d []FlagItem
	}
	//FlagItem element of flag model
	FlagItem struct {
		req   bool
		name  string
		value interface{}
		usage string
		call  func(getter ArgGetter) (interface{}, error)
	}
)

// FlagsGetter getter interface
type FlagsGetter interface {
	Info(cb func(bool, string, interface{}, string))
	Call(g ArgGetter, cb func(interface{})) error
}

// FlagsSetter setter interface
type FlagsSetter interface {
	StringVar(name string, value string, usage string)
	String(name string, usage string)
	IntVar(name string, value int64, usage string)
	Int(name string, usage string)
	FloatVar(name string, value float64, usage string)
	Float(name string, usage string)
	Bool(name string, usage string)
}

// NewFlags init new flag
func NewFlags() *Flags {
	return &Flags{
		d: make([]FlagItem, 0),
	}
}

// Count of flags
func (f *Flags) Count() int {
	return len(f.d)
}

// Info about command
func (f *Flags) Info(cb func(req bool, name string, v interface{}, usage string)) {
	for _, item := range f.d {
		cb(item.req, item.name, item.value, item.usage)
	}
}

func (f *Flags) Call(g ArgGetter, cb func(interface{})) error {
	for _, item := range f.d {
		v, err := item.call(g)
		if err != nil {
			return err
		}
		cb(v)
	}
	return nil
}

// StringVar flag decoder with default value
func (f *Flags) StringVar(name string, value string, usage string) {
	f.d = append(f.d, FlagItem{
		req:   false,
		name:  name,
		value: value,
		usage: usage,
		call: func(getter ArgGetter) (interface{}, error) {
			if val := getter.Get(name); val != nil {
				return *val, nil
			}
			return value, nil
		},
	})
}

// String flag decoder
func (f *Flags) String(name string, usage string) {
	f.d = append(f.d, FlagItem{
		req:   true,
		name:  name,
		usage: usage,
		call: func(getter ArgGetter) (interface{}, error) {
			if val := getter.Get(name); val != nil && len(*val) > 0 {
				return *val, nil
			}
			return nil, fmt.Errorf("--%s is not found", name)
		},
	})
}

// IntVar flag decoder with default value
func (f *Flags) IntVar(name string, value int64, usage string) {
	f.d = append(f.d, FlagItem{
		req:   false,
		value: value,
		name:  name,
		usage: usage,
		call: func(getter ArgGetter) (interface{}, error) {
			if val := getter.Get(name); val != nil && len(*val) > 0 {
				return strconv.ParseInt(*val, 10, 64)
			}
			return value, nil
		},
	})
}

// Int flag decoder
func (f *Flags) Int(name string, usage string) {
	f.d = append(f.d, FlagItem{
		req:   true,
		value: 0,
		name:  name,
		usage: usage,
		call: func(getter ArgGetter) (interface{}, error) {
			if val := getter.Get(name); val != nil && len(*val) > 0 {
				return strconv.ParseInt(*val, 10, 64)
			}
			return nil, fmt.Errorf("--%s is not found", name)
		},
	})
}

// FloatVar flag decoder with default value
func (f *Flags) FloatVar(name string, value float64, usage string) {
	f.d = append(f.d, FlagItem{
		req:   false,
		value: value,
		name:  name,
		usage: usage,
		call: func(getter ArgGetter) (interface{}, error) {
			if val := getter.Get(name); val != nil && len(*val) > 0 {
				return strconv.ParseFloat(*val, 64)
			}
			return value, nil
		},
	})
}

// Float flag decoder
func (f *Flags) Float(name string, usage string) {
	f.d = append(f.d, FlagItem{
		req:   true,
		value: 0.0,
		name:  name,
		usage: usage,
		call: func(getter ArgGetter) (interface{}, error) {
			if val := getter.Get(name); val != nil && len(*val) > 0 {
				return strconv.ParseFloat(*val, 64)
			}
			return nil, fmt.Errorf("--%s is not found", name)
		},
	})
}

// Bool flag decoder
func (f *Flags) Bool(name string, usage string) {
	f.d = append(f.d, FlagItem{
		req:   false,
		value: false,
		name:  name,
		usage: usage,
		call: func(getter ArgGetter) (interface{}, error) {
			if val := getter.Get(name); val != nil {
				return true, nil
			}
			return false, nil
		},
	})
}
