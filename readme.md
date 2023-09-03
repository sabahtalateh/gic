# GIC - Global Init Container

- `(G)` It stored in a [global variable](https://github.com/sabahtalateh/gic/blob/main/container.go#L60) `globC`
- `(I)` Components added in `init` function
- `(C)` ontainer

## Example

Example project with Postgres DB can be found here https://github.com/sabahtalateh/gicex

## Concept

`Golang` has `init` mechanism. Package's `init` functions automatically called in hierarchy. `init` mechanism solves two problems (at least): 
- Dependencies cycling problem. As project will not be compiled if cycles exists
- Unused components will not be included into container. As they never imported

**NOTE** This last point is a bit tricky as we need to retrieve `component` from `container` somewhere in program (or just import `component's` package) so it to be initialized

Let's use `init` mechanism to compose dependency injection container

We will add `struct instances` (A.K.A. `components`) into container within `init` function with `gic.Add`. Components will be
added at the same file the structs are defined. It will take us to the point where we have all the `components` of some package by
just importing it. Components will not be created on `gic.Add` but only initialization functions passed to `gic.Add` will be added
to [their array](https://github.com/sabahtalateh/gic/blob/main/add.go#L130)

Then in program entry point we call `gic.Init` to [create components](https://github.com/sabahtalateh/gic/blob/main/init.go#L12)
with function passed to `gic.Add`

That's all, we have initialized global container and can use `gic.Get/gic.GetE` to retrieve component

There is also such feature as `stage` which allows us to execute some action on all the `components`
implementing `stage`. [Read more](https://github.com/sabahtalateh/gic#stages)

## Components

### Add

To add component use `gic.Add` from `init` function. If called from another function `gic.Add` will panic. Checked
with `runtime.Caller`.

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/internal/greeter.go)

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

Then in your app entry point initialize global container. As said above `gic.Init()` will call initialization functions added
in `init` functions. Functions will be called in adding order which is equal to `init`-calls order

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/example_test.go)

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

Now we can get our component from container with `gic.Get` or `gic.GetE` (`gic.Get` will panic on errors). Provide component type
you want to get and component id if it was added with `gic.WithID` id (see: https://github.com/sabahtalateh/gic#id)

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/example_test.go)

```go
package main

import "github.com/sabahtalateh/gic/tests/example/internal"

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

To be able to have instances of one struct with different parameters (or just multiple times) `gic.ID` used. After `gic.Init` use `gic.GetE(gic.WithID..)` to get
component by ID

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/internal/greeter.go)

```go
package internal

import (
	"fmt"

	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/example/internal"
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

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/example_test.go)

```go
func main() {
// ...

	g, err = gic.GetE[*internal.Greeter](gic.WithID(internal.RussianGreeter))

	g, err = gic.GetE[*internal.Greeter](gic.WithID(internal.ChineseGreeter))

}
```

## Stages

Container has two predefined `stages`: `Start` and `Stop`. It can be useful for opening/closing db client sockets,
starting/stopping event consumers and workers and so on. To implement `Start` or `Stop` for some component use `gic.WithStart`
and `gic.WithStop`. Pass function accepting `context.Context` and `component`. Context can be set from outside to stop `stage` execution by timeout

Custom `stages` also can be registered 

### Implementing Start & Stop

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/internal/numberseater.go)

```go
package internal

import (
	"context"
	"sync"

	"github.com/sabahtalateh/gic"
)

type NumbersEater struct {
	mu    sync.Mutex
	c     chan int
	eaten []int
}

func (n *NumbersEater) Start() {
	go func() {
		for number := range n.c {
			n.mu.Lock()
			n.eaten = append(n.eaten, number)
			n.mu.Unlock()
		}
	}()
}

func (n *NumbersEater) Stop() {
	close(n.c)
}

func (n *NumbersEater) Feed(num int) {
	n.c <- num
	return
}

func (n *NumbersEater) Eaten() []int {
	return n.eaten
}

func init() {
	gic.Add[*NumbersEater](
		gic.WithInit(func() *NumbersEater {
			return &NumbersEater{c: make(chan int)}
		}),
		// Implement Start
		gic.WithStart(func(_ context.Context, ne *NumbersEater) error {
			ne.Start()
			return nil
		}),
		// Implement Stop
		gic.WithStop(func(ctx context.Context, ne *NumbersEater) error {
			ne.Stop()
			return nil
		}),
	)
}
```

### Running Start & Stop

`Start` and `Stop` `stages` runs manually after `gic.Init`

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/example_test.go)
```go
func main() {
	// ...
	err = gic.Init()
	// ...
	
	// Control Start timeout with context
	err = gic.Start(context.Background())
	// ...

	ne, err := gic.GetE[*internal.NumbersEater]()
	// ...
	
	ne.Feed(1)
	ne.Feed(2)
	// ...

	// Control Stop timeout with context
	err = gic.Start(context.Background())
	// ...
}
```

### Add custom stage

To add custom `stage` manually use `gic.RegisterStage`

Stage configuration options:
- `stage` should have unique ID (`gic.WithID`)
- By default `stage` will be run in `parallel` on all implementing components (see: https://github.com/sabahtalateh/gic/blob/main/stage.go). May be disabled with `gic.WithDisableParallel`
- By default `stage` will be run without order. May be changed with `gic.WithInitOrder` (same order as components were initialized) and `gic.WithReverseInitOrder` (reverse  to initialization order). 

**NOTE** order will not take effect until `parallel` disabled with `gic.WithDisableParallel`

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/internal/mystage.go)
```go
var MyStage = gic.RegisterStage(
	gic.WithID(gic.ID("MyStage")),
	gic.WithDisableParallel(),
	gic.WithInitOrder(),
)
```

### Implement custom stage

To implement stage use `gic.WithStageImpl`

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/internal/mystage.go)
```go
type Dummy struct {
	X int
}

func init() {
	gic.Add[*Dummy](
		gic.WithInit(func() *Dummy {
			return &Dummy{}
		}),
		gic.WithStageImpl(MyStage, func(ctx context.Context, d *Dummy) error {
			d.X = 999
			return nil
		}),
	)
}
```

### Run custom stage

To run custom stage use `gic.RunStage` after `gic.Init`

(see: https://github.com/sabahtalateh/gic/blob/main/tests/example/example_test.go)
```go
func main() {
	_ = gic.Init()
	// ...
	err = gic.RunStage(context.Background(), internal.MyStage)
	// ...   
}
```
