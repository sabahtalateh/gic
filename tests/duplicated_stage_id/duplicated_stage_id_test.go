package duplicated_stage_id

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
	gic.RegisterStage(gic.WithID(gic.ID("ID1")))
	gic.RegisterStage(gic.WithID(gic.ID("ID1")))
}

func TestStageIDUnique(t *testing.T) {
	require.ErrorIs(t, err, gic.ErrStageRegistered)
}
