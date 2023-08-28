package gic

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type stageImpl[T any] struct {
	s       stage
	onStage func(context.Context, T) error
}

type addOptions[T any] struct {
	id        id
	init      func() T
	initE     func() (T, error)
	stageImps map[id]func(context.Context, T) error
}

type addOption[T any] interface{ addOption() }

type withInitE[T any] struct{ f func() (T, error) }
type withInit[T any] struct{ f func() T }
type withStageImpl[T any] struct{ stageImpl[T] }

func (w withInitE[T]) addOption()     {}
func (w withInit[T]) addOption()      {}
func (w withStageImpl[T]) addOption() {}

func WithInitE[T any](f func() (T, error)) withInitE[T] { return withInitE[T]{f: f} }
func WithInit[T any](f func() T) withInit[T]            { return withInit[T]{f: f} }
func WithStageImpl[T any](s stage, onStage func(context.Context, T) error) withStageImpl[T] {
	return withStageImpl[T]{stageImpl[T]{s: s, onStage: onStage}}
}

// Add adds component init function into container
// Added functions will be run on Init call
// This function be called only from init() function so it panics on error
func Add[T any](opts ...addOption[T]) {
	err := add[T](globC, makeAddOptions[T](opts...))
	check(err)
}

func makeAddOptions[T any](opts ...addOption[T]) addOptions[T] {
	ao := addOptions[T]{stageImps: map[id]func(context.Context, T) error{}}

	// TODO check options repeat (2 calls to ic.Start for example)
	for _, opt := range opts {
		switch o := opt.(type) {
		case withID:
			ao.id = o.id
		case withInit[T]:
			ao.init = o.f
		case withInitE[T]:
			ao.initE = o.f
		case withStageImpl[T]:
			ao.stageImps[o.s.id] = o.onStage
		default:
			panic("unsupported option")
		}
	}

	return ao
}

func add[T any](c *Container, ao addOptions[T]) error {
	var (
		t    T
		err  error
		call = makeCaller()
	)

	if err = checkCallFromInit(call); err != nil {
		return err
	}

	// common approach to get type of interfaces
	// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
	rt := reflect.TypeOf(&t).Elem()
	if err = checkAddType(rt); err != nil {
		return err
	}

	if ao.init == nil && ao.initE == nil {
		return ErrNoInit
	}

	if ao.init != nil && ao.initE != nil {
		return ErrBothInit
	}

	c.initFns = append(c.initFns, func(*Container) error { return initFn[T](c, rt, ao, call) })

	return nil
}

func initFn[T any](c *Container, typ reflect.Type, ao addOptions[T], addCall *caller) error {
	if comps, ok := c.components[typ]; ok {
		if err := checkAdd(comps, ao.id); err != nil {
			return err
		}
	} else {
		c.components[typ] = map[id]*component{}
	}

	var (
		t   T
		err error
	)

	// reset got components before execute init function
	if c.dump != nil {
		c.dump.got = nil
	}

	if ao.init != nil {
		t = ao.init()
	}

	if ao.initE != nil {
		t, err = ao.initE()
		if err != nil {
			return err
		}
	}

	comp := &component{id: ao.id, caller: addCall, c: t}
	c.components[typ][ao.id] = comp

	logID := ao.id
	if logID == "" {
		logID = "(No ID)"
	}
	c.LogInfof("%s[id=%s] initialized", typ, logID)

	dumpComponent(c, comp)

	for stgID, f := range ao.stageImps {
		f := f
		stg, ok := c.stages[stgID]
		if !ok {
			return errors.Join(ErrStageNotRegistered, fmt.Errorf("%s", stgID))
		}
		c.stageFns[stgID] = append(c.stageFns[stgID], func(ctx context.Context) error { return f(ctx, t) })
		c.LogInfof("%s[id=%s] implementing stage: %s", typ, logID, stgID)

		dumpStageImpl(c, &stg, comp)
	}

	return nil
}

func checkAdd(comps map[id]*component, id id) error {
	// FORBIDDEN to have same id for different components
	if comp, ok := comps[id]; ok {
		return errIDInUse(id, comp.caller, makeCaller())
	}

	return nil
}

func checkAddType(t reflect.Type) error {
	if t.Kind() == reflect.Interface {
		return errors.Join(
			ErrInterface,
			fmt.Errorf("attempting to add %s interface\n%s", t, makeCaller()),
		)
	}

	return nil
}
