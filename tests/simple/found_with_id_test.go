package simple

import (
	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
	"testing"
)

type Component struct {
	prop string
}

var SomeComponent = gic.ID("SomeComponent")

func init() {
	gic.Add[*Component](
		gic.WithID(SomeComponent),
		gic.WithInit(func() *Component {
			return &Component{prop: "value"}
		}),
	)
}

func TestComponentFoundByID(t *testing.T) {
	_ = gic.Init()

	comp, err := gic.GetE[*Component](gic.WithID(SomeComponent))
	require.NoError(t, err)
	require.NotNil(t, comp)
	require.Equal(t, comp.prop, "value")
}
