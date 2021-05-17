package common

// WorkerPool represents our background workers
// which gonna do the waiting for context idiom function or combiner.
//
// The goal is to not allocate goroutines to wait multiple paths.
type WorkerPool struct {
	partitioner func(arg interface{}) int
	channels    []chan func()
}

// NewWorkerPool creates our WorkerPool object.
// The total number of workers is `numOfPartition * numOfWorker`
func NewWorkerPool(
	numOfPartition int,
	numOfWorker int) (*WorkerPool, error) {

	if numOfPartition <= 1 {
		return nil, ErrPartitionNumberTooLow
	}
	if numOfWorker <= 0 {
		return nil, ErrNumOfWorkerLessThanEqualZero
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
		partitioner: GetDefaultPartitioner(numOfPartition),
		channels:    channels,
	}, nil
}

// Submit select a random channel, and then put fn there
func (wp *WorkerPool) Submit(fn func()) {
	wp.channels[wp.partitioner(nil)] <- fn
}
