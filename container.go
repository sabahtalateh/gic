package gic

import (
	"context"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

type component struct {
	id     id
	caller *caller
	c      any
}

// Container keeps components in form of
//
//	{
//	   Type1: {
//	     ID1: {..}
//	     ID2: {..}
//	   }
//	   Type2: {..}
//	}
type Container struct {
	mu          sync.Mutex
	initialized bool

	logger *zap.SugaredLogger

	initFns  []func(*Container) error             // keeps init functions
	stages   map[id]stage                         // keeps registered stages
	stageFns map[id][]func(context.Context) error // keeps stages function grouped by stage id

	components map[reflect.Type]map[id]*component

	dump *dump
}

type GlobalContainerOption interface {
	applyGlobalContainerOption(*Container)
}

func ConfigureGlobalContainer(opts ...GlobalContainerOption) error {
	globC.mu.Lock()
	defer globC.mu.Unlock()
	if globC.initialized {
		return ErrInitialized
	}

	for _, opt := range opts {
		opt.applyGlobalContainerOption(globC)
	}

	return nil
}

var globC = &Container{
	stages:     map[id]stage{},
	stageFns:   map[id][]func(context.Context) error{},
	components: map[reflect.Type]map[id]*component{},
}

var start = RegisterStage(WithID("start"))
var stop = RegisterStage(WithID("stop"))

// WithStart adds component start function
// Added function will be executed on Start call
func WithStart[T any](f func(ctx context.Context, t T) error) withStageImpl[T] {
	return WithStageImpl[T](start, f)
}

// WithStop adds component stop function
// Added function will be executed on Stop call
func WithStop[T any](f func(ctx context.Context, t T) error) withStageImpl[T] {
	return WithStageImpl[T](stop, f)
}

// Start executes all functions added with WithStart
func Start(ctx context.Context) error {
	return RunStage(ctx, start)
}

// Stop executes all functions added with WithStop
func Stop(ctx context.Context) error {
	return RunStage(ctx, stop)
}
