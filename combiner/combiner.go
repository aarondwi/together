package combiner

import (
	"context"

	com "github.com/aarondwi/together/common"
	e "github.com/aarondwi/together/engine"
)

/*
 * This sub-packages implement promise-like functionality for multiple BatchResults.
 * The goal is to allow users to reduce latency on a per request basis,
 * by doing all of them in parallel.
 *
 * Roughly adapted from https://github.com/chebyrash/promise/blob/master/promise.go,
 * with few semantic changes and optimization in place (especially with worker pool)
 *
 * Notes do to the nature of the problem,
 * this implementation has quite a few allocations on hot path.
 * Use this implementation sparingly.
 */

type resolutionHelper struct {
	index  int
	result interface{}
	err    error
}

type Combiner struct {
	wp *com.WorkerPool
}

func NewCombiner(wp *com.WorkerPool) (*Combiner, error) {
	if wp == nil {
		return nil, com.ErrNilWorkerPool
	}
	return &Combiner{wp: wp}, nil
}

// All waits for all BatchResult to be returned, or one error be returned.
//
// If all succeed, this will return array of response in the same order as params, with nil error
// Else, it returns the first error only.
func (c *Combiner) All(brs []e.BatchResult) ([]interface{}, error) {
	if len(brs) == 0 {
		return nil, nil
	}

	results := make([]interface{}, len(brs))

	// we buffer the channels, to prevent goroutine leak
	resultCh := make(chan resolutionHelper, len(brs))
	errCh := make(chan error, len(brs))

	for i, br := range brs {
		c.wp.Submit(func(j int, br e.BatchResult) func() {
			return func() {
				res, err := br.GetResult()
				if err != nil {
					errCh <- err
					return
				}
				resultCh <- resolutionHelper{index: j, result: res}
			}
		}(i, br))
	}

	for i := 0; i < len(brs); i++ {
		select {
		case resHelper := <-resultCh:
			results[resHelper.index] = resHelper.result
		case err := <-errCh:
			return nil, err
		}
	}
	return results, nil
}

// AllWithContext is the same as `All` call, but with context idiom.
//
// Note unless you need the context idiom, it is preferable
// to use `All` call instead, as it has less allocations (so it is faster)
func (c *Combiner) AllWithContext(
	ctx context.Context, brs []e.BatchResult) ([]interface{}, error) {
	// fast path, ctx already done
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(brs) == 0 {
		return nil, nil
	}

	results := make([]interface{}, len(brs))

	// we buffer the channels, to prevent goroutine leak
	resultCh := make(chan resolutionHelper, len(brs))
	errCh := make(chan error, len(brs))

	for i, br := range brs {
		c.wp.Submit(func(ctx context.Context, j int, br e.BatchResult) func() {
			return func() {
				res, err := br.GetResultWithContext(ctx)
				if err != nil {
					errCh <- err
					return
				}
				resultCh <- resolutionHelper{index: j, result: res}
			}
		}(ctx, i, br))
	}

	for i := 0; i < len(brs); i++ {
		select {
		case resHelper := <-resultCh:
			results[resHelper.index] = resHelper.result
		case err := <-errCh:
			return nil, err
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
	// we buffer the channels, to prevent goroutine leak
	resultCh := make(chan interface{}, len(brs))
	errCh := make(chan resolutionHelper, len(brs))

	for i, br := range brs {
		c.wp.Submit(func(j int, br e.BatchResult) func() {
			return func() {
				res, err := br.GetResult()
				if err != nil {
					errCh <- resolutionHelper{index: j, err: err}
					return
				}
				resultCh <- res
			}
		}(i, br))
	}

	for i := 0; i < len(brs); i++ {
		select {
		case result := <-resultCh:
			return result, nil
		case resHelper := <-errCh:
			errs[resHelper.index] = resHelper.err
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

	if len(brs) == 0 {
		return nil, nil
	}
	// fast path, ctx already done
	select {
	case <-ctx.Done():
		errs := make([]error, 0, len(brs))
		for i := 0; i < len(brs); i++ {
			errs = append(errs, ctx.Err())
		}
		return nil, errs
	default:
	}

	errs := make([]error, len(brs))

	// we buffer the channels, to prevent goroutine leak
	resultCh := make(chan interface{}, len(brs))
	errCh := make(chan resolutionHelper, len(brs))

	for i, br := range brs {
		c.wp.Submit(func(ctx context.Context, j int, br e.BatchResult) func() {
			return func() {
				res, err := br.GetResultWithContext(ctx)
				if err != nil {
					errCh <- resolutionHelper{index: j, err: err}
					return
				}
				resultCh <- res
			}
		}(ctx, i, br))
	}

	for i := 0; i < len(brs); i++ {
		select {
		case result := <-resultCh:
			return result, nil
		case resHelper := <-errCh:
			errs[resHelper.index] = resHelper.err
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

	// we buffer the channels, to prevent goroutine leak
	resultCh := make(chan resolutionHelper, len(brs))
	errCh := make(chan resolutionHelper, len(brs))

	for i, br := range brs {
		c.wp.Submit(func(j int, br e.BatchResult) func() {
			return func() {
				res, err := br.GetResult()
				if err != nil {
					errCh <- resolutionHelper{index: j, err: err}
					return
				}
				resultCh <- resolutionHelper{index: j, result: res}
			}
		}(i, br))
	}

	for i := 0; i < len(brs); i++ {
		select {
		case resHelper := <-resultCh:
			results[resHelper.index] = resHelper.result
			errs[resHelper.index] = nil
		case errHelper := <-errCh:
			results[errHelper.index] = nil
			errs[errHelper.index] = errHelper.err
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
	if len(brs) == 0 {
		return nil, nil
	}
	// fast path, ctx already done
	select {
	case <-ctx.Done():
		errs := make([]error, 0, len(brs))
		for i := 0; i < len(brs); i++ {
			errs = append(errs, ctx.Err())
		}
		return nil, errs
	default:
	}

	results := make([]interface{}, len(brs))
	errs := make([]error, len(brs))

	// we buffer the channels, to prevent goroutine leak
	resultCh := make(chan resolutionHelper, len(brs))
	errCh := make(chan resolutionHelper, len(brs))

	for i, br := range brs {
		c.wp.Submit(func(ctx context.Context, j int, br e.BatchResult) func() {
			return func() {
				res, err := br.GetResultWithContext(ctx)
				if err != nil {
					errCh <- resolutionHelper{index: j, err: err}
					return
				}
				resultCh <- resolutionHelper{index: j, result: res}
			}
		}(ctx, i, br))
	}

	for i := 0; i < len(brs); i++ {
		select {
		case resHelper := <-resultCh:
			results[resHelper.index] = resHelper.result
			errs[resHelper.index] = nil
		case errHelper := <-errCh:
			results[errHelper.index] = nil
			errs[errHelper.index] = errHelper.err
		}
	}
	return results, errs
}
