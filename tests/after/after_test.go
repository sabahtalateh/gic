package after

import (
	"testing"
)

type Component struct {
	msg string
}

func (c *Component) Start() {
	c.msg = "started"
}

func init() {
	// _ = gic.Init()
	// gic.Add[*Component](gic.noID, func() (*Component, error) {
	// 	return &Component{}, nil
	// }, func(_ *gic.Container, c *Component) error {
	// 	c.Start()
	// 	return nil
	// })
}

func TestAfter(t *testing.T) {
	// c := gic.Get[*Component](gic.noID)
	// require.Equal(t, "started", c.msg)
}
