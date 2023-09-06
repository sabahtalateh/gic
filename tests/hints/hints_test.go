package hints

import (
	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
	"testing"
)

type A struct{}

func init() {
	gic.Add[*A](
		gic.WithInit(func() *A { return &A{} }),
	)

	gic.Add[string](
		gic.WithInit(func() string { return "" }),
	)
}

func TestHint(t *testing.T) {
	err := gic.Init()
	require.NoError(t, err)

	_, err = gic.Get[A]()
	require.ErrorIs(t, err, gic.ErrNotFound)
	require.Contains(t, err.Error(), "type hints.A not found but *hints.A found")
}
