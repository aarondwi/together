package together

import (
	"sync"
)

// Batch is our wrapper, promise implementation.
// This implementation is *NOT* goroutine-safe
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

// GetResult waits until the batch is done, then match the result for each caller
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
