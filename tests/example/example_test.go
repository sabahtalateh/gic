package example

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
	"github.com/sabahtalateh/gic/tests/example/internal"
)

func TestExample(t *testing.T) {
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
