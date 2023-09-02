package add_not_in_init

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
)

func TestAddNotInInitPanics(t *testing.T) {
	defer func() {
		var err error
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = e
			}
		}

		require.ErrorIs(t, err, gic.ErrNotFromInit)
	}()

	gic.Add[string](gic.WithInit(func() string { return "" }))
}
