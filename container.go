package gic

import (
	"context"
	"reflect"
	"sync"

	"go.uber.org/zap"
)

type component struct {
	id     id
	caller caller
	c      any
}

// container keeps components in form of
//
//	{
//	   Type1: {
//	     ID1: {..}
//	     ID2: {..}
//	   }
//	   Type2: {..}
//	}
type container struct {
	mu          sync.Mutex
	initialized bool

	logger *zap.SugaredLogger

	initFns   []func(*container) error                 // init functions in add order
	initsDone map[int]struct{}                         // indexes of executed initFns
	initLocs  map[reflect.Type]map[string]initLocation // initFns indices by component type and id

	stages   map[string]stage                         // registered stages
	stageFns map[string][]func(context.Context) error // stages function grouped by stage id

	components map[reflect.Type]map[string]*component

	dump *dump
}

var globC = &container{
	initsDone:  map[int]struct{}{},
	initLocs:   map[reflect.Type]map[string]initLocation{},
	stages:     map[string]stage{},
	stageFns:   map[string][]func(context.Context) error{},
	components: map[reflect.Type]map[string]*component{},
}

// container has 2 predefined stages.
var start = RegisterStage(WithID(ID("Start")))
var stop = RegisterStage(WithID(ID("Stop")))

type GlobalContainerOption interface {
	applyGlobalContainerOption(*container)
}

// ConfigureGlobalContainer configures dump (see: WithDump) and logger (see: WithZapSugaredLogger).
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

// WithStart adds component start function
// Added function will be executed on Start call.
func WithStart[T any](f func(ctx context.Context, t T) error) withStageImpl[T] {
	return WithStageImpl[T](start, f)
}

// WithStop adds component stop function
// Added function will be executed on Stop call.
func WithStop[T any](f func(ctx context.Context, t T) error) withStageImpl[T] {
	return WithStageImpl[T](stop, f)
}

// Start executes all functions added with WithStart.
func Start(ctx context.Context) error {
	return RunStage(ctx, start)
}

// Stop executes all functions added with WithStop.
func Stop(ctx context.Context) error {
	return RunStage(ctx, stop)
}
