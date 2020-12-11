# go-app

[![Coverage Status](https://coveralls.io/repos/github/deweppro/go-app/badge.svg?branch=master)](https://coveralls.io/github/deweppro/go-app?branch=master)
[![Release](https://img.shields.io/github/release/deweppro/go-app.svg?style=flat-square)](https://github.com/deweppro/go-app/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/deweppro/go-app)](https://goreportcard.com/report/github.com/deweppro/go-app)
[![Build Status](https://travis-ci.com/deweppro/go-app.svg?branch=master)](https://travis-ci.com/deweppro/go-app)

## simple

***config.yaml***

```yaml
env: dev
log: /var/log/simple.log
pig: /var/run/simple.pid
```

***main.go***

```go
package main

import (
	"fmt"

	"github.com/deweppro/go-app"
)

var _ app.Servicer = (*Simple)(nil)

type (
	//Simple model
	Simple struct{}
	//SimpleConfig config model
	SimpleConfig struct {
		Env string `yaml:"env"`
	}
)

//NewSimple init Simple
func NewSimple(_ *SimpleConfig) *Simple {
	fmt.Println("call NewSimple")
	return &Simple{}
}

//Up  method for start Simple in DI container
func (s *Simple) Up() error {
	fmt.Println("call *Simple.Up")
	return nil
}

//Down  method for stop Simple in DI container
func (s *Simple) Down() error {
	fmt.Println("call *Simple.Down")
	return nil
}

func main() {
	app.New().
		ConfigFile(
			"./config.yaml",
			&SimpleConfig{},
		).
		Modules(
			NewSimple,
		).
		Run()
}
```

## HowTo

***Run the app***
```go
app.New()
    .ConfigFile(<path to config file: string>, <config objects separate by comma: ...interface{}>)
    .Modules(<config objects separate by comma: ...interface{}>)
    .Run()
```

***Supported types for initialization***

* Function that returns an object or interface

*All incoming dependencies will be injected automatically*
```go
type Simple1 struct{}
func NewSimple1(_ *logger.Logger) *Simple1 { return &Simple1{} }
```

*Returns the interface*
```go
type Simple2 struct{}
type Simple2Interface interface{
    Get() string
}
func NewSimple2() Simple2Interface { return &Simple2{} }
func (s2 *Simple2) Get() string { 
    return "Hello world"
}
```

*If the object has the `Up() error` and `Down() error` methods, they will be called `Up() error`  when the app starts, and `Down() error` when it finishes. This allows you to automatically start and stop routine processes inside the module*

```go
var _ app.Servicer = (*Simple3)(nil)
type Simple3 struct{}
func NewSimple3(_ *Simple4) *Simple3 { return &Simple3{} }
func (s3 *Simple3) Up() error { return nil }
func (s3 *Simple3) Down() error { return nil }
```

* Named type

```go
type HelloWorld string
```

* Object structure

```go
type Simple4 struct{
    S1 *Simple1
    S2 Simple2Interface
    HW HelloWorld
}
```

* Object reference or type

```go
s1 := &Simple1{}
hw := HelloWorld("Hello!!")
```


## Example of initialization of all types

```go
func main() {
	
    s1 := &Simple1{}
    hw := HelloWorld("Hello!!")

    app.New().
        ConfigFile(
            "config.yaml",
            &debug.ConfigDebug{},
        ).
        Modules(
            debug.New,
            NewSimple2,
            NewSimple3,
            Simple4{}
            s1, hw,
        ).
        Run()
}
```
