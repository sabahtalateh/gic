package gic

import (
	"fmt"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	// TODO link to doc in errors texts.
	ErrInitialized        = fmt.Errorf("container already initialized")
	ErrNotInitialized     = fmt.Errorf("container not initialized. call gic.Init")
	ErrNoInit             = fmt.Errorf("no init function. use gic.WithInit or gic.WithInitE")
	ErrBothInit           = fmt.Errorf("both gic.Init and git.InitE used. use one")
	ErrComponentAdded     = fmt.Errorf("component added")
	ErrNotFound           = fmt.Errorf("not found")
	ErrInterface          = fmt.Errorf("container can not keep interfaces")
	ErrNotFromInit        = fmt.Errorf("component should be added from init function")
	ErrEmptyStageName     = fmt.Errorf("stage id must be set")
	ErrStageRegistered    = fmt.Errorf("stage already registered")
	ErrStageNotRegistered = fmt.Errorf("stage not registered")
)
