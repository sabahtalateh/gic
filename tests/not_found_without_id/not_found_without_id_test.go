package not_found_without_id

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
)

type Component struct {
	prop string
}

// not added to container
type Component2 struct {
	prop string
}

// component added without this id
var SomeComponent = gic.ID("SomeComponent")

func init() {
	gic.Add[*Component](
		gic.WithInit(func() *Component {
			return &Component{prop: "value"}
		}),
	)
}

func TestComponentNotFoundWithoutID(t *testing.T) {
	var err error

	_ = gic.Init()

	_, err = gic.GetE[*Component](gic.WithID(SomeComponent))
	require.ErrorIs(t, err, gic.ErrNotFound)

	_, err = gic.GetE[*Component2]()
	require.ErrorIs(t, err, gic.ErrNotFound)

	_, err = gic.GetE[*Component2](gic.WithID(SomeComponent))
	require.ErrorIs(t, err, gic.ErrNotFound)
}
