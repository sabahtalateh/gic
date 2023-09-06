package multiple_unordered_inits

import (
	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
	"testing"
)

type A struct {
	a string
}

type B struct {
	A A
}

func init() {
	gic.Add[B](
		gic.WithInit[B](func() B { return B{A: gic.MustGet[A]()} }),
	)
}

func init() {
	gic.Add[A](
		gic.WithInit(func() A { return A{a: "hello"} }),
	)
}

func TestMultipleUnorderedInits(t *testing.T) {
	err := gic.Init()
	require.NoError(t, err)

	b, err := gic.Get[B]()
	require.NoError(t, err)
	require.Equal(t, "hello", b.A.a)
}
