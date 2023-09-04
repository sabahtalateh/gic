package found_without_id

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
)

type Component struct {
	prop string
}

func init() {
	gic.Add[*Component](
		gic.WithInit(func() *Component {
			return &Component{prop: "value"}
		}),
	)
}

func TestComponentFoundWithoutID(t *testing.T) {
	_ = gic.Init()

	comp, err := gic.Get[*Component]()
	require.NoError(t, err)
	require.NotNil(t, comp)
	require.Equal(t, comp.prop, "value")
}
