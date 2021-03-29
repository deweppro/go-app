package application

import (
	"fmt"
	"reflect"

	"github.com/deweppro/go-algorithms/graph/kahn"
	"github.com/pkg/errors"
)

var types = []string{
	"invalid",
	"int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"uintptr", "float32", "float64", "complex64", "complex128",
	"bool", "array", "chan", "func", "interface", "map",
	"ptr", "slice", "string", "struct", "unsafe.Pointer",
}

type (
	// DI - managing dependencies
	DI struct {
		kahn *kahn.Graph
		srv  *services
		all  map[string]interface{}
	}
	typer interface {
		Elem() reflect.Type
		String() string
		PkgPath() string
		Name() string
	}
)

// NewDI - create new dependency manager
func NewDI() *DI {
	dep := &DI{
		kahn: kahn.New(),
		srv:  newServices(),
		all:  make(map[string]interface{}),
	}
	dep.all["error"] = new(error)
	return dep
}

func (d *DI) addr(t typer) string {
	if isDefaultType(t.Name()) {
		return t.String()
	}
	if len(t.PkgPath()) > 0 {
		return t.PkgPath() + ":" + t.String()
	}
	return t.Elem().PkgPath() + ":" + t.String()
}

// Register - register a new dependency
func (d *DI) Register(items ...interface{}) error {
	if d.srv.IsUp() {
		return ErrDepRunning
	}

	for _, item := range items {
		ref := reflect.TypeOf(item)

		switch ref.Kind() {

		case reflect.Func:
			for i := 0; i < ref.NumOut(); i++ {
				n := d.addr(ref.Out(i))
				if n == "error" {
					continue
				}
				d.all[n] = item
			}

		case reflect.Ptr, reflect.Struct:
			d.all[d.addr(ref)] = item

		default:
			if !isDefaultType(ref.Name()) {
				d.all[d.addr(ref)] = item
			}
		}
	}

	return nil
}

// Build - initialize dependencies
func (d *DI) Build() error {
	for out, item := range d.all {
		ref := reflect.TypeOf(item)

		switch ref.Kind() {

		case reflect.Func:
			if ref.NumIn() == 0 {
				if err := d.kahn.Add("error", out); err != nil {
					return errors.Wrapf(err, "cant add [error->%s] to graph", out)
				}
			}
			for i := 0; i < ref.NumIn(); i++ {
				in := d.addr(ref.In(i))
				if _, ok := d.all[in]; !ok {
					return fmt.Errorf("type is not found %s for %s", in, out)
				}
				if err := d.kahn.Add(in, out); err != nil {
					return errors.Wrapf(err, "cant add [%s->%s] to graph", in, out)
				}
			}

		case reflect.Struct:
			if ref.NumField() == 0 {
				if err := d.kahn.Add("error", out); err != nil {
					return errors.Wrapf(err, "cant add [error->%s] to graph", out)
				}
			}
			for i := 0; i < ref.NumField(); i++ {
				in := d.addr(ref.Field(i).Type)
				if _, ok := d.all[in]; !ok {
					return fmt.Errorf("type is not found %s for %s", in, out)
				}
				if err := d.kahn.Add(in, out); err != nil {
					return errors.Wrapf(err, "cant add [%s->%s] to graph", in, out)
				}
			}
		}
	}

	if err := d.kahn.Build(); err != nil {
		return errors.Wrap(err, "cant build graph")
	}

	names := make(map[string]struct{})
	for _, name := range d.kahn.Result() {
		names[name] = struct{}{}
	}

	for _, name := range d.kahn.Result() {
		if _, ok := names[name]; !ok {
			continue
		}

		if item, ok := d.all[name]; ok {
			if values, err := d.di(item); err == nil {
				for _, value := range values {
					if value.Type().AssignableTo(srvType) {
						if err = d.srv.Add(value.Interface().(Servicer)); err != nil {
							return errors.Wrap(err, "cant add element in graph")
						}
					}
					name = d.addr(value.Type())
					delete(names, name)
					d.all[d.addr(value.Type())] = value.Interface()
				}
			} else if !errors.Is(err, ErrBadAction) {
				return errors.Wrapf(err, "cant initialize %s", name)
			}

		} else {
			return errors.Wrapf(ErrDepUnknown, "dep: %s", name)
		}
	}

	return nil
}

// Down - stop all services in dependencies
func (d *DI) Down() error {
	return d.srv.Down()
}

// Up - start all services in dependencies
func (d *DI) Up() error {
	return d.srv.Up()
}

// Inject - obtained dependence
func (d *DI) Inject(item interface{}) error {
	_, err := d.di(item)
	return err
}

func (d *DI) di(item interface{}) ([]reflect.Value, error) {
	ref := reflect.TypeOf(item)
	args := make([]reflect.Value, 0)

	switch ref.Kind() {

	case reflect.Func:
		for i := 0; i < ref.NumIn(); i++ {
			in := d.addr(ref.In(i))
			if el, ok := d.all[in]; ok {
				args = append(args, reflect.ValueOf(el))
			} else {
				return nil, errors.Wrapf(ErrDepUnknown, "dep: %s", in)
			}
		}

	case reflect.Struct:
		value := reflect.New(ref)
		for i := 0; i < ref.NumField(); i++ {
			in := d.addr(ref.Field(i).Type)
			if el, ok := d.all[in]; ok {
				value.Elem().FieldByName(ref.Field(i).Name).Set(reflect.ValueOf(el))
			} else {
				return nil, errors.Wrap(ErrDepUnknown, in)
			}
		}
		return append(args, value), nil

	default:
		return nil, ErrBadAction
	}

	return reflect.ValueOf(item).Call(args), nil
}

func isDefaultType(name string) bool {
	for _, el := range types {
		if el == name {
			return true
		}
	}
	return false
}
