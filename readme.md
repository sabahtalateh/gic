# GIC - Global Init Container

- `(G)` It stored in a [global variable](https://github.com/sabahtalateh/gic/blob/main/container.go#L60) `globC`
- `(I)` Components added in `init` function
- `(C)` ontainer

## Concept

`Golang` has `init` mechanism. Package's `init` functions automatically called for all the project in hierarchy. `init` mechanism also solves dependencies cycling problem as project will not be compiled if cycles exists

Let's use this mechanism to compose dependency injection container

We will add `struct instances` (A.K.A. `components`) into container within `init` function with `gic.Add`. Components will be added at the same file the structs are defined. It will take us to the point where we have all the `components` of some package by just importing it (which is cool in my opinion :-)). Components will not be created yet but only initialization functions passed to `gic.Add` will be added to [their array](https://github.com/sabahtalateh/gic/blob/main/add.go#L130)

Then in `main` function (or whatever function runs on your application startup) we call `gic.Init` to [call functions](https://github.com/sabahtalateh/gic/blob/main/init.go#L12)

That's all, we have initialized global container and can use `gic.Get/gic.GetE` to retrieve component

There is also such feature as `stage` which allows us to execute some action on all the `components` implementing that `stage`. [Read more](https://github.com/sabahtalateh/gic#stages)

## Components

### Add

To add component use `gic.Add` from `init` function. It checked with `runtime.Caller`.

```go
package services

import (
	"github.com/sabahtalateh/gic"
)

type SomeService struct {
	param string
}

func (s *SomeService) Hello() {
	println("Hello " + s.param)
}

func init() {
	gic.Add[*SomeService](
		gic.WithInit(func() *SomeService {
			return &SomeService{param: "World!"}
		}),
	)
}

```

### Init

Then in your `main` function initialize global container

```go
package main

import (
	"log"
	
	"github.com/sabahtalateh/gic"
)

func main() {
	err := gic.Init()
	if err != nil {
		log.Fatal(err)
	}
}
```

### Get

Now we can get our component from container

```go
package main

import (
	"log"
	
	"github.com/sabahtalateh/gic"
)

func main() {
	// ...
	srv, err := gic.GetE[*SomeService]()
	if err != nil {
		log.Fatal(err)
	}
	srv.Hello()
}
```
```shell
Hello World!
```

### ID

To create an instance of one struct several times with different parameters `gic.ID` is used

```go
import (
	"github.com/sabahtalateh/gic"
)

type DB struct {
	dsn string
}

var Write = gic.ID("WriteDB")
var Read = gic.ID("ReadDB")

func init() {
	gic.Add[*DB](
		gic.WithID(Write),
		gic.WithInit(func() *DB { return &DB{dsn: "write dsn"} }),
	)

	gic.Add[*DB](
		gic.WithID(Read),
		gic.WithInit(func() *DB { return &DB{dsn: "read dsn"} }),
	)
}
```
Previously created ids also used to get components
```go
type Service struct {
	writeDB *DB
	readDB  *DB
}

func init() {
	gic.Add[*Service](
		gic.WithInit(func() *Service {
			return &Service{
				writeDB: gic.Get[*DB](gic.WithID(Write)),
				readDB:  gic.Get[*DB](gic.WithID(Read)),
			}
		}),
	)
}
```

## Stages

### Add stage

### Implement stage

### Run stage

