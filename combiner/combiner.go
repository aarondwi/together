package combiner

import (
	"context"

	e "github.com/aarondwi/together/engine"
	WP "github.com/aarondwi/together/workerpool"
)

type resolutionHelper struct {
	index  int
	result interface{}
	err    error
}

// Combiner implements promise-like functionality for multiple BatchResults.
// The goal is to allow users to reduce latency on a per request basis,
// by doing all of them in parallel.
//
// Roughly adapted from https://github.com/chebyrash/promise/blob/master/promise.go,
// with few semantic changes and optimization in place (especially with worker pool)
//
// If the given WorkerPool is nil, it will create a goroutine for each BatchResults given.
// So it is recommended to supply a non-nil, big enough worker-pool.
//
// On every call, buffered channel is created, to prevent goroutine leak.
//
// Notes do to the nature of the problem,
// this implementation has quite a few allocations on hot path.
// Use this implementation sparingly.
type Combiner struct {
	wp *WP.WorkerPool
}

// NewCombiner creates our combiner, given the WorkerPool.
func NewCombiner(wp *WP.WorkerPool) *Combiner {
	return &Combiner{wp: wp}
}

func (c *Combiner) applyWorker(
	j int, br e.BatchResult,
	ch chan<- resolutionHelper) {

	res, err := br.GetResult()
	if err != nil {
		ch <- resolutionHelper{index: j, err: err}
		return
	}
	ch <- resolutionHelper{index: j, result: res}
}

func (c *Combiner) apply(
	brs []e.BatchResult, ch chan<- resolutionHelper) {

	if c.wp == nil {
		for i, br := range brs {
			func(j int, br e.BatchResult) {
				go c.applyWorker(j, br, ch)
			}(i, br)
		}
	} else {
		for i, br := range brs {
			func(j int, br e.BatchResult) {
				c.wp.Submit(func() {
					c.applyWorker(j, br, ch)
				})
			}(i, br)
		}
	}
}

func (c *Combiner) applyWithCtxWorker(
	ctx context.Context,
	j int, br e.BatchResult,
	ch chan<- resolutionHelper) {

	res, err := br.GetResultWithContext(ctx)
	if err != nil {
		ch <- resolutionHelper{index: j, err: err}
		return
	}
	ch <- resolutionHelper{index: j, result: res}
}

func (c *Combiner) applyWithCtx(
	ctx context.Context,
	brs []e.BatchResult, ch chan<- resolutionHelper) {

	if c.wp == nil {
		for i, br := range brs {
			func(j int, br e.BatchResult) {
				go c.applyWithCtxWorker(ctx, j, br, ch)
			}(i, br)
		}
	} else {
		for i, br := range brs {
			func(j int, br e.BatchResult) {
				c.wp.Submit(func() {
					c.applyWithCtxWorker(ctx, j, br, ch)
				})
			}(i, br)
		}
	}
}

// AllSuccesses waits for all BatchResult to be returned, or one error be returned.
//
// If all succeed, this will return array of response in the same order as params, with nil error.
// Else, it returns the first error only.
func (c *Combiner) AllSuccesses(brs []e.BatchResult) ([]interface{}, error) {
	if len(brs) == 0 {
		return nil, nil
	}

	results := make([]interface{}, len(brs))
	ch := make(chan resolutionHelper, len(brs))
	c.apply(brs, ch)

	for i := 0; i < len(brs); i++ {
		resHelper := <-ch
		if resHelper.result != nil {
			results[resHelper.index] = resHelper.result
		} else {
			return nil, resHelper.err
		}
	}
	return results, nil
}

// AllSuccessesWithContext is the same as `All` call, but with context idiom.
//
// Note unless you need the context idiom, it is preferable
// to use `All` call instead, as it has less allocations (so it is faster)
func (c *Combiner) AllSuccessesWithContext(
	ctx context.Context, brs []e.BatchResult) ([]interface{}, error) {

	// fast path
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if len(brs) == 0 {
		return nil, nil
	}

	// we manage our own cancellation
	// we defer on return, so either all completes, and just be cautious by cancelling
	// or error returns and we cancel all pending calls
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]interface{}, len(brs))
	ch := make(chan resolutionHelper, len(brs))
	c.applyWithCtx(ctx2, brs, ch)

	for i := 0; i < len(brs); i++ {
		resHelper := <-ch
		if resHelper.result != nil {
			results[resHelper.index] = resHelper.result
		} else {
			return nil, resHelper.err
		}
	}
	return results, nil
}

