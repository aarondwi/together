package workerpool

import "errors"

// ErrNumberOfWorkerLessThanEqualZero is returned
// when given numOfWorker is <= zero
var ErrNumberOfWorkerLessThanEqualZero = errors.New(
	"numOfWorker is expected to be > 0")

// ErrNumberOfWorkerPartitionTooLow is returned
// when Cluster is not given enough numOfPartition, which is > 1.
//
// if <=1, then don't use cluster at all.
var ErrNumberOfWorkerPartitionTooLow = errors.New(
	"Cluster's partition number should be > 1")

// WorkerPool represents our background workers
// which gonna do the waiting for context idiom or combiner implementation.
//
// The goal is to have few hot, already allocated goroutines. `allowRuntimeCreation` manages this behavior
//
// The recommendation is to have ~10K goroutines here.
// Assuming each goroutine takes 2KB (cause just a waiting goroutine, very hard to get bigger),
// even 10K only uses 20MB of memory
type WorkerPool struct {
	partitioner          func(arg interface{}) int
	channels             []chan func()
	allowRuntimeCreation bool
}

// NewWorkerPool creates our WorkerPool object.
// The total number of workers is `numOfPartition * numOfWorker`
//
// If only 1 partition, this is just an unnecessary abstraction,
// and will return error
func NewWorkerPool(
	numOfPartition int,
	numOfWorker int,
	allowRuntimeCreation bool) (*WorkerPool, error) {

	if numOfPartition <= 1 {
		return nil, ErrNumberOfWorkerPartitionTooLow
	}
	if numOfWorker <= 0 {
		return nil, ErrNumberOfWorkerLessThanEqualZero
	}
	channels := make([]chan func(), 0, numOfPartition)
	for i := 0; i < numOfPartition; i++ {
		ch := make(chan func(), numOfWorker)
		channels = append(channels, ch)
	}

	for i := 0; i < numOfPartition; i++ {
		for j := 0; j < numOfWorker; j++ {
			go func(ch chan func()) {
				for fn := range ch {
					fn()
				}
			}(channels[i])
		}
	}

	return &WorkerPool{
		partitioner:          GetDefaultPartitioner(numOfPartition),
		channels:             channels,
		allowRuntimeCreation: allowRuntimeCreation,
	}, nil
}

// Submit select a random channel, and then put fn there
//
// For now, if the given channel is full AND `allowRuntimeCreation`, create a new goroutine instead.
// In the future, will change to somehow track and reuse new goroutines.
func (wp *WorkerPool) Submit(fn func()) {
	if !wp.allowRuntimeCreation {
		wp.channels[wp.partitioner(nil)] <- fn
	} else {
		select {
		case wp.channels[wp.partitioner(nil)] <- fn:
			return
		default:
			if wp.allowRuntimeCreation {
				go fn()
			}
		}
	}
}
