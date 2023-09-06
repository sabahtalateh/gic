package gic

import (
	"errors"
	"fmt"
	"reflect"
)

type getOpts struct{ id id }
type getOpt interface{ applyGetOption(*getOpts) }

// Get return component from container
// errors: ErrNotFound
func Get[T any](opts ...getOpt) (T, error) {
	oo := getOpts{}
	for _, opt := range opts {
		opt.applyGetOption(&oo)
	}

	return get[T](globC, oo)
}

// MustGet return component from container. panics on Get error
// errors: ErrNotFound
func MustGet[T any](opts ...getOpt) T {
	t, err := Get[T](opts...)
	check(err)
	return t
}

// errors: ErrNotFound
func get[T any](c *container, opts getOpts) (t T, err error) {
	lookFor := reflect.TypeOf(&t).Elem()
	if err = checkGetType(lookFor); err != nil {
		return t, err
	}

	initLoc, err := c.initLocation(lookFor, opts.id.v)
	if err != nil {
		return t, err
	}

	// If init function for component exists but was not executed yet then execute it
	if _, done := c.initsDone[initLoc.fn]; !done {
		if len(c.initFns)-1 < initLoc.fn {
			return t, fmt.Errorf("impossible")
		}
		if err = c.initFns[initLoc.fn](c); err != nil {
			return t, err
		}
		c.initsDone[initLoc.fn] = struct{}{}
	}

	comps, ok := c.components[lookFor]
	if !ok {
		return t, errNotFound(lookFor, opts.id.v)
	}

	comp, ok := comps[opts.id.v]
	if !ok {
		return t, errNotFound(lookFor, opts.id.v)
	}

	if c.dump != nil && !c.initialized {
		c.dump.got = append(c.dump.got, comp)
	}
	return comp.c.(T), nil
}

func checkGetType(t reflect.Type) error {
	if t.Kind() == reflect.Interface {
		return errors.Join(
			ErrInterface,
			fmt.Errorf("attempting to get %s interface\n%s", t, makeCaller()),
		)
	}

	return nil
}
