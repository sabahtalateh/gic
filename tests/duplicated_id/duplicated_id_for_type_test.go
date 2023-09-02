package duplicated_id

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
)

type Component struct{ prop string }

var SomeID = gic.ID("SomeID")

func init() {
	gic.Add[*Component](
		gic.WithID(SomeID),
		gic.WithInit(func() *Component { return &Component{prop: "v1"} }),
	)

	gic.Add[*Component](
		gic.WithID(SomeID),
		gic.WithInit(func() *Component { return &Component{prop: "v2"} }),
	)
}

func TestErrorOnDuplicatedID(t *testing.T) {
	err := gic.Init()
	require.ErrorIs(t, err, gic.ErrComponentAdded)
}