// Race waits until either one result is returned.
//
// If all failed, it returns all the errors, in the same order as params
func (c *Combiner) Race(brs []e.BatchResult) (interface{}, []error) {
	if len(brs) == 0 {
		return nil, nil
	}

	errs := make([]error, len(brs))
	ch := make(chan resolutionHelper, len(brs))
	c.apply(brs, ch)

	for i := 0; i < len(brs); i++ {
		resHelper := <-ch
		if resHelper.err != nil {
			errs[resHelper.index] = resHelper.err
		} else {
			return resHelper.result, nil
		}
	}

	return nil, errs
}

// RaceWithContext is the same as `Race` call, but with context idiom.
//
// Note unless you need the context idiom, it is preferable
// to use `Race` call instead, as it has less allocations (so it is faster)
func (c *Combiner) RaceWithContext(
	ctx context.Context, brs []e.BatchResult) (interface{}, []error) {

	// fast path
	select {
	case <-ctx.Done():
		errs := make([]error, 0, len(brs))
		for i := 0; i < len(brs); i++ {
			errs = append(errs, ctx.Err())
		}
		return nil, errs
	default:
	}
	if len(brs) == 0 {
		return nil, nil
	}

	// we manage our own cancellation
	// we defer on return, so either one completes and we cancel all pending calls
	// or all fail, and just be cautious by cancelling
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make([]error, len(brs))
	ch := make(chan resolutionHelper, len(brs))
	c.applyWithCtx(ctx2, brs, ch)

	for i := 0; i < len(brs); i++ {
		resHelper := <-ch
		if resHelper.err != nil {
			errs[resHelper.index] = resHelper.err
		} else {
			return resHelper.result, nil
		}
	}

	return nil, errs
}

// Every waits until all results/errors to be returned.
//
// Result will be nil if err exists, else err will be nil
// in the same order as params (so may be nil on one index, and not nil on another)
func (c *Combiner) Every(brs []e.BatchResult) ([]interface{}, []error) {
	if len(brs) == 0 {
		return nil, nil
	}

	results := make([]interface{}, len(brs))
	errs := make([]error, len(brs))
	ch := make(chan resolutionHelper, len(brs))
	c.apply(brs, ch)

	for i := 0; i < len(brs); i++ {
		resHelper := <-ch
		if resHelper.result != nil {
			results[resHelper.index] = resHelper.result
			errs[resHelper.index] = nil
		} else {
			results[resHelper.index] = nil
			errs[resHelper.index] = resHelper.err
		}
	}
	return results, errs
}

// EveryWithContext is the same as `Every` call, but with context idiom.
//
// Note unless you need the context idiom, it is preferable
// to use `Every` call instead, as it has less allocations (so it is faster)
func (c *Combiner) EveryWithContext(
	ctx context.Context, brs []e.BatchResult) ([]interface{}, []error) {

	// fast path
	select {
	case <-ctx.Done():
		errs := make([]error, 0, len(brs))
		for i := 0; i < len(brs); i++ {
			errs = append(errs, ctx.Err())
		}
		return nil, errs
	default:
	}
	if len(brs) == 0 {
		return nil, nil
	}

	results := make([]interface{}, len(brs))
	errs := make([]error, len(brs))
	ch := make(chan resolutionHelper, len(brs))
	c.applyWithCtx(ctx, brs, ch)

	for i := 0; i < len(brs); i++ {
		resHelper := <-ch
		if resHelper.result != nil {
			results[resHelper.index] = resHelper.result
			errs[resHelper.index] = nil
		} else {
			results[resHelper.index] = nil
			errs[resHelper.index] = resHelper.err
		}
	}
	return results, errs
}
