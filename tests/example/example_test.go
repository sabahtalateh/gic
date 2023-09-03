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

func TestGreeter(t *testing.T) {
	var err error

	err = gic.Init()
	require.Nil(t, err)

	var g *internal.Greeter
	g, err = gic.GetE[*internal.Greeter]()
	require.Nil(t, err)
	require.Equal(t, "Hello World!", g.Greet("World"))

	g, err = gic.GetE[*internal.Greeter](gic.WithID(internal.RussianGreeter))
	require.Nil(t, err)
	require.Equal(t, "Привет Мир!", g.Greet("Мир"))

	g, err = gic.GetE[*internal.Greeter](gic.WithID(internal.ChineseGreeter))
	require.Nil(t, err)
	require.Equal(t, "你好 世界!", g.Greet("世界"))
}

func TestStartStop(t *testing.T) {
	var err error

	err = gic.Init()
	require.Nil(t, err)

	err = gic.Start(context.Background())
	require.Nil(t, err)

	var ne *internal.NumbersEater
	ne, err = gic.GetE[*internal.NumbersEater]()
	require.Nil(t, err)

	ne.Feed(1)
	ne.Feed(2)

	require.Eventually(t, func() bool {
		return slices.Contains(ne.Eaten(), 1) &&
			slices.Contains(ne.Eaten(), 2)
	}, 1*time.Second, 10*time.Millisecond)

	err = gic.Stop(context.Background())
	require.Nil(t, err)
}

func TestMyStage(t *testing.T) {
	var err error

	err = gic.Init()
	require.Nil(t, err)

	err = gic.RunStage(context.Background(), internal.MyStage)
	require.Nil(t, err)

	require.Equal(t, 999, gic.Get[*internal.Dummy]().X)
}
