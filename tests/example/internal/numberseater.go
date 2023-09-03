package internal

import (
	"context"
	"sync"

	"github.com/sabahtalateh/gic"
)

type NumbersEater struct {
	mu    sync.Mutex
	c     chan int
	eaten []int
}

func (n *NumbersEater) Start() {
	go func() {
		for number := range n.c {
			n.mu.Lock()
			n.eaten = append(n.eaten, number)
			n.mu.Unlock()
		}
	}()
}

func (n *NumbersEater) Stop() {
	close(n.c)
}

func (n *NumbersEater) Feed(num int) {
	n.c <- num
	return
}

func (n *NumbersEater) Eaten() []int {
	return n.eaten
}

func init() {
	gic.Add[*NumbersEater](
		gic.WithInit(func() *NumbersEater {
			return &NumbersEater{c: make(chan int)}
		}),
		gic.WithStart(func(_ context.Context, ne *NumbersEater) error {
			ne.Start()
			return nil
		}),
		gic.WithStop(func(ctx context.Context, ne *NumbersEater) error {
			ne.Stop()
			return nil
		}),
	)
}
