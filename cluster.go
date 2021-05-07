package together

import (
	"context"
	"errors"
	"time"
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

// ErrNilWhichPartitionFunc is returned
// when non `ToPartition` function is called when whichPartition is nil
var ErrNilWhichPartitionFunc = errors.New(
	"nil whichPartition func. Please use `SubmitToPartition` or `SubmitToPartitionWithContext` instead")

// Cluster allows you to scale engines on multi-core machine.
// While not perfect (cause using locks), with good partitioning scheme,
// it will scale to hella lots of goroutines submitting work.
//
// Note that this implementation is goroutine-safe, even without lock/atomic.
// This is because all variables can't (and won't) be changed after creation, only be read.
type Cluster struct {
	numOfPartition int
	engines        []*Engine
	whichPartition func(arg interface{}) int
}

// NewCluster creates our cluster.
//
// 4 last params are the exact same as a single engine,
// and directly applied to each engine
//
// Notes for `whichPartition` param, the given function
// should return value in range [0, numOfPartition),
// and it should be goroutine-safe
func NewCluster(
	// cluster params
	numOfPartition int,
	whichPartition func(arg interface{}) int,

	// engine params
	numOfWorker int,
	argSizeLimit int,
	waitDuration time.Duration,
	fn WorkerFn) (*Cluster, error) {

	if numOfPartition <= 1 {
		return nil, ErrPartitionNumberTooLow
	}

	engines := make([]*Engine, 0, numOfPartition)
	for i := 0; i < numOfPartition; i++ {
		e, err := NewEngine(
			numOfWorker, argSizeLimit,
			waitDuration, fn)
		if err != nil {
			return nil, err
		}
		engines = append(engines, e)
	}

	return &Cluster{
		numOfPartition: numOfPartition,
		engines:        engines,
		whichPartition: whichPartition,
	}, nil
}

// Submit selects and puts arg into one of the engines.
//
// This call will internally calls `whichPartition` func,
// to decide which partition to send arg to.
//
// Should not be called if you passed nil to workerFn
func (c *Cluster) Submit(
	arg interface{}) (interface{}, error) {

	if c.whichPartition == nil {
		return nil, ErrNilWhichPartitionFunc
	}
	partitionNum := c.whichPartition(arg)
	if partitionNum < 0 || partitionNum >= c.numOfPartition {
		return nil, ErrPartitionNumOutOfRange
	}

	return c.engines[partitionNum].Submit(arg)
}

// Submit selects and puts arg into engine number `partitionNum`.
func (c *Cluster) SubmitToPartition(
	partitionNum int,
	arg interface{}) (interface{}, error) {

	if partitionNum < 0 || partitionNum >= c.numOfPartition {
		return nil, ErrPartitionNumOutOfRange
	}

	return c.engines[partitionNum].Submit(arg)
}

// SubmitWithContext selects and puts arg into engine number `partitionNum`.
//
// This call will internally calls `whichPartition` func,
// to decide which partition to send arg to.
//
// It is recommended to use `Submit` instead, if you don't need
// the context idiom. It has less allocation, so it is faster.
//
// Should not be called if you passed nil to workerFn
func (c *Cluster) SubmitWithContext(
	ctx context.Context,
	arg interface{}) (interface{}, error) {

	if c.whichPartition == nil {
		return nil, ErrNilWhichPartitionFunc
	}
	partitionNum := c.whichPartition(arg)
	if partitionNum < 0 || partitionNum >= c.numOfPartition {
		return nil, ErrPartitionNumOutOfRange
	}

	return c.engines[partitionNum].SubmitWithContext(ctx, arg)
}

// SubmitToPartitionWithContext selects and puts arg into engine number `partitionNum`.
//
// It is recommended to use `SubmitToPartition` instead, if you don't need
// the context idiom. It has less allocation, so it is faster.
func (c *Cluster) SubmitToPartitionWithContext(
	ctx context.Context,
	partitionNum int,
	arg interface{}) (interface{}, error) {

	if partitionNum < 0 || partitionNum >= c.numOfPartition {
		return nil, ErrPartitionNumOutOfRange
	}

	return c.engines[partitionNum].SubmitWithContext(ctx, arg)
}
