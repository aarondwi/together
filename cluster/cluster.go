package cluster

import (
	"errors"
	"sync"

	e "github.com/aarondwi/together/engine"
	WP "github.com/aarondwi/together/workerpool"
)

// ErrPartitionNumberTooLow is returned
// when Cluster is not given enough numOfPartition, which is > 1.
//
// if <=1, then don't use cluster at all.
var ErrPartitionNumberTooLow = errors.New(
	"Cluster's partition number should be > 1")

// ErrPartitionNumOutOfRange is returned
// when the resulting partitionNum is outside slice's range
var ErrPartitionNumOutOfRange = errors.New(
	"Cluster's partition number should be in range [0, numOfPartition)")

// ErrNilPartitionerFunc is returned
// when non `ToPartition` function is called on nil partitioner
var ErrNilPartitionerFunc = errors.New(
	"nil partitioner func. Please use `SubmitToPartition` or `SubmitToPartitionWithContext` instead")

// Cluster allows you to scale together's engines on multi-core machine.
// While not perfect (cause using locks), with good partitioning scheme,
// it will scale to hella lots of goroutines submitting work.
//
// Note that this implementation is goroutine-safe, even without lock/atomic.
// This is because all variables can't (and won't, for now) be changed after creation, only be read.
//
// Note that lots of sync.Pool in the definition is used for `SubmitMany()`, as it allocates approximately (5 + 3 * numOfPartition) slices.
// All of them also return pointer to slices, as slice object is another allocation on heap
type Cluster struct {
	numOfPartition int
	engines        []*e.Engine
	partitioner    func(arg interface{}) int

	currentIdxForPartitionsPool sync.Pool
	positionsPool               sync.Pool

	temporaryResultSlicePool sync.Pool
	temporaryResultsPool     sync.Pool
	temporaryArgSlicePool    sync.Pool
	temporaryArgsPool        sync.Pool
}

// NewCluster creates our cluster.
//
// 3 last params are the exact same as a single engine,
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
	ec e.EngineConfig,
	fn e.WorkerFn,
	wp *WP.WorkerPool) (*Cluster, error) {

	if numOfPartition <= 1 {
		return nil, ErrPartitionNumberTooLow
	}

	engines := make([]*e.Engine, 0, numOfPartition)
	for i := 0; i < numOfPartition; i++ {
		e, err := e.NewEngine(
			e.EngineConfig{
				NumOfWorker:  ec.NumOfWorker,
				ArgSizeLimit: ec.ArgSizeLimit,
				WaitDuration: ec.WaitDuration},
			fn, wp)
		if err != nil {
			return nil, err
		}
		engines = append(engines, e)
	}

	return &Cluster{
		numOfPartition: numOfPartition,
		engines:        engines,
		partitioner:    partitioner,
		currentIdxForPartitionsPool: sync.Pool{
			// zero value of int is zero
			New: func() interface{} {
				s := make([]int, numOfPartition)
				return &s
			},
		},
		positionsPool: sync.Pool{
			New: func() interface{} {
				s := make([]int, 0, 64) // good buffer size
				return &s
			},
		},
		temporaryResultSlicePool: sync.Pool{
			New: func() interface{} {
				s := make([]*[]e.BatchResult, 0, numOfPartition)
				return &s
			},
		},
		temporaryResultsPool: sync.Pool{
			New: func() interface{} {
				s := make([]e.BatchResult, 0, 64)
				return &s
			},
		},
		temporaryArgSlicePool: sync.Pool{
			New: func() interface{} {
				s := make([]*[]interface{}, 0, numOfPartition)
				return &s
			},
		},
		temporaryArgsPool: sync.Pool{
			New: func() interface{} {
				s := make([]interface{}, 0, 64) // good buffer size
				return &s
			},
		},
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
		return e.EmptyBatchResult, ErrNilPartitionerFunc
	}
	return c.SubmitToPartition(c.partitioner(arg), arg)
}

