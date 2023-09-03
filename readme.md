# GIC - Global Init Container

- `(G)` It stored in a [global variable](https://github.com/sabahtalateh/gic/blob/main/container.go#L60) `globC`
- `(I)` Components added in `init` function
- `(C)` ontainer

## Concept

`Golang` has `init` mechanism. Package's `init` functions automatically called in hierarchy. `init` mechanism also solves dependencies cycling problem as project will not be compiled if cycles exists

Let's use this mechanism to compose dependency injection container

We will add `struct instances` (A.K.A. `components`) into container within `init` function with `gic.Add`. Components will be added at the same file the structs are defined. It will take us to the point where we have all the `components` of some package by just importing it. Components will not be created on `gic.Add` but only initialization functions passed to `gic.Add` will be added to [their array](https://github.com/sabahtalateh/gic/blob/main/add.go#L130)

Then in program entry point we call `gic.Init` to [create components](https://github.com/sabahtalateh/gic/blob/main/init.go#L12) with function passed to `gic.Add` 

That's all, we have initialized global container and can use `gic.Get/gic.GetE` to retrieve component

There is also such feature as `stage` which allows us to execute some action on all the `components` implementing `stage`. [Read more](https://github.com/sabahtalateh/gic#stages)

## Components

### Add

To add component use `gic.Add` from `init` function. If called from another function `gic.Add` will panic. Checked with `runtime.Caller`.

```go
package internal

import "github.com/sabahtalateh/gic"

type Greeter struct {
	greet string
}

func (g *Greeter) Greet(whom string) string {
	return fmt.Sprintf("%s %s!", g.greet, whom)
}

func init() {
	gic.Add[*Greeter](
		gic.WithInit(func() *Greeter { return &Greeter{greet: "Hello"} }),
	)
}
```

### Init

Then in your `main` function initialize global container. As said above `gic.Init()` will call initialization functions added in `init` functions. Functions will be called in adding order which is equal to `init`-calls order

```go
package main

func main() {
	err := gic.Init()
	if err != nil {
	    panic(err)
	}
}
```

### Get

Now we can get our component from container with `gic.Get` or `gic.GetE`. Provide component type you want to get and component id if it was added with `gic.WithID` id (see: https://github.com/sabahtalateh/gic#id)

```go
package main

func main() {
	// ...
	
	g, err := gic.Get[*internal.Greeter]()
	if err != nil {
		panic(err)
	}
	println(g.Greet("World"))
}
```
```shell
Hello World!
```

### ID

To create an instances of one struct with different parameters `gic.ID` is used. After creation use `gic.WithID` to get component by ID

```go
package internal

import (
	"fmt"
	
	"github.com/sabahtalateh/gic"
)

type Greeter struct {
	greet string
}

func (g *Greeter) Greet(whom string) string {
	return fmt.Sprintf("%s %s!", g.greet, whom)
}

var RussianGreeter = gic.ID("RussianGreeter")
var ChineseGreeter = gic.ID("ChineseGreeter")

func init() {
	gic.Add[*Greeter](
		gic.WithID(RussianGreeter),
		gic.WithInit(func() *Greeter { return &Greeter{greet: "Привет"} }),
	)

	gic.Add[*Greeter](
		gic.WithID(ChineseGreeter),
		gic.WithInit(func() *Greeter { return &Greeter{greet: "你好"} }),
	)
}
```

```go
func main() {
	// ...
	
	g, err = gic.GetE[*internal.Greeter](gic.WithID(internal.RussianGreeter))

	g, err = gic.GetE[*internal.Greeter](gic.WithID(internal.ChineseGreeter))
	
}
```


## Stages

### Add stage

### Implement stage

### Run stage

