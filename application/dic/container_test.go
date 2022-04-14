package dic_test

import (
	"fmt"
	"testing"

	ctx2 "github.com/deweppro/go-app/application/ctx"
	"github.com/deweppro/go-app/application/dic"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

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
	T0  *t0
	T1  *t1
	T2  *t2
	T7  *t7
	T44 t44
}

type t44 struct {
	Env string
}

type t5 struct{}

func newT5() *t5         { return &t5{} }
func (t5 *t5) V() string { return "t5V" }

type t6 struct{ T4 t4 }

func newT6(t4 t4) *t6    { return &t6{T4: t4} }
func (t6 *t6) V() string { return "t6V" }

type t7 struct{}

func newT7() *t7         { return &t7{} }
func (t7 *t7) V() string { return "t7V" }

type t8 struct{}

func newT8() (*t8, error) { return &t8{}, nil }
func (t8 *t8) V() string  { return "t8V" }

type hello string

var a = hello("hhhh")

type ii interface {
	V() string
}

func newT7i(_ hello) ii {
	return &t7{}
}

func TestUnit_Dependencies(t *testing.T) {
	ctx := ctx2.New()
	dep := dic.New()

	require.NoError(t, dep.Register([]interface{}{
		newT1, newT2, newT5, newT6, newT7(), newT8,
		a, newT7i, newT0, t44{Env: "aaa"},
	}...))

	require.NoError(t, dep.Build())
	require.NoError(t, dep.Up(ctx))
	require.Error(t, dep.Up(ctx))

	require.NoError(t, dep.Inject(func(a *t6, b ii, c hello, d *t8) {
		require.Equal(t, "t6V", a.V())
		require.Equal(t, "t0V", a.T4.T0.V())
		require.Equal(t, "t1V", a.T4.T1.V())
		require.Equal(t, "t7V", b.V())
		require.Equal(t, "t8V", d.V())
		require.Equal(t, hello("hhhh"), c)
	}))

	require.Error(t, dep.Inject(func(a string, b int, c bool) {

	}))

	require.NoError(t, dep.Down(ctx))
	require.Error(t, dep.Down(ctx))
}

type demo1 struct{}
type demo2 struct{}
type demo3 struct{}

func newDemo() (*demo1, *demo2, *demo3) { return &demo1{}, &demo2{}, &demo3{} }
func (d *demo1) Up() error {
	fmt.Println("demo1 up")
	return nil
}
func (d *demo1) Down() error {
	fmt.Println("demo1 down")
	return nil
}

func TestUnit_Dependencies2(t *testing.T) {
	ctx := ctx2.New()
	dep := dic.New()
	require.NoError(t, dep.Register([]interface{}{
		newDemo,
	}...))
	require.NoError(t, dep.Build())
	require.NoError(t, dep.Up(ctx))
	require.Error(t, dep.Up(ctx))
	require.NoError(t, dep.Down(ctx))
	require.Error(t, dep.Down(ctx))
}

type demo4 struct{}

func newDemo4() (*demo4, error) { return nil, fmt.Errorf("fail init constructor demo4") }

func TestUnit_Dependencies3(t *testing.T) {
	dep := dic.New()
	require.NoError(t, dep.Register([]interface{}{
		newDemo4,
	}...))
	err := dep.Build()
	require.Error(t, err)
	fmt.Println(err.Error())
	require.Contains(t, err.Error(), "initialize error <github.com/deweppro/go-app/application/dic_test.demo4>: fail init constructor demo4")
}
