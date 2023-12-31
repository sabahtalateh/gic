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

type addOption[T any] interface {
	addOption()
	addOptionCallInfo() string
}

type withInitE[T any] struct {
	f func() (T, error)
	c caller
}
type withInit[T any] struct {
	f func() T
	c caller
}
type withStageImpl[T any] struct {
	impl stageImpl[T]
	c    caller
}

func (w withInitE[T]) addOption()     {}
func (w withInit[T]) addOption()      {}
func (w withStageImpl[T]) addOption() {}

func (w withInitE[T]) addOptionCallInfo() string {
	var t T
	return fmt.Sprintf("gic.WithInitE[%s]\n%s", reflect.TypeOf(t), w.c)
}
func (w withInit[T]) addOptionCallInfo() string {
	var t T
	return fmt.Sprintf("gic.WithInit[%s]\n%s", reflect.TypeOf(t), w.c)
}
func (w withStageImpl[T]) addOptionCallInfo() string {
	var t T
	return fmt.Sprintf("gic.WithStageImpl[%s]\n%s", reflect.TypeOf(t), w.c)
}

// WithInitE set component initialization function which can return error.
func WithInitE[T any](f func() (T, error)) withInitE[T] { return withInitE[T]{f: f, c: makeCaller()} }

// WithInit set component initialization function.
func WithInit[T any](f func() T) withInit[T] { return withInit[T]{f: f, c: makeCaller()} }

// WithStageImpl set stage implementation function.
func WithStageImpl[T any](s stage, onStage func(context.Context, T) error) withStageImpl[T] {
	return withStageImpl[T]{impl: stageImpl[T]{s: s, onStage: onStage}, c: makeCaller()}
}

// Add adds component init function into container
// Added functions will be run on Init call
// This function to be called only from init() function so it panics on error.
func Add[T any](opts ...addOption[T]) {
	err := add[T](globC, makeAddOptions[T](opts...))
	check(err)
}

func makeAddOptions[T any](opts ...addOption[T]) addOptions[T] {
	ao := addOptions[T]{stageImps: map[id]func(context.Context, T) error{}}

	// TODO check options repeat (2 calls to ic.Start for example).
	for _, opt := range opts {
		switch o := opt.(type) {
		case withID:
			ao.id = o.id
		case withInit[T]:
			ao.init = o.f
		case withInitE[T]:
			ao.initE = o.f
		case withStageImpl[T]:
			ao.stageImps[o.impl.s.id] = o.impl.onStage
		default:
			var t T
			panic(fmt.Sprintf(
				"inconsistent type parameters on gic.Add[%s] and on %s",
				reflect.TypeOf(t), opt.addOptionCallInfo(),
			))
		}
	}

	return ao
}

func add[T any](c *container, ao addOptions[T]) error {
	var (
		t    T
		err  error
		call = makeCaller()
	)

	// no need to be locked as can be used only from init
	if c.initialized {
		return ErrInitialized
	}

	if call.found {
		if err = checkCallFromInit(call); err != nil {
			return err
		}
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

	if err = checkAdd[T](c.initLocs, rt, ao.id); err != nil {
		return err
	}

	fn := func(*container) error { return initFn[T](c, rt, ao, call) }
	c.initFns = append(c.initFns, fn)
	c.setInitLocation(rt, ao.id.v, len(c.initFns)-1, call)

	return nil
}

func initFn[T any](c *container, typ reflect.Type, ao addOptions[T], caller caller) error {
	if _, ok := c.components[typ]; !ok {
		c.components[typ] = map[string]*component{}
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

	comp := &component{id: ao.id, caller: caller, c: t}
	c.components[typ][ao.id.v] = comp

	c.LogInfof("%s%s initialized", typ, strID(ao.id.v))

	dumpComponent(c, comp)

	for stgID, f := range ao.stageImps {
		f := f
		stg, ok := c.stages[stgID.v]
		if !ok {
			return errors.Join(ErrStageNotRegistered, fmt.Errorf("%s", stgID))
		}
		c.stageFns[stgID.v] = append(c.stageFns[stgID.v], func(ctx context.Context) error { return f(ctx, t) })
		c.LogInfof("%s%s implementing stage: %s", typ, strID, strID(ao.id.v))

		dumpStageImpl(c, &stg, comp)
	}

	return nil
}

func checkAdd[T any](locs map[reflect.Type]map[string]initLocation, typ reflect.Type, id id) error {
	ids, ok := locs[typ]
	if !ok {
		return nil
	}

	var t T

	// FORBIDDEN to have same id for same type
	if comp, ok := ids[id.v]; ok {
		idStr := " "
		if id.v != "" {
			idStr = fmt.Sprintf(" with id %s ", id.v)
		}

		return errors.Join(
			ErrComponentIDInUse,
			fmt.Errorf(
				"component of type %s%salready added\n%s\n%s", reflect.TypeOf(t), idStr, comp.caller, makeCaller(),
			),
		)
	}

	return nil
}

func checkAddType(t reflect.Type) error {
	if t.Kind() == reflect.Interface {
		return errors.Join(ErrInterface, fmt.Errorf("attempting to add %s interface\n%s", t, makeCaller()))
	}

	return nil
}
