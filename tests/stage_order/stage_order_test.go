package stage_order

import (
	"context"
	"testing"

	"github.com/sabahtalateh/gic"
	"github.com/stretchr/testify/require"
)

var direct = gic.RegisterStage(
	gic.WithID(gic.ID("direct")),
	gic.WithDisableParallel(),
	gic.WithInitOrder(),
)

var reverse = gic.RegisterStage(
	gic.WithID(gic.ID("reverse")),
	gic.WithDisableParallel(),
	gic.WithReverseInitOrder(),
)

var directOut []int
var reverseOut []int

func init() {
	gic.Add[string](
		gic.WithID(gic.ID("1")),
		gic.WithInit(func() string { return "1" }),
		gic.WithStageImpl(direct, func(ctx context.Context, s string) error {
			directOut = append(directOut, 1)
			return nil
		}),
		gic.WithStageImpl(reverse, func(ctx context.Context, s string) error {
			reverseOut = append(reverseOut, 1)
			return nil
		}),
	)

	gic.Add[string](
		gic.WithID(gic.ID("2")),
		gic.WithInit(func() string { return "2" }),
		gic.WithStageImpl(direct, func(ctx context.Context, s string) error {
			directOut = append(directOut, 2)
			return nil
		}),
		gic.WithStageImpl(reverse, func(ctx context.Context, s string) error {
			reverseOut = append(reverseOut, 2)
			return nil
		}),
	)
}

func TestStageOrder(t *testing.T) {
	var err error

	_ = gic.Init()
	err = gic.RunStage(context.Background(), direct)
	require.NoError(t, err)

	err = gic.RunStage(context.Background(), reverse)
	require.NoError(t, err)

	require.Equal(t, []int{1, 2}, directOut)
	require.Equal(t, []int{2, 1}, reverseOut)
}
