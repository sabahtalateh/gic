package not_found_with_id

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
)

type Component struct {
	prop string
}

var SomeComponent = gic.ID("SomeComponent")
var NotExists = gic.ID("NotExists")

func init() {
	gic.Add[*Component](
		gic.WithID(SomeComponent),
		gic.WithInit(func() *Component {
			return &Component{prop: "value"}
		}),
	)
}

func TestComponentNotFoundWithID(t *testing.T) {
	_ = gic.Init()

	_, err := gic.Get[*Component](gic.WithID(NotExists))
	require.ErrorIs(t, err, gic.ErrNotFound)
}
