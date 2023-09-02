package list

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/gic"
)

type Deps struct {
	Deps []Dep
}

type Dep interface {
	Dep() string
}

var D1 = gic.ID("Dep1")

type Dep1 struct{}

func (d *Dep1) Dep() string {
	return "dep1"
}

var D2 = gic.ID("Dep2")

type Dep2 struct{}

func (d *Dep2) Dep() string {
	return "dep2"
}

func init() {
	gic.Add[*Dep1](
		gic.WithID(D1),
		gic.WithInit(func() *Dep1 { return &Dep1{} }),
	)

	gic.Add[*Dep2](
		gic.WithID(D2),
		gic.WithInit(func() *Dep2 { return &Dep2{} }),
	)

	gic.Add[*Deps](
		gic.WithInit(func() *Deps {
			return &Deps{Deps: gic.List[Dep](
				gic.Get[*Dep1](gic.WithID(D1)),
				gic.Get[*Dep2](gic.WithID(D2)),
			)}
		}),
	)
}

func TestList(t *testing.T) {
	_ = gic.Init()
	c := gic.Get[*Deps]()
	var out []string
	for _, dep := range c.Deps {
		out = append(out, dep.Dep())
	}
	require.Equal(t, []string{"dep1", "dep2"}, out)
}
