package engine

import (
	"errors"
	"sync"
	"time"

	WP "github.com/aarondwi/together/workerpool"
)

// ErrEngineBrokenSizes is returned when the soft/hard limit is relatively broken
//
// Please check EngineConfig docs for more
var ErrEngineBrokenSizes = errors.New(
	"Please check EngineConfig docs, and ensure your config match the requirements")

// ErrResultNotFound is returned
// when an ID is not found in the results map
var ErrResultNotFound = errors.New(
	"A result is not found on the resulting map. " +
		"Please check your code to ensure all ids are returned with their corresponding results.")

// ErrNilWorkerFn is returned when`workerFn` is nil
var ErrNilWorkerFn = errors.New("workerFn can't be nil")

// ErrArgsBiggerThanHardLimit is returned when the number of given args is bigger than hard limit
var ErrArgsBiggerThanHardLimit = errors.New(
	"number of args passed is bigger than the given hard limit")

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
	mu                  *sync.Mutex
	newBatch            *sync.Cond
	batchReady          *sync.Cond
	batchSlotsAvailable *sync.Cond
	batchID             uint64
	taskID              uint64
	fn                  WorkerFn
	numOfWorker         int

	// we move from channel, as it is bounded, can cause a deadlock when can't insert anymore to buffer
	// we either need to back to before, which only track the time/softLimit, or going unbounded
	batchChan    llBatch
	batchWaiting int // we track number of batches in batchChan, so can apply custom backpressure
	currentBatch *Batch
	softLimit    int
	hardLimit    int
	waitDuration time.Duration
	wp           *WP.WorkerPool
}

// EngineConfig is the config object for our engine
//
// SoftLimit is the limit where the batch can already be taken by worker,
// but not forcefully sent to. Worker gonna take this by itself after a timeout
// (triggered by timeoutWatchdog)
//
// While HardLimit is which can't be passed. Useful for example like AWS SQS,
// which a single call can only have at most 100 messages
//
// If SoftLimit is empty, it is set to HardLimit / 2.
//
// If HardLimit is empty, it is set to SoftLimit * 2.
//
// After those conversion, both should still be bigger than 1, and SoftLimit <= HardLimit
type EngineConfig struct {
	NumOfWorker  int
	SoftLimit    int
	HardLimit    int
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
	if fn == nil {
		return nil, ErrNilWorkerFn
	}

	// only change if not set
	var softLimit, hardLimit int
	if ec.SoftLimit == 0 {
		hardLimit = ec.HardLimit
		softLimit = ec.HardLimit / 2
	} else if ec.HardLimit == 0 {
		softLimit = ec.SoftLimit
		hardLimit = ec.SoftLimit * 2
	} else {
		softLimit = ec.SoftLimit
		hardLimit = ec.HardLimit
	}

	if (softLimit > hardLimit) || (softLimit <= 1) || (hardLimit <= 1) {
		return nil, ErrEngineBrokenSizes
	}

	mu := &sync.Mutex{}
	newBatch := sync.NewCond(mu)
	batchReady := sync.NewCond(mu)
	batchSlotsAvailable := sync.NewCond(mu)
	e := &Engine{
		mu:                  mu,
		newBatch:            newBatch,
		batchReady:          batchReady,
		batchSlotsAvailable: batchSlotsAvailable,
		fn:                  fn,
		batchChan:           llBatch{},
		numOfWorker:         ec.NumOfWorker,
		softLimit:           softLimit,
		hardLimit:           hardLimit,
		waitDuration:        ec.WaitDuration,
		wp:                  wp,
	}
	for i := 0; i < e.numOfWorker; i++ {
		go e.worker()
	}
	go e.timeoutWatchdog()
	return e, nil
}

func (e *Engine) worker() {
	for {
		e.mu.Lock()
		for (e.currentBatch == nil ||
			(e.currentBatch.argSize < e.softLimit && time.Since(e.currentBatch.createdAt) < e.waitDuration)) &&
			e.batchWaiting == 0 {
			e.batchReady.Wait()
		}
		if e.batchWaiting == 0 { // a new batch has passed softLimit
			e.readyToWork()
		}
		b := e.batchChan.get()
		if b == nil {
			panic("Shouldn't be nil after reaching here, `Wait` code must be broken")
		}
		e.batchWaiting--
		if e.batchWaiting < e.numOfWorker {
			e.batchSlotsAvailable.Broadcast()
		}
		e.mu.Unlock()

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
	e.batchChan.add(e.currentBatch)
	e.batchWaiting++
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
	for {
		e.mu.Lock()
		for e.currentBatch == nil {
			e.newBatch.Wait()
		}
		IDToTrack := e.currentBatch.ID
		e.mu.Unlock()

		time.Sleep(e.waitDuration)

		e.mu.Lock()
		if e.currentBatch != nil &&
			e.currentBatch.ID == IDToTrack {
			e.batchReady.Signal()
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

// this function should only be called while holding the mutex
func (e *Engine) waitBatchSlotsAvailable() {
	for e.batchWaiting > e.numOfWorker {
		e.batchSlotsAvailable.Wait()
	}
}

// This function should only be called when holding the mutex
func (e *Engine) checkSoftLimitReadiness() {
	if e.currentBatch.argSize >= e.softLimit {
		e.batchReady.Signal()
	}
}

// Submit puts arg to current batch to be worked on by background goroutine.
func (e *Engine) Submit(arg interface{}) BatchResult {
	e.mu.Lock()
	e.waitBatchSlotsAvailable()

	e.ensureBatch()
	if e.currentBatch.argSize >= e.hardLimit { // cause only gonna increase by 1
		e.readyToWork()
		e.batchReady.Signal()
		e.ensureBatch()
	}

	e.taskID++
	br := e.currentBatch.Put(e.taskID, arg)

	e.checkSoftLimitReadiness()
	e.mu.Unlock()

	return br
}

// SubmitMany puts args into batches, which may be one or more
//
// Atomicity should not be assumed, as this implementation will try to pack as much as possible
// to ensure we get the best possible savings + efficiency.
//
// It would allocate a slice to accomodate the result, and internally call `SubmitManyInto`
func (e *Engine) SubmitMany(args []interface{}) []BatchResult {
	result := make([]BatchResult, 0, len(args))
	e.SubmitManyInto(args, &result)
	return result
}

// SubmitManyInto puts args into batches, which may be one or more, and store them into result
//
// Atomicity should not be assumed, as this implementation will try to pack as much as possible
// to ensure we get the best possible savings + efficiency.
//
// Mainly used to control result's slice allocation
//
// Note that the result will be appended to, so can (and will) be re-allocated for more spaces
func (e *Engine) SubmitManyInto(args []interface{}, result *[]BatchResult) {
	l := len(args)
	e.mu.Lock()
	e.waitBatchSlotsAvailable()
	for l > 0 {
		e.ensureBatch()
		numArgsTaken := l
		emptySlots := e.hardLimit - e.currentBatch.argSize
		if emptySlots < numArgsTaken {
			numArgsTaken = emptySlots
		}

		startPos := len(args) - l
		for _, arg := range args[startPos : startPos+numArgsTaken] {
			e.taskID++
			br := e.currentBatch.Put(e.taskID, arg)
			*result = append(*result, br)
		}

		if e.currentBatch.argSize >= e.hardLimit {
			e.readyToWork()
			e.batchReady.Signal()
		} else {
			e.checkSoftLimitReadiness()
		}
		l -= numArgsTaken
	}
	e.mu.Unlock()
}
