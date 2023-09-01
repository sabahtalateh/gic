package simple

import (
	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
	"testing"
)

type Component2 struct {
	prop string
}

func init() {
	gic.Add[*Component2](
		gic.WithInit(func() *Component2 {
			return &Component2{prop: "value"}
		}),
	)
}

func TestComponentFoundWithoutID(t *testing.T) {
	_ = gic.Init()

	comp, err := gic.GetE[*Component2]()
	require.NoError(t, err)
	require.NotNil(t, comp)
	require.Equal(t, comp.prop, "value")
}
