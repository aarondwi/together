package combiner

import (
	"context"
	"errors"

	e "github.com/aarondwi/together/engine"
)

/*
 * This sub-packages implement promise-like functionality for multiple BatchResults.
 * The goal is to allow users to reduce waiting time, by doing all of them in parallel.
 *
 * Roughly adapted from https://github.com/chebyrash/promise/blob/master/promise.go
 *
 * Notes do to the nature of the problem,
 * this implementation has few allocations on the path.
 * Use this implementation sparingly.
 */

// ErrCombinerNumWorkerZeroOrLess is returned when
// numOfWorker given to NewCombiner is <= zero
var ErrCombinerNumWorkerZeroOrLess = errors.New("Number of worker should be > 0")
var emptySlicesResponse = make([]interface{}, 0)

type resolutionHelper struct {
	index  int
	result interface{}
	err    error
}

type Combiner struct{}

func NewCombiner(numOfWorker int) (*Combiner, error) {
	if numOfWorker <= 0 {
		return nil, ErrCombinerNumWorkerZeroOrLess
	}

	return nil, nil
}

// All waits for all BatchResult to be returned, or one error be returned.
//
// If all succeed, this will return array of response in the same order as params, with nil error
// Else, it returns the first error with empty slices.
func All(brs []e.BatchResult) ([]interface{}, error) {
	if len(brs) == 0 {
		return emptySlicesResponse, nil
	}

	results := make([]interface{}, len(brs))

	// we buffer the channels, to prevent goroutine leak
	resultCh := make(chan resolutionHelper, len(brs))
	errCh := make(chan error, len(brs))

	for i, br := range brs {
		go func(j int, br e.BatchResult) {
			res, err := br.GetResult()
			if err != nil {
				errCh <- err
			}
			resultCh <- resolutionHelper{index: j, result: res}
		}(i, br)
	}

	for i := 0; i < len(brs); i++ {
		select {
		case resHelper := <-resultCh:
			results[resHelper.index] = resHelper.result
		case err := <-errCh:
			return emptySlicesResponse, err
		}
	}
	return results, nil
}

// AllWithContext is the same as `All` call, but with context idiom.
//
// Note unless you need the context idiom, it is preferable
// to use `All` call instead, as it has less allocations (so it is faster)
func AllWithContext(
	ctx context.Context, brs []e.BatchResult) ([]interface{}, error) {
	return nil, nil
}

// Race waits until either one result or one error is returned.
func Race(brs []e.BatchResult) (interface{}, error) {
	return nil, nil
}

// RaceWithContext is the same as `Race` call, but with context idiom.
//
// Note unless you need the context idiom, it is preferable
// to use `Race` call instead, as it has less allocations (so it is faster)
func RaceWithContext(
	ctx context.Context, brs []e.BatchResult) (interface{}, error) {
	return nil, nil
}

// AllSettled waits until all results/errors to be returned.
//
// Result will be nil if err exists, else err will be nil
// in the same order as params (so may be nil on one index, and not nil on another)
func AllSettled(brs []e.BatchResult) ([]interface{}, []error) {
	return nil, nil
}

// AllSettledWithContext is the same as `AllSettled` call, but with context idiom.
//
// Note unless you need the context idiom, it is preferable
// to use `AllSettled` call instead, as it has less allocations (so it is faster)
func AllSettledWithContext(
	ctx context.Context, brs []e.BatchResult) ([]interface{}, []error) {
	return nil, nil
}
