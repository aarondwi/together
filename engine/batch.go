package engine

import (
	"context"
	"sync"
)

// Batch is our wrapper, promise implementation.
//
// This implementation is *NOT* goroutine-safe,
// but should be okay as it is always created without pointer (on stack)
type Batch struct {
	ID      uint64
	args    map[uint64]interface{}
	argSize int
	wg      sync.WaitGroup
	results map[uint64]interface{}
	err     error
}

type BatchResult struct {
	id    uint64
	batch *Batch
}

var EmptyBatchResult = BatchResult{}

// NewBatch creates a new batch
//
// Once taken to work on, nothing should be put anymore
func NewBatch(id uint64) *Batch {
	b := &Batch{
		ID: id,
		// how to pool map?
		args: make(map[uint64]interface{}),
	}
	b.wg.Add(1)
	return b
}

// Put into the Batch.args
func (b *Batch) Put(id uint64, arg interface{}) BatchResult {
	b.args[id] = arg
	b.argSize++
	return BatchResult{
		batch: b,
		id:    id,
	}
}

// GetResult waits until the batch is done, then match the result for each caller.
//
// This is not automatically called inside Engine's `Submit()` call.
// This is by design, to allow user to specify when to wait,
// and also makes it easier to be reused inside another utility functions.
func (br *BatchResult) GetResult() (interface{}, error) {
	br.batch.wg.Wait()
	if br.batch.err != nil {
		return nil, br.batch.err
	}
	res, ok := br.batch.results[br.id]
	if !ok {
		return nil, ErrResultNotFound
	}
	return res, nil
}

// GetResultWithContext waits until either the batch or ctx is done,
// then match the result for each caller.
//
// Note that unless you need to use the context idiom, it is recommended
// to use `GetResult()` call instead, as it has much, much less allocation (only interface{} typecasting).
// This API need to create another goroutine, and 2 channels to manage its functionality.
// (And of course, using context's `WithCancel` or `WithTimeout` also creates goroutines)
//
// Beware that because how go handles local variable,
// if you are using multiple `GetResultWithContext` in single function,
// be sure to assign it to different local variable.
// See https://stackoverflow.com/questions/25919213/why-does-go-handle-closures-differently-in-goroutines
// for details
func (br *BatchResult) GetResultWithContext(
	ctx context.Context) (interface{}, error) {
	// fast path, ctx already done
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// size 1, prevent goroutine leak
	// if either the context or the other goroutine done first
	resultCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	go func() {
		res, err := br.GetResult()
		br.batch = nil
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- res
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultCh:
		return res, nil
	case err := <-errCh:
		return nil, err
	}
}
