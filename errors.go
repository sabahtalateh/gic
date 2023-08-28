package gic

import (
	"errors"
	"fmt"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	// TODO link to doc in errors texts
	ErrInitialized        = fmt.Errorf("already initialized")
	ErrNotInitialized     = fmt.Errorf("not initialized. call gic.Init")
	ErrNoInit             = fmt.Errorf("no init function. use gic.WithInit or gic.WithInitE")
	ErrBothInit           = fmt.Errorf("both Init and InitE set. use one")
	ErrIDInUse            = fmt.Errorf("component id in use")
	ErrNotFound           = fmt.Errorf("not found")
	ErrInterface          = fmt.Errorf("container can not keep interfaces")
	ErrNotFromInit        = fmt.Errorf("component should be added from init function")
	ErrEmptyStageName     = fmt.Errorf("stage id must be set")
	ErrStageRegistered    = fmt.Errorf("stage already registered")
	ErrStageNotRegistered = fmt.Errorf("stage not registered")
)

func errIDInUse(id id, where *caller, c *caller) error {
	return errors.Join(ErrIDInUse, fmt.Errorf("component id %d in use\n%s\n%s", id, where, c))
}
