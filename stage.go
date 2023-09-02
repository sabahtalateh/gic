package gic

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
)

type StageRunOrder uint64

const (
	// NoOrder has same effect as InitOrder
	// Needs to differ from InitOrder and ReverseInitOrder when stage is running concurrently
	NoOrder StageRunOrder = iota
	// InitOrder as components was initialized
	InitOrder
	// ReverseInitOrder in reverse of init order. what a surprise :-)
	ReverseInitOrder
)

// stage created with RegisterStageE and
// WithID, WithInitOrder, WithReverseInitOrder, WithDisableParallel
type stage struct {
	id              id
	order           StageRunOrder // order has effect only if disableParallel is true
	disableParallel bool          // see: runInOrder, runInParallel
}

type stageOption interface{ applyStageOption(*stage) }

type withInitOrder struct{}
type withReverseInitOrder struct{}
type withDisableParallel struct{}

func (w withInitOrder) applyStageOption(s *stage)        { s.order = InitOrder }
func (w withReverseInitOrder) applyStageOption(s *stage) { s.order = ReverseInitOrder }
func (w withDisableParallel) applyStageOption(s *stage)  { s.disableParallel = true }

// WithInitOrder set stage init order to InitOrder
func WithInitOrder() withInitOrder { return withInitOrder{} }

// WithReverseInitOrder set stage init order to ReverseInitOrder
func WithReverseInitOrder() withReverseInitOrder { return withReverseInitOrder{} }

// WithDisableParallel
// If passed stage implementations will be executed in order (see: runInOrder)
// If Not passed stage implementations will be executed within goroutines (see: runInParallel)
func WithDisableParallel() withDisableParallel { return withDisableParallel{} }

// RegisterStageE add stage to container global variable (see: var globC = ...)
// No stages can be implemented on component before it registered
// errors: ErrEmptyStageName, ErrStageRegistered
func RegisterStageE(opts ...stageOption) (stage, error) {
	globC.mu.Lock()
	defer globC.mu.Unlock()
	if globC.initialized {
		return stage{}, errors.Join(ErrInitialized, fmt.Errorf("stages can not be registered after gic.Init called"))
	}

	s := stage{}
	for _, opt := range opts {
		opt.applyStageOption(&s)
	}

	if s.id.v == "" {
		return s, ErrEmptyStageName
	}

	if _, ok := globC.stages[s.id.v]; ok {
		return s, errors.Join(ErrStageRegistered, fmt.Errorf("%s", s.id))
	}

	if !s.disableParallel && s.order != NoOrder {
		globC.LogWarnf("order will not take effect without gic.WithDisableParallel")
	}

	globC.stages[s.id.v] = s

	globC.LogInfof("stage %s registered", s.id)

	return s, nil
}

// RegisterStage panics on RegisterStageE error
func RegisterStage(opts ...stageOption) stage {
	stg, err := RegisterStageE(opts...)
	check(err)
	return stg
}

// RunStage StageRunOrder has no effect if disableParallel is false
// errors: ErrStageNotRegistered
func RunStage(ctx context.Context, s stage) error {
	globC.mu.Lock()
	if !globC.initialized {
		globC.mu.Unlock()
		return errors.Join(ErrNotInitialized, fmt.Errorf("before run stage container must be initialized with gic.Init"))
	}
	globC.mu.Unlock()

	stg, ok := globC.stages[s.id.v]
	if !ok {
		return errors.Join(ErrStageNotRegistered, fmt.Errorf("%s", s.id))
	}

	if stg.disableParallel {
		return runInOrder(ctx, globC, s)
	}

	return runInParallel(ctx, globC, s)
}

func runInOrder(ctx context.Context, c *container, s stage) error {
	fns := c.stageFns[s.id.v]
	if s.order != ReverseInitOrder {
		for i := 0; i < len(fns); i++ {
			if err := fns[i](ctx); err != nil {
				return err
			}
		}
	} else {
		for i := len(fns) - 1; i >= 0; i-- {
			if err := fns[i](ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func runInParallel(ctx context.Context, c *container, s stage) error {
	fns := c.stageFns[s.id.v]
	eg, egCtx := errgroup.WithContext(ctx)

	for i := 0; i < len(fns); i++ {
		i := i
		eg.Go(func() error {
			if err := fns[i](egCtx); err != nil {
				return err
			}
			return nil
		})
	}

	return eg.Wait()
}
