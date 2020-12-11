/**
 * Copyright 2020 Mikhail Knyazhev <markus621@gmail.com>. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

//Modules DI container
type Modules []interface{}

//Add object to container
func (m Modules) Add(v ...interface{}) Modules {
	for _, mod := range v {
		switch mod.(type) {
		case Modules:
			m = m.Add(mod.(Modules)...)
		default:
			m = append(m, mod)
		}
	}
	return m
}
