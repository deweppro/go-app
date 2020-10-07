/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
 */

package app

import (
	"fmt"
	"reflect"

	"github.com/deweppro/go-algorithms/graph/kahn"
	"github.com/pkg/errors"
)

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

func (_di *DI) addr(t typer) string {
	if isDefaultType(t.Name()) {
		return t.String()
	}
	if len(t.PkgPath()) > 0 {
		return t.PkgPath() + ":" + t.String()
	}
	return t.Elem().PkgPath() + ":" + t.String()
}

// Register - register a new dependency
func (_di *DI) Register(items []interface{}) error {
	if _di.srv.IsUp() {
		return ErrDepRunning
	}

	for _, item := range items {
		ref := reflect.TypeOf(item)

		switch ref.Kind() {

		case reflect.Func:
			for i := 0; i < ref.NumOut(); i++ {
				n := _di.addr(ref.Out(i))
				if n == "error" {
					continue
				}
				_di.all[n] = item
			}

		case reflect.Struct:
			_di.all[_di.addr(reflect.New(reflect.TypeOf(item)).Type())] = item

		case reflect.Ptr:
			_di.all[_di.addr(ref)] = item

		default:
			if !isDefaultType(ref.Name()) {
				_di.all[_di.addr(ref)] = item
			}
		}
	}

	return nil
}

// Build - initialize dependencies
func (_di *DI) Build() error {
	for out, item := range _di.all {
		ref := reflect.TypeOf(item)

		switch ref.Kind() {

		case reflect.Func:
			if ref.NumIn() == 0 {
				if err := _di.kahn.Add("error", out); err != nil {
					return errors.Wrapf(err, "cant add [error->%s] to graph", out)
				}
			}
			for i := 0; i < ref.NumIn(); i++ {
				in := _di.addr(ref.In(i))
				if _, ok := _di.all[in]; !ok {
					return fmt.Errorf("type is not found %s for %s", in, out)
				}
				if err := _di.kahn.Add(in, out); err != nil {
					return errors.Wrapf(err, "cant add [%s->%s] to graph", in, out)
				}
			}

		case reflect.Struct:
			if ref.NumField() == 0 {
				if err := _di.kahn.Add("error", out); err != nil {
					return errors.Wrapf(err, "cant add [error->%s] to graph", out)
				}
			}
			for i := 0; i < ref.NumField(); i++ {
				in := _di.addr(ref.Field(i).Type)
				if _, ok := _di.all[in]; !ok {
					return fmt.Errorf("type is not found %s for %s", in, out)
				}
				if err := _di.kahn.Add(in, out); err != nil {
					return errors.Wrapf(err, "cant add [%s->%s] to graph", in, out)
				}
			}
		}

	}

	if err := _di.kahn.Build(); err != nil {
		return errors.Wrap(err, "cant build graph")
	}

	for _, name := range _di.kahn.Result() {
		if item, ok := _di.all[name]; ok {

			if values, err := _di.di(item); err == nil {
				for _, value := range values {
					if value.Type().AssignableTo(srvType) {
						if err = _di.srv.Add(value.Interface().(ServiceInterface)); err != nil {
							return errors.Wrap(err, "cant add element in graph")
						}
					}

					_di.all[name] = value.Interface()
				}
			} else if !errors.Is(err, ErrBadAction) {
				return errors.Wrapf(err, "cant initialize %s", name)
			}

		} else {
			return ErrDepUnknown
		}
	}

	return nil
}

// Down - stop all services in dependencies
func (_di *DI) Down() error {
	return _di.srv.Down()
}

// Up - start all services in dependencies
func (_di *DI) Up() error {
	return _di.srv.Up()
}

// Inject - obtained dependence
func (_di *DI) Inject(item interface{}) error {
	_, err := _di.di(item)
	return err
}

func (_di *DI) di(item interface{}) ([]reflect.Value, error) {
	ref := reflect.TypeOf(item)
	args := make([]reflect.Value, 0)

	switch ref.Kind() {

	case reflect.Func:
		for i := 0; i < ref.NumIn(); i++ {
			in := _di.addr(ref.In(i))
			if el, ok := _di.all[in]; ok {
				args = append(args, reflect.ValueOf(el))
			} else {
				return nil, errors.Wrap(ErrDepUnknown, in)
			}
		}

	case reflect.Struct:
		value := reflect.New(ref)
		for i := 0; i < ref.NumField(); i++ {
			in := _di.addr(ref.Field(i).Type)
			if el, has := _di.all[in]; has {
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
