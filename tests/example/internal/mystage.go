package internal

import (
	"context"
	"github.com/sabahtalateh/gic"
)

var MyStage = gic.RegisterStage(
	gic.WithID(gic.ID("MyStage")),
	gic.WithDisableParallel(),
	gic.WithInitOrder(),
)

type Dummy struct {
	X int
}

func init() {
	gic.Add[*Dummy](
		gic.WithInit(func() *Dummy {
			return &Dummy{}
		}),
		gic.WithStageImpl(MyStage, func(ctx context.Context, d *Dummy) error {
			d.X = 999
			return nil
		}),
	)
}
