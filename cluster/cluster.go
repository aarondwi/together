package cluster

import (
	"time"

	com "github.com/aarondwi/together/common"
	e "github.com/aarondwi/together/engine"
)

// Cluster allows you to scale together's engines on multi-core machine.
// While not perfect (cause using locks), with good partitioning scheme,
// it will scale to hella lots of goroutines submitting work.
//
// Note that this implementation is goroutine-safe, even without lock/atomic.
// This is because all variables can't (and won't) be changed after creation, only be read.
type Cluster struct {
	numOfPartition int
	engines        []*e.Engine
	partitioner    func(arg interface{}) int
}

// NewCluster creates our cluster.
//
// 5 last params are the exact same as a single engine,
// and directly applied to each engine
//
// Notes for `partitioner` param, the given function
// should return value in range [0, numOfPartition),
// and it should be goroutine-safe.
//
// Note that the WorkerPool is shared to all engines
func NewCluster(
	// cluster params
	numOfPartition int,
	partitioner func(arg interface{}) int,

	// engine params
	numOfWorker int,
	argSizeLimit int,
	waitDuration time.Duration,
	fn e.WorkerFn,
	wp *com.WorkerPool) (*Cluster, error) {

	if numOfPartition <= 1 {
		return nil, com.ErrPartitionNumberTooLow
	}

	engines := make([]*e.Engine, 0, numOfPartition)
	for i := 0; i < numOfPartition; i++ {
		e, err := e.NewEngine(
			numOfWorker, argSizeLimit,
			waitDuration, fn, wp)
		if err != nil {
			return nil, err
		}
		engines = append(engines, e)
	}

	return &Cluster{
		numOfPartition: numOfPartition,
		engines:        engines,
		partitioner:    partitioner,
	}, nil
}

// Submit selects and puts arg into one of the engines.
//
// This call will redirect calls to `SubmitToPartition` via `partitioner`
//
// Should not be called if you passed nil to workerFn
func (c *Cluster) Submit(
	arg interface{}) (e.BatchResult, error) {

	if c.partitioner == nil {
		return e.EmptyBatchResult, com.ErrNilPartitionerFunc
	}
	return c.SubmitToPartition(c.partitioner(arg), arg)
}

// Submit selects and puts arg into engine number `partitionNum`.
func (c *Cluster) SubmitToPartition(
	partitionNum int,
	arg interface{}) (e.BatchResult, error) {

	if partitionNum < 0 || partitionNum >= c.numOfPartition {
		return e.EmptyBatchResult, com.ErrPartitionNumOutOfRange
	}
	return c.engines[partitionNum].Submit(arg), nil
}
