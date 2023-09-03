package init_required

import (
	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
	"testing"
)

var err error

func init() {
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = e
			}
		}
	}()

	gic.Add[string](
		gic.WithID(gic.ID("1")),
	)
}

func TestInitRequire(t *testing.T) {
	require.ErrorIs(t, err, gic.ErrNoInit)
}
