/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

type modules struct {
	data []interface{}
}

func newModules() *modules {
	return &modules{
		data: []interface{}{},
	}
}

func (m *modules) Add(a ...interface{}) *modules {
	m.data = append(m.data, a...)
	return m
}

func (m *modules) Get() []interface{} {
	return m.data
}
