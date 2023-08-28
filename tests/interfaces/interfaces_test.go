package interfaces

import (
	"testing"

	"github.com/sabahtalateh/gic"
)

func init() {
	_ = gic.Init()
}

type Int interface {
	F()
}

type Impl struct{}

func (i Impl) F() {}

func TestCanNotKeepInterface(t *testing.T) {
	// err := gic.ComponentE[Int](gic.ID(), func() (Int, error) { return &Impl{}, nil })
	// require.ErrorIs(t, err, gic.ErrInterface)
}
