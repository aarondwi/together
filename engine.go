package together

import (
	"context"
	"sync"
	"time"
)

// Engine is our batch-controller, loosely adapted from
// https://github.com/grab/async/blob/master/batch.go,
// designed for business-logic use-case.
//
// User just need to specify the config on `NewEngine` call,
// and then use `Submit()` call on logic code.
//
// This implementation is goroutine-safe
type Engine struct {
	mu           *sync.Mutex
	newBatch     *sync.Cond
	batchID      uint64
	taskID       uint64
	fn           func(map[uint64]interface{}) (map[uint64]interface{}, error)
	batchChan    chan *Batch
	currentBatch *Batch
	argSizeLimit int
	waitDuration time.Duration
}

// NewEngine creates the engine with specified config
//
// The given `fn` should be goroutine-safe
func NewEngine(
	numOfWorker int,
	argSizeLimit int,
	waitDuration time.Duration,
	fn func(map[uint64]interface{}) (map[uint64]interface{}, error)) (*Engine, error) {
	if numOfWorker <= 0 {
		return nil, ErrNumOfWorkerLessThanEqualZero
	}
	if argSizeLimit <= 1 {
		return nil, ErrArgSizeLimitLessThanEqualOne
	}
	mu := &sync.Mutex{}
	newBatch := sync.NewCond(mu)
	e := &Engine{
		mu:           mu,
		newBatch:     newBatch,
		fn:           fn,
		batchChan:    make(chan *Batch, numOfWorker),
		argSizeLimit: argSizeLimit,
		waitDuration: waitDuration,
	}
	for i := 0; i < numOfWorker; i++ {
		go e.worker()
	}
	go e.batchLatencyChecker()
	return e, nil
}

func (e *Engine) worker() {
	for {
		b := <-e.batchChan
		m, err := e.fn(b.args)
		b.results = m
		b.err = err
		b.wg.Done()
	}
}

// This function should only be called when holding the mutex
//
// Note that, it may still hold the mutex when the chan is full,
// but IMO that is a better option, cause either:
//
// 1. backend can't keep up
//
// 2. just bad engine configuration (numOfWorker too low, etc)
//
// to not balloon the memory requirement
func (e *Engine) readyToWork() {
	e.batchChan <- e.currentBatch
	e.currentBatch = nil
}

// this implementation may skip some id to check
// but that means the `Submit()` call is much faster
// than the wait.
//
// In that case, no need to wait for those.
func (e *Engine) batchLatencyChecker() {
	var IDToTrack uint64
	for {
		e.mu.Lock()
		for e.currentBatch == nil {
			e.newBatch.Wait()
		}
		IDToTrack = e.currentBatch.ID
		e.mu.Unlock()

		time.Sleep(e.waitDuration)

		e.mu.Lock()
		if e.currentBatch != nil &&
			e.currentBatch.ID == IDToTrack {
			e.readyToWork()
		}
		e.mu.Unlock()
	}
}

func (e *Engine) putToBatch(arg interface{}) BatchResult {
	e.mu.Lock()
	if e.currentBatch == nil {
		e.batchID++
		e.currentBatch = NewBatch(e.batchID)
		e.newBatch.Signal()
	}
	e.taskID++
	br := e.currentBatch.Put(e.taskID, arg)
	if e.currentBatch.argSize == e.argSizeLimit {
		e.readyToWork()
	}
	e.mu.Unlock()

	return br
}

// Submit puts arg to current batch to be worked on by background goroutine.
//
// This call will also waits for the results
func (e *Engine) Submit(arg interface{}) (interface{}, error) {
	br := e.putToBatch(arg)
	res, err := br.GetResult()
	br.batch = nil
	return res, err
}

// SubmitWithContext puts arg to current batch to be worked on by background goroutine,
// using golang's context idiom.
//
// Note that unless you need to use the context idiom, it is recommended
// to use `Submit()` call instead, as it has much, much less allocation (only interface{} typecasting).
// This API need to create another goroutine, and 2 channels to manage its functionality.
// (And of course, using context's `WithCancel` or `WithTimeout` also creates goroutines)
func (e *Engine) SubmitWithContext(
	ctx context.Context,
	arg interface{}) (interface{}, error) {

	// fast path, ctx already done
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	br := e.putToBatch(arg)

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
