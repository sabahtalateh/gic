package kind_of_real_project

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic/tests/kind_of_real_project/run"
)

func TestKindOfRealProject(t *testing.T) {
	out := run.Run()
	require.Equal(t, []string{
		"sending message to Ivan",
		"sending message to Petr",
		"sending message to Anonymous",
	}, out)
}
