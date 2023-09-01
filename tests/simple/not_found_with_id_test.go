package simple

import (
	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
	"testing"
)

type Component3 struct {
	prop string
}

var SomeComponent3 = gic.ID("SomeComponent3")
var NotExists = gic.ID("NotExists")

func init() {
	gic.Add[*Component3](
		gic.WithID(SomeComponent3),
		gic.WithInit(func() *Component3 {
			return &Component3{prop: "value"}
		}),
	)
}

func TestComponentNotFoundWithID(t *testing.T) {
	_ = gic.Init()

	_, err := gic.GetE[*Component3](gic.WithID(NotExists))
	require.ErrorIs(t, err, gic.ErrNotFound)
}
