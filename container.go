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

// сontainer keeps components in form of
//
//	{
//	   Type1: {
//	     ID1: {..}
//	     ID2: {..}
//	   }
//	   Type2: {..}
//	}
type сontainer struct {
	mu          sync.Mutex
	initialized bool

	logger *zap.SugaredLogger

	initFns  []func(*сontainer) error             // keeps init functions
	stages   map[id]stage                         // keeps registered stages
	stageFns map[id][]func(context.Context) error // keeps stages function grouped by stage id

	components map[reflect.Type]map[id]*component

	dump *dump
}

type GlobalContainerOption interface {
	applyGlobalContainerOption(*сontainer)
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

var globC = &сontainer{
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
