package simple

import (
	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
	"testing"
)

type Component4 struct {
	prop string
}

// not added to container
type Component5 struct {
	prop string
}

// component added without this id
var SomeComponent4 = gic.ID("SomeComponent4")

func init() {
	gic.Add[*Component4](
		gic.WithInit(func() *Component4 {
			return &Component4{prop: "value"}
		}),
	)
}

func TestComponentNotFoundWithoutID(t *testing.T) {
	_ = gic.Init()

	_, err := gic.GetE[*Component5]()
	require.ErrorIs(t, err, gic.ErrNotFound)

	_, err = gic.GetE[*Component4](gic.WithID(SomeComponent4))
	require.ErrorIs(t, err, gic.ErrNotFound)
}
