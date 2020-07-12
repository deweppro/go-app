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

package initializer

import (
	errors2 "errors"
	"fmt"
	"reflect"

	"github.com/deweppro/algorithms/graph/kahn"
	"github.com/pkg/errors"
)

// Dependencies - managing dependencies
type Dependencies struct {
	kahn *kahn.Kahn
	srv  *services
	all  map[string]interface{}
}

// New - create new dependency manager
func New() *Dependencies {
	dep := &Dependencies{
		kahn: kahn.New(),
		srv:  newServices(),
		all:  make(map[string]interface{}),
	}
	dep.all["error"] = new(error)
	return dep
}

// Register - register a new dependency
func (dep *Dependencies) Register(items []interface{}) error {
	if dep.srv.IsUp() {
		return errorDepRunning
	}

	for _, item := range items {
		ref := reflect.TypeOf(item)

		switch ref.Kind() {

		case reflect.Func:
			for i := 0; i < ref.NumOut(); i++ {
				t := ref.Out(i)
				n := t.String()

				if n == "error" {
					continue
				}

				dep.all[n] = item
			}

		case reflect.Struct:
			dep.all[reflect.New(reflect.TypeOf(item)).Type().String()] = item

		case reflect.Ptr:
			dep.all[ref.String()] = item

		default:
			if !isDefaultType(ref.Name()) {
				dep.all[ref.String()] = item
			}
		}
	}

	return nil
}

// Build - initialize dependencies
func (dep *Dependencies) Build() error {
	for out, item := range dep.all {

		ref := reflect.TypeOf(item)

		switch ref.Kind() {

		case reflect.Func:
			if ref.NumIn() == 0 {
				if err := dep.kahn.Add("error", out); err != nil {
					return errors.Wrapf(err, "cant add [error->%s] to graph", out)
				}
			}

			for i := 0; i < ref.NumIn(); i++ {
				it := ref.In(i)
				in := it.String()

				if _, ok := dep.all[in]; !ok {
					return fmt.Errorf("type is not found %s for %s", in, out)
				}

				if err := dep.kahn.Add(in, out); err != nil {
					return errors.Wrapf(err, "cant add [%s->%s] to graph", in, out)
				}
			}

		case reflect.Struct:
			if ref.NumField() == 0 {
				if err := dep.kahn.Add("error", out); err != nil {
					return errors.Wrapf(err, "cant add [error->%s] to graph", out)
				}
			}

			for i := 0; i < ref.NumField(); i++ {
				it := ref.Field(i)
				in := it.Type.String()

				if _, ok := dep.all[in]; !ok {
					return fmt.Errorf("type is not found %s for %s", in, out)
				}

				if err := dep.kahn.Add(in, out); err != nil {
					return errors.Wrapf(err, "cant add [%s->%s] to graph", in, out)
				}
			}
		}

	}

	if err := dep.kahn.Build(); err != nil {
		return errors.Wrap(err, "cant build graph")
	}

	for _, name := range dep.kahn.Result() {
		if item, ok := dep.all[name]; ok {

			if values, err := dep.di(item); err == nil {
				for _, value := range values {
					if value.Type().AssignableTo(srvType) {
						dep.srv.Add(value.Interface().(iserv))
					}

					dep.all[name] = value.Interface()
				}
			} else if !errors2.Is(err, errorBadAction) {
				return errors.Wrapf(err, "cant initialize %s", name)
			}

		} else {
			return errorDepUnknown
		}
	}

	return nil
}

// Down - stop all services in dependencies
func (dep *Dependencies) Down() error {
	return dep.srv.Down()
}

// Up - start all services in dependencies
func (dep *Dependencies) Up() error {
	return dep.srv.Up()
}

// Inject - obtained dependence
func (dep *Dependencies) Inject(item interface{}) error {
	_, err := dep.di(item)
	return err
}

func (dep *Dependencies) di(item interface{}) ([]reflect.Value, error) {
	ref := reflect.TypeOf(item)
	args := make([]reflect.Value, 0)

	switch ref.Kind() {

	case reflect.Func:
		for i := 0; i < ref.NumIn(); i++ {
			in := ref.In(i).String()
			if el, ok := dep.all[in]; ok {
				args = append(args, reflect.ValueOf(el))
			} else {
				return nil, errors.Wrap(errorDepUnknown, in)
			}
		}

	case reflect.Struct:
		value := reflect.New(ref)
		for i := 0; i < ref.NumField(); i++ {
			in := ref.Field(i).Type.String()
			if el, has := dep.all[in]; has {
				value.Elem().FieldByName(ref.Field(i).Name).Set(reflect.ValueOf(el))
			} else {
				return nil, errors.Wrap(errorDepUnknown, in)
			}
		}
		return append(args, value), nil

	default:
		return nil, errorBadAction
	}

	return reflect.ValueOf(item).Call(args), nil
}
