package dic

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/deweppro/go-algorithms/graph/kahn"
	"github.com/deweppro/go-app/application/ctx"
	e "github.com/deweppro/go-app/application/error"
	"github.com/deweppro/go-app/application/service"
	"github.com/deweppro/go-app/internal"
	"github.com/deweppro/go-errors"
)

// Dic dependency injection container
type Dic struct {
	kahn *kahn.Graph
	srv  *service.Tree
	list *diData
}

// New dic constructor
func New() *Dic {
	return &Dic{
		kahn: kahn.New(),
		srv:  service.New(),
		list: newDiData(),
	}
}

// Down - stop all services in dependencies
func (v *Dic) Down(ctx ctx.Context) error {
	return v.srv.Down(ctx)
}

// Up - start all services in dependencies
func (v *Dic) Up(ctx ctx.Context) error {
	return v.srv.Up(ctx)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Register - register a new dependency
func (v *Dic) Register(items ...interface{}) error {
	if v.srv.IsUp() {
		return e.ErrDepRunning
	}

	for _, item := range items {
		ref := reflect.TypeOf(item)
		switch ref.Kind() {

		case reflect.Struct:
			if err := v.list.Add(item, item, typeExist); err != nil {
				return err
			}

		case reflect.Func:
			for i := 0; i < ref.NumIn(); i++ {
				in := ref.In(i)
				if in.Kind() == reflect.Struct {
					if err := v.list.Add(in, reflect.New(in).Elem().Interface(), typeNewIfNotExist); err != nil {
						return err
					}
				}

			}
			if ref.NumOut() == 0 {
				if err := v.list.Add(ref, item, typeNew); err != nil {
					return err
				}
				continue
			}
			for i := 0; i < ref.NumOut(); i++ {
				if err := v.list.Add(ref.Out(i), item, typeNew); err != nil {
					return err
				}
			}

		default:
			if err := v.list.Add(item, item, typeExist); err != nil {
				return err
			}
		}
	}

	return nil
}

// Build - initialize dependencies
func (v *Dic) Build() error {
	if v.srv.IsUp() {
		return e.ErrDepRunning
	}

	err := v.list.foreach(v.calcFunc, v.calcStruct, v.calcOther)
	if err != nil {
		return errors.WrapMessage(err, "building dependency graph")
	}

	if err = v.kahn.Build(); err != nil {
		return errors.WrapMessage(err, "dependency graph calculation")
	}

	return v.exec()
}

// Inject - obtained dependence
func (v *Dic) Inject(item interface{}) error {
	_, err := v.callArgs(item)
	return err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var empty = "EMPTY"

func (v *Dic) calcFunc(outAddr string, outRef reflect.Type) error {
	if outRef.NumIn() == 0 {
		if err := v.kahn.Add(empty, outAddr); err != nil {
			return errors.WrapMessage(err, "cant add [->%s] to graph", outAddr)
		}
	}

	for i := 0; i < outRef.NumIn(); i++ {
		inRef := outRef.In(i)
		inAddr, _ := internal.GetAddr(inRef)

		if _, err := v.list.Get(inAddr); err != nil {
			return errors.WrapMessage(err, "cant add [%s->%s] to graph", inAddr, outAddr)
		}
		if err := v.kahn.Add(inAddr, outAddr); err != nil {
			return errors.WrapMessage(err, "cant add [%s->%s] to graph", inAddr, outAddr)
		}
	}

	return nil
}

func (v *Dic) calcStruct(outAddr string, outRef reflect.Type) error {
	if outRef.NumField() == 0 {
		if err := v.kahn.Add(empty, outAddr); err != nil {
			return errors.WrapMessage(err, "cant add [->%s] to graph", outAddr)
		}
		return nil
	}
	for i := 0; i < outRef.NumField(); i++ {
		inRef := outRef.Field(i).Type
		inAddr, _ := internal.GetAddr(inRef)

		if _, err := v.list.Get(inAddr); err != nil {
			return errors.WrapMessage(err, "cant add [%s->%s] to graph", inAddr, outAddr)
		}
		if err := v.kahn.Add(inAddr, outAddr); err != nil {
			return errors.WrapMessage(err, "cant add [%s->%s] to graph", inAddr, outAddr)
		}
	}
	return nil
}

func (v *Dic) calcOther(_ string, _ reflect.Type) error {
	return nil
}

func (v *Dic) callFunc(item interface{}) ([]reflect.Value, error) {
	ref := reflect.TypeOf(item)
	args := make([]reflect.Value, 0, ref.NumIn())

	for i := 0; i < ref.NumIn(); i++ {
		inRef := ref.In(i)
		inAddr, _ := internal.GetAddr(inRef)
		vv, err := v.list.Get(inAddr)
		if err != nil {
			return nil, err
		}
		args = append(args, reflect.ValueOf(vv))
	}

	args = reflect.ValueOf(item).Call(args)
	for _, arg := range args {
		if err, ok := arg.Interface().(error); ok && err != nil {
			return nil, err
		}
	}

	return args, nil
}

func (v *Dic) callStruct(item interface{}) ([]reflect.Value, error) {
	ref := reflect.TypeOf(item)
	value := reflect.New(ref)
	args := make([]reflect.Value, 0, ref.NumField())

	for i := 0; i < ref.NumField(); i++ {
		inRef := ref.Field(i)
		inAddr, _ := internal.GetAddr(inRef.Type)
		vv, err := v.list.Get(inAddr)
		if err != nil {
			return nil, err
		}
		value.Elem().FieldByName(inRef.Name).Set(reflect.ValueOf(vv))
	}

	return append(args, value.Elem()), nil
}

func (v *Dic) callArgs(item interface{}) ([]reflect.Value, error) {
	ref := reflect.TypeOf(item)

	switch ref.Kind() {
	case reflect.Func:
		return v.callFunc(item)
	case reflect.Struct:
		return v.callStruct(item)
	default:
		return []reflect.Value{reflect.ValueOf(item)}, nil
	}
}

func (v *Dic) exec() error {
	names := make(map[string]struct{})
	for _, name := range v.kahn.Result() {
		if name == empty {
			continue
		}
		names[name] = struct{}{}
	}

	for _, name := range v.kahn.Result() {
		if _, ok := names[name]; !ok {
			continue
		}
		if v.list.HasType(name, typeExist) {
			continue
		}

		item, err := v.list.Get(name)
		if err != nil {
			return err
		}

		args, err := v.callArgs(item)
		if err != nil {
			return errors.WrapMessage(err, "initialize error <%s>", name)
		}

		for _, arg := range args {
			addr, _ := internal.GetAddr(arg.Type())
			if vv, ok := service.AsService(arg); ok {
				if err = v.srv.Add(vv); err != nil {
					return errors.WrapMessage(err, "service initialization error <%s>", addr)
				}
			}
			if vv, ok := service.AsServiceCtx(arg); ok {
				if err = v.srv.Add(vv); err != nil {
					return errors.WrapMessage(err, "service initialization error <%s>", addr)
				}
			}
			delete(names, addr)
			if arg.Type().String() == "error" {
				continue
			}
			if err = v.list.Add(arg.Type(), arg.Interface(), typeExist); err != nil {
				return errors.WrapMessage(err, "initialize error <%s>", addr)
			}
		}
		delete(names, name)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	typeNew int = iota
	typeNewIfNotExist
	typeExist
)

type (
	diElement struct {
		Value interface{}
		Type  int
	}
	diData struct {
		data map[string]*diElement
		mux  sync.RWMutex
	}
)

func newDiData() *diData {
	return &diData{
		data: make(map[string]*diElement),
	}
}

func (v *diData) Add(place, value interface{}, t int) error {
	v.mux.Lock()
	defer v.mux.Unlock()

	ref, ok := place.(reflect.Type)
	if !ok {
		ref = reflect.TypeOf(place)
	}

	addr, ok := internal.GetAddr(ref)
	if !ok {
		if addr != "error" {
			return fmt.Errorf("dependency <%s> is not supported", addr)
		}
		//return nil
	}

	if vv, ok := v.data[addr]; ok {
		if t == typeNewIfNotExist {
			return nil
		}
		if vv.Type == typeExist {
			return fmt.Errorf("dependency <%s> already initiated", addr)
		}
	}
	v.data[addr] = &diElement{
		Value: value,
		Type:  t,
	}

	return nil
}

func (v *diData) Get(addr string) (interface{}, error) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if vv, ok := v.data[addr]; ok {
		return vv.Value, nil
	}
	return nil, fmt.Errorf("dependency <%s> not initiated", addr)
}

func (v *diData) HasType(addr string, t int) bool {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if vv, ok := v.data[addr]; ok {
		return vv.Type == t
	}
	return false
}

func (v *diData) Step(addr string) (int, error) {
	v.mux.RLock()
	defer v.mux.RUnlock()

	if vv, ok := v.data[addr]; ok {
		return vv.Type, nil
	}
	return 0, fmt.Errorf("dependency <%s> not initiated", addr)
}

func (v *diData) foreach(kFunc, kStruct, kOther func(addr string, ref reflect.Type) error) error {
	v.mux.RLock()
	defer v.mux.RUnlock()

	for addr, item := range v.data {
		if item.Type == typeExist {
			continue
		}

		ref := reflect.TypeOf(item.Value)
		var err error
		switch ref.Kind() {
		case reflect.Func:
			err = kFunc(addr, ref)
		case reflect.Struct:
			err = kStruct(addr, ref)
		default:
			err = kOther(addr, ref)
		}

		if err != nil {
			return err
		}
	}
	return nil
}
