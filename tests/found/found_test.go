package found

import (
	"testing"
)

type Component struct{}

// var Existing = gic.ID()
// var NonExisting = gic.ID()

func init() {
	// _ = gic.Init()
	// gic.Add[*Component](Existing, func() (*Component, error) { return &Component{}, nil })
}

func TestComponentFound(t *testing.T) {
	// comp, err := gic.GetE[*Component](Existing)
	// require.NoError(t, err)
	// require.NotNil(t, comp)
}

func TestComponentNotFound(t *testing.T) {
	// comp, err := gic.GetE[*Component](NonExisting)
	// require.ErrorIs(t, err, gic.ErrNotFound)
	// require.Nil(t, comp)
}
