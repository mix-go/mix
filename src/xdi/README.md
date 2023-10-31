> Produced by OpenMix: [https://openmix.org](https://openmix.org/mix-go)

## Mix XDI

DI, IoC container

## Overview

A library for creating objects and managing their dependencies. This library can be used for managing dependencies in a unified way, managing global objects, and refreshing dynamic configurations.

## Installation

```
go get github.com/mix-go/xdi
```

## Quick start

Create a singleton through dependency configuration

```go
package main

import (
    "github.com/mix-go/xdi"
)

type Foo struct {
}

func init() {
    obj := &xdi.Object{
        Name: "foo",
        New: func() (interface{}, error) {
            i := &Foo{}
            return i, nil
        },
    }
    if err := xdi.Provide(obj); err != nil {
        panic(err)
    }
}

func main() {
    var foo *Foo
    if err := xdi.Populate("foo", &foo); err != nil {
        panic(err)
    }
    // use foo
}
```

## Reference

Refer to another dependency configuration instance in the dependency configuration

```go
package main

import (
    "github.com/mix-go/xdi"
)

type Foo struct {
    Bar *Bar
}

type Bar struct {
}

func init() {
    objs := []*xdi.Object{
        {
            Name: "foo",
            New: func() (interface{}, error) {
                // reference bar
                var bar *Bar
                if err := xdi.Populate("bar", &bar); err != nil {
                    return nil, err
                }

                i := &Foo{
                    Bar: bar,
                }
                return i, nil
            },
        },
        {
            Name: "bar",
            New: func() (interface{}, error) {
                i := &Bar{}
                return i, nil
            },
            NewEverytime: true,
        },
    }
    if err := xdi.Provide(objs...); err != nil {
        panic(err)
    }
}

func main() {
    var foo *Foo
    if err := xdi.Populate("foo", &foo); err != nil {
        panic(err)
    }
    // use foo
}
```

## Refresh singleton

When the configuration information changes during program execution, `Refresh()` can refresh the singleton instance to switch to using the new configuration. It is commonly used in microservice configuration centers.

```go
obj, err := xdi.DefaultContainer.Object("foo")
if err != nil {
    panic(err)
}
if err := obj.Refresh(); err != nil {
    panic(err)
}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
