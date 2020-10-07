/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"github.com/pkg/errors"
)

var (
	ErrDepRunning    = errors.New("dependencies is already running")
	ErrDepNotRunning = errors.New("dependencies are not running yet")
	ErrDepEmpty      = errors.New("dependencies is empty")
	ErrDepUnknown    = errors.New("unknown dependency")
	ErrBadAction     = errors.New("is not a supported action")
)

var types = []string{
	"invalid",
	"bool",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uintptr",
	"float32",
	"float64",
	"complex64",
	"complex128",
	"array",
	"chan",
	"func",
	"interface",
	"map",
	"ptr",
	"slice",
	"string",
	"struct",
	"unsafe.Pointer",
}

func isDefaultType(name string) bool {
	for _, el := range types {
		if el == name {
			return true
		}
	}
	return false
}
