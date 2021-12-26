package engine

import (
	"errors"
	"sync"
	"time"

	WP "github.com/aarondwi/together/workerpool"
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
// Internally, this implementation uses a slice, to hold inflight data.
// In the future, if the need arises, will add map-based batch.
// Useful typically for case where multiple keys could be the same.
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
	wp           *WP.WorkerPool
}

// EngineConfig is the config object for our engine
//
// Note that ArgSizeLimit is only a soft limit.
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
	wp *WP.WorkerPool) (*Engine, error) {
	if ec.NumOfWorker <= 0 {
		return nil, WP.ErrNumberOfWorkerLessThanEqualZero
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
// This is a conscious decision, to not balloon the memory requirement, at least for now.
func (e *Engine) readyToWork() {
	e.batchChan <- e.currentBatch
	e.currentBatch = nil
}

// This implementation may skip some id to check
// but that means the `Submit()` call is much faster than the wait.
// In that case, no need to wait for those batches.
//
// We do not need to precisely wait this
// as we are `Wait`-ing on the condition variable, which should be precise enough for this use case.
// The diff with the precise checking should only be on low single digit ms.
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

// This function should only be called when holding the mutex
func (e *Engine) ensureBatch() {
	if e.currentBatch == nil {
		e.batchID++
		e.currentBatch = NewBatch(e.batchID, e.wp)
		e.newBatch.Signal()
	}
}

// This function should only be called when holding the mutex
func (e *Engine) checkReadiness() {
	if e.currentBatch.argSize >= e.argSizeLimit {
		e.readyToWork()
	}
}

// Submit puts arg to current batch to be worked on by background goroutine.
func (e *Engine) Submit(arg interface{}) BatchResult {
	e.mu.Lock()
	e.ensureBatch()

	e.taskID++
	br := e.currentBatch.Put(e.taskID, arg)

	e.checkReadiness()
	e.mu.Unlock()

	return br
}

// SubmitMany puts bunch of args into current batch.
//
// It would allocate a slice on hot path to accomodate the result
func (e *Engine) SubmitMany(args []interface{}) []BatchResult {
	// create slices outside of mutex, reduce holding time (even only ns)
	result := make([]BatchResult, 0, len(args))

	e.mu.Lock()
	e.ensureBatch()

	for _, arg := range args {
		e.taskID++
		br := e.currentBatch.Put(e.taskID, arg)
		result = append(result, br)
	}

	e.checkReadiness()
	e.mu.Unlock()

	return result
}

// SubmitManyInto puts args to current batch, and store BatchResult into result
//
// Mainly used to control result's slice allocation
//
// Note that the result will be appended to, so will (and can) be re-allocated for more spaces
func (e *Engine) SubmitManyInto(args []interface{}, result *[]BatchResult) {
	e.mu.Lock()
	e.ensureBatch()

	for _, arg := range args {
		e.taskID++
		br := e.currentBatch.Put(e.taskID, arg)
		*result = append(*result, br)
	}

	e.checkReadiness()
	e.mu.Unlock()
}
