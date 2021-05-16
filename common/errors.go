package common

import "errors"

/*
 * This package has multiple common errors
 * used by this and other sub-packages
 */

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

// ErrNumOfWorkerLessThanEqualZero is returned
// when given numOfWorker is <= zero
var ErrNumOfWorkerLessThanEqualZero = errors.New(
	"numOfWorker is expected to be > 0")

// ErrNilWorkerPool is returned
// when function that needs background goroutine to wait
// is called without initializing worker pool
var ErrNilWorkerPool = errors.New(
	"Nil Worker pool. Can't use this functionality")
