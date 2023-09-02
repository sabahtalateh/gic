package interfaces

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
)

var initErr error
var wg sync.WaitGroup

func init() {
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				switch err := r.(type) {
				case error:
					initErr = err
				}
			}
			wg.Done()
		}()
		gic.Add[string](
			gic.WithInit(func() string { return "" }),
		)
	}()
}

func TestAddInInitGoroutinePanics(t *testing.T) {
	wg.Wait()
	require.ErrorIs(t, initErr, gic.ErrNotFromInit)
}
