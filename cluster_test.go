package together

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestClusterValidation(t *testing.T) {
	_, err := NewCluster(1, nil, 2, 10, time.Duration(time.Second), nil)
	if err == nil || err != ErrPartitionNumberTooLow {
		log.Fatal("Should return ErrPartitionNumberTooLow cause only given 1, but it is not")
	}

	_, err = NewCluster(2, nil, 2, 10, time.Duration(time.Second), nil)
	if err == nil || err != ErrNilWorkerFn {
		log.Fatal("Should return ErrNilWorkerFn cause given nil, but it is not")
	}
}

func TestClusterSubmitNilPartitionerFn(t *testing.T) {
	c, err := NewCluster(
		// cluster params
		2, nil,
		// per-engine param
		2, 10, time.Duration(time.Second),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		})
	if err != nil {
		log.Fatalf("Should not fail, but we got %v", err)
	}

	_, err = c.Submit(1)
	if err == nil || err != ErrNilPartitionerFunc {
		log.Fatal("Should error because nil partitioner func, but it is not")
	}

	_, err = c.SubmitWithContext(context.Background(), 1)
	if err == nil || err != ErrNilPartitionerFunc {
		log.Fatal("Should error because nil partitioner func, but it is not")
	}
}

func TestClusterSubmitOutsideRange(t *testing.T) {
	c, err := NewCluster(
		// cluster params
		2, func(arg interface{}) int {
			return 2
		},
		// per-engine param
		2, 10, time.Duration(time.Second),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		})
	if err != nil {
		log.Fatalf("Should not fail, but we got %v", err)
	}

	// without toPartition idiom
	_, err = c.Submit(1)
	if err == nil || err != ErrPartitionNumOutOfRange {
		log.Fatal("Should error because out of range, but it is not")
	}

	_, err = c.SubmitWithContext(context.Background(), 1)
	if err == nil || err != ErrPartitionNumOutOfRange {
		log.Fatal("Should error because out of range, but it is not")
	}

	// with toPartition idiom
	_, err = c.SubmitToPartition(2, 1)
	if err == nil || err != ErrPartitionNumOutOfRange {
		log.Fatal("Should error because out of range, but it is not")
	}

	_, err = c.SubmitToPartitionWithContext(context.Background(), 2, 1)
	if err == nil || err != ErrPartitionNumOutOfRange {
		log.Fatal("Should error because out of range, but it is not")
	}
}

func TestClusterSubmit(t *testing.T) {
	c, err := NewCluster(
		// cluster params
		2, func(arg interface{}) int {
			return arg.(int)
		},
		// per-engine param
		2, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		})
	if err != nil {
		log.Fatalf("Should not fail, but we got %v", err)
	}

	// without toPartition idiom
	res, err := c.Submit(1)
	if err != nil {
		log.Fatalf("Should not error, but instead we got %v", err)
	}
	if res.(int) != 1 {
		log.Fatalf("Should be equal to 1, but instead we got %d", res.(int))
	}

	res, err = c.SubmitWithContext(context.Background(), 1)
	if err != nil {
		log.Fatalf("Should not error, but instead we got %v", err)
	}
	if res.(int) != 1 {
		log.Fatalf("Should be equal to 1, but instead we got %d", res.(int))
	}

	// with toPartition idiom
	res, err = c.SubmitToPartition(0, 1)
	if err != nil {
		log.Fatalf("Should not error, but instead we got %v", err)
	}
	if res.(int) != 1 {
		log.Fatalf("Should be equal to 1, but instead we got %d", res.(int))
	}

	res, err = c.SubmitToPartitionWithContext(context.Background(), 0, 1)
	if err != nil {
		log.Fatalf("Should not error, but instead we got %v", err)
	}
	if res.(int) != 1 {
		log.Fatalf("Should be equal to 1, but instead we got %d", res.(int))
	}
}