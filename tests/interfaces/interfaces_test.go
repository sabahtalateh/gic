package interfaces

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
)

type Iface interface {
	F()
}

type Impl struct{}

func (i Impl) F() {}

var initErr error

func init() {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				initErr = e
			}
		}
	}()
	gic.Add[Iface](gic.WithInit(func() Iface { return &Impl{} }))
}

func TestAddInterfacePanics(t *testing.T) {
	require.ErrorIs(t, initErr, gic.ErrInterface)
}
