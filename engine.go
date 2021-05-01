package together

import (
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

// this function should only be called
// when holding the mutex
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

// Submit arg to current batch to be worked on by background goroutine.
//
// This call will also waits for the results
func (e *Engine) Submit(arg interface{}) (interface{}, error) {
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

	res, err := br.GetResult()
	br.batch = nil
	return res, err
}
