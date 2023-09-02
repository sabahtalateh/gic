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

There is also such feature as `stage` which allows us to execute some action on all the `components` implementing that `stage`. More details

## Add component

To add component use `gic.Add` from 

## Stage

### Add stage

### Implement stage

### Run stage

