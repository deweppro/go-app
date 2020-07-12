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
	"testing"

	"github.com/stretchr/testify/require"
)

type t0 struct{}

func newT0() *t0           { return &t0{} }
func (t0 *t0) Up() error   { return nil }
func (t0 *t0) Down() error { return nil }
func (t0 *t0) V() string   { return "t0V" }

type t1 struct {
	t0 *t0
}

func newT1(t0 *t0) *t1     { return &t1{t0: t0} }
func (t1 *t1) Up() error   { return nil }
func (t1 *t1) Down() error { return nil }
func (t1 *t1) V() string   { return "t1V" }

type t2 struct {
	t0 *t0
	t1 *t1
}

func newT2(t1 *t1, t0 *t0) *t2             { return &t2{t0: t0, t1: t1} }
func (t2 *t2) Up() error                   { return nil }
func (t2 *t2) Down() error                 { return nil }
func (t2 *t2) V() (string, string, string) { return "t2V", t2.t1.V(), t2.t0.V() }

type t4 struct {
	T0 *t0
	T1 *t1
	T2 *t2
	T7 *t7
}

type t5 struct{}

func newT5() *t5         { return &t5{} }
func (t5 *t5) V() string { return "t5V" }

type t6 struct{ T4 *t4 }

func newT6(t4 *t4) *t6   { return &t6{T4: t4} }
func (t6 *t6) V() string { return "t6V" }

type t7 struct{}

func newT7() *t7         { return &t7{} }
func (t7 *t7) V() string { return "t7V" }

type hello string

type ii interface {
	V() string
}

func newT7i() ii { return &t7{} }

func TestUnit_Dependencies(t *testing.T) {
	dep := New()

	var a hello
	a = "hhhh"
	b := "gggg"

	require.NoError(t, dep.Register([]interface{}{
		newT1, newT2, newT5, newT6, t4{}, newT7(),
		1, "hello", true, a, b, newT7i, newT0,
	}))

	require.NoError(t, dep.Build())

	require.NoError(t, dep.Up())
	require.Error(t, dep.Up())

	require.NoError(t, dep.Inject(func(a *t6, b ii, c hello) {
		require.Equal(t, "t6V", a.V())
		require.Equal(t, "t0V", a.T4.T0.V())
		require.Equal(t, "t1V", a.T4.T1.V())
		require.Equal(t, "t7V", b.V())
		require.Equal(t, hello("hhhh"), c)
	}))

	require.Error(t, dep.Inject(func(a string, b int, c bool) {

	}))

	require.NoError(t, dep.Down())
	require.Error(t, dep.Down())
}
