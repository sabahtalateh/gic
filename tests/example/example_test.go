package example

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/example/internal"
)

func init() {
	if err := gic.Init(); err != nil {
		panic(err)
	}
}

func TestGreeter(t *testing.T) {
	var err error

	var g *internal.Greeter
	g, err = gic.Get[*internal.Greeter]()
	require.NoError(t, err)
	require.Equal(t, "Hello World!", g.Greet("World"))

	g, err = gic.Get[*internal.Greeter](gic.WithID(internal.RussianGreeter))
	require.NoError(t, err)
	require.Equal(t, "Привет Мир!", g.Greet("Мир"))

	g, err = gic.Get[*internal.Greeter](gic.WithID(internal.ChineseGreeter))
	require.NoError(t, err)
	require.Equal(t, "你好 世界!", g.Greet("世界"))
}

func TestStartStop(t *testing.T) {
	var err error

	err = gic.Start(context.Background())
	require.NoError(t, err)

	var ne *internal.NumbersEater
	ne, err = gic.Get[*internal.NumbersEater]()
	require.NoError(t, err)

	ne.Feed(1)
	ne.Feed(2)

	require.Eventually(t, func() bool {
		return slices.Contains(ne.Eaten(), 1) &&
			slices.Contains(ne.Eaten(), 2)
	}, 1*time.Second, 10*time.Millisecond)

	err = gic.Stop(context.Background())
	require.NoError(t, err)
}

func TestMyStage(t *testing.T) {
	var err error

	err = gic.RunStage(context.Background(), internal.MyStage)
	require.NoError(t, err)

	require.Equal(t, 999, gic.MustGet[*internal.Dummy]().X)
}