// SubmitMany puts args into their correct partitions.
// To do so, it calls partitioner for each arg, and puts them into their respective temporary data
//
// There are 2 other ways to do so:
//
// 1. Do checking on each arg against each engine. This remove all needs to have intermediary state and so allocation, but need do to that while holding locks
//
// 2. The same as current way, but as bitmap. Only needs 1/8 of memory compared to this implementation, but harder to read
//
// This implementation currently is very, very pointer heavy, due to the need to reduce allocation (and thus GC pressure).
// Be sure not to check this while tired, high, or anything like that
func (c *Cluster) SubmitMany(
	args []interface{}) ([]e.BatchResult, error) {

	if c.partitioner == nil {
		return e.EmptyBatchResultSlice, ErrNilPartitionerFunc
	}

	// these variables are used
	// to track which arg got where
	// as we need to return the BatchResult in the same order as args
	//
	// This is done to avoid looping for every args and every partition, which is (numOfPartition * len(args))
	currentIdxForPartitions := c.currentIdxForPartitionsPool.Get().(*[]int)
	positions := c.positionsPool.Get().(*[]int)

	temporaryResults := c.temporaryResultSlicePool.Get().(*([]*[]e.BatchResult))
	for i := 0; i < c.numOfPartition; i++ {
		*temporaryResults = append(*temporaryResults, c.temporaryResultsPool.Get().(*[]e.BatchResult))
	}

	temporaryArgs := c.temporaryArgSlicePool.Get().(*([]*[]interface{}))
	for i := 0; i < c.numOfPartition; i++ {
		*temporaryArgs = append(*temporaryArgs, c.temporaryArgsPool.Get().(*[]interface{}))
	}

	// this one can't be pooled transparently, need caller (user) to participate
	result := make([]e.BatchResult, 0, len(args))

	// assign each arg to the correct partition
	for _, arg := range args {
		idx := c.partitioner(arg)
		*(*temporaryArgs)[idx] = append(*(*temporaryArgs)[idx], arg)
		*positions = append(*positions, idx)
	}
	for i, t := range *temporaryArgs {
		c.engines[i].SubmitManyInto(*t, (*temporaryResults)[i])
	}

	// this is the reordering part
	// just as a reminder, positions hold index of which partition
	for i := 0; i < len(args); i++ {
		result = append(result,
			(*(*temporaryResults)[(*positions)[i]])[(*currentIdxForPartitions)[(*positions)[i]]])
		(*currentIdxForPartitions)[(*positions)[i]]++
	}

	/////////////////////////////////////////////////////////////
	// Return all to pools
	/////////////////////////////////////////////////////////////
	for i := 0; i < c.numOfPartition; i++ {
		(*currentIdxForPartitions)[i] = 0
	}
	c.currentIdxForPartitionsPool.Put(currentIdxForPartitions)

	*positions = (*positions)[:0]
	c.positionsPool.Put(positions)

	for i := 0; i < c.numOfPartition; i++ {
		*(*temporaryArgs)[i] = (*(*temporaryArgs)[i])[:0]
		c.temporaryArgsPool.Put((*temporaryArgs)[i])
	}
	*temporaryArgs = (*temporaryArgs)[:0]
	c.temporaryArgSlicePool.Put(temporaryArgs)

	for i := 0; i < c.numOfPartition; i++ {
		*(*temporaryResults)[i] = (*(*temporaryResults)[i])[:0]
		c.temporaryResultsPool.Put((*temporaryResults)[i])
	}
	*temporaryResults = (*temporaryResults)[:0]
	c.temporaryResultSlicePool.Put(temporaryResults)
	/////////////////////////////////////////////////////////////

	return result, nil
}

// SubmitToPartition puts arg into engine number `partitionNum`.
func (c *Cluster) SubmitToPartition(
	partitionNum int,
	arg interface{}) (e.BatchResult, error) {

	if partitionNum < 0 || partitionNum >= c.numOfPartition {
		return e.EmptyBatchResult, ErrPartitionNumOutOfRange
	}
	return c.engines[partitionNum].Submit(arg), nil
}

// SubmitManyToPartition puts bunch of args at once into engine number `partitionNum`.
func (c *Cluster) SubmitManyToPartition(
	partitionNum int,
	args []interface{}) ([]e.BatchResult, error) {

	if partitionNum < 0 || partitionNum >= c.numOfPartition {
		return e.EmptyBatchResultSlice, ErrPartitionNumOutOfRange
	}
	return c.engines[partitionNum].SubmitMany(args), nil
}
