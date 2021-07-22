package engine

import (
	"errors"
	"sync"
	"time"

	com "github.com/aarondwi/together/common"
)

// ErrArgSizeLimitLessThanEqualOne is returned
// when given argSizeLimit <= 1
var ErrArgSizeLimitLessThanEqualOne = errors.New(
	"argSizeLimit is expected to be > 1")

// ErrResultNotFound is returned
// when an ID is not found in the results map
var ErrResultNotFound = errors.New(
	"A result is not found on the resulting map. " +
		"Please check your code to ensure all ids are returned with their corresponding results.")

// ErrNilWorkerFn is returned when`workerFn` is nil
var ErrNilWorkerFn = errors.New("workerFn can't be nil")

type WorkerFn func(map[uint64]interface{}) (map[uint64]interface{}, error)

// Engine is our batch-controller, loosely adapted from
// https://github.com/grab/async/blob/master/batch.go,
// designed specifically for business-logic use-case.
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
	fn           WorkerFn
	batchChan    chan *Batch
	currentBatch *Batch
	argSizeLimit int
	waitDuration time.Duration
	wp           *com.WorkerPool
}

// EngineConfig is the config object for our engine
type EngineConfig struct {
	NumOfWorker  int
	ArgSizeLimit int
	WaitDuration time.Duration
}

// NewEngine creates the engine based on the given EngineConfig
//
// The given `fn` should be goroutine-safe
func NewEngine(
	ec EngineConfig,
	fn WorkerFn,
	wp *com.WorkerPool) (*Engine, error) {
	if ec.NumOfWorker <= 0 {
		return nil, com.ErrNumOfWorkerLessThanEqualZero
	}
	if ec.ArgSizeLimit <= 1 {
		return nil, ErrArgSizeLimitLessThanEqualOne
	}
	if fn == nil {
		return nil, ErrNilWorkerFn
	}
	mu := &sync.Mutex{}
	newBatch := sync.NewCond(mu)
	e := &Engine{
		mu:       mu,
		newBatch: newBatch,
		fn:       fn,

		// we allow one buffer for each worker
		batchChan:    make(chan *Batch, ec.NumOfWorker),
		argSizeLimit: ec.ArgSizeLimit,
		waitDuration: ec.WaitDuration,
		wp:           wp,
	}
	for i := 0; i < ec.NumOfWorker; i++ {
		go e.worker()
	}
	go e.timeoutWatchdog()
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
// This is a conscious decision, to not balloon the memory requirement
func (e *Engine) readyToWork() {
	e.batchChan <- e.currentBatch
	e.currentBatch = nil
}

// this implementation may skip some id to check
// but that means the `Submit()` call is much faster
// than the wait.
//
// In that case, no need to wait for those batches.
func (e *Engine) timeoutWatchdog() {
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

// Submit puts arg to current batch to be worked on by background goroutine.
func (e *Engine) Submit(arg interface{}) BatchResult {
	e.mu.Lock()
	if e.currentBatch == nil {
		e.batchID++
		e.currentBatch = NewBatch(e.batchID, e.wp)
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
