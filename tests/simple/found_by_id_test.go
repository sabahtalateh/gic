package simple

import (
	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
	"testing"
)

type Component struct{}

var SomeComponent = gic.ID("SomeComponent")

func init() {
	gic.Add[*Component](
		gic.WithID(SomeComponent),
		gic.WithInit(func() []Component {
			return nil
			// return Component{}
		}),
	)
}

func TestComponentFound(t *testing.T) {
	_ = gic.Init()

	comp, err := gic.GetE[*Component](gic.WithID(SomeComponent))
	require.NoError(t, err)
	require.NotNil(t, comp)
}

func TestComponentNotFound(t *testing.T) {
	// comp, err := gic.GetE[*Component](NonExisting)
	// require.ErrorIs(t, err, gic.ErrNotFound)
	// require.Nil(t, comp)
}
