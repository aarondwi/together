package cluster

import (
	"math/rand"
	"testing"
	"time"

	e "github.com/aarondwi/together/engine"
	WP "github.com/aarondwi/together/workerpool"
)

func TestClusterValidation(t *testing.T) {
	var wp, _ = WP.NewWorkerPool(4, 10, false)
	_, err := NewCluster(
		1, nil,
		e.EngineConfig{
			NumOfWorker:  2,
			SoftLimit:    10,
			WaitDuration: time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != ErrPartitionNumberTooLow {
		t.Fatal("Should return ErrPartitionNumberTooLow cause only given 1, but it is not")
	}

	_, err = NewCluster(
		2, nil,
		e.EngineConfig{
			NumOfWorker:  2,
			SoftLimit:    10,
			WaitDuration: time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != e.ErrNilWorkerFn {
		t.Fatal("Should return ErrNilWorkerFn cause given nil, but it is not")
	}
}

func TestClusterSubmitNilPartitionerFn(t *testing.T) {
	c, err := NewCluster(
		2, nil,
		e.EngineConfig{
			NumOfWorker:  2,
			SoftLimit:    10,
			WaitDuration: time.Duration(time.Second)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		}, nil)
	if err != nil {
		t.Fatalf("Should not fail, but we got %v", err)
	}

	_, err = c.Submit(1)
	if err == nil || err != ErrNilPartitionerFunc {
		t.Fatal("Should error because nil partitioner func, but it is not")
	}
}

func TestClusterSubmitOutsideRange(t *testing.T) {
	c, err := NewCluster(
		2, func(arg interface{}) int {
			return 2
		},
		e.EngineConfig{
			NumOfWorker:  2,
			SoftLimit:    10,
			WaitDuration: time.Duration(time.Second)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		}, nil)
	if err != nil {
		t.Fatalf("Should not fail, but we got %v", err)
	}

	_, err = c.Submit(1)
	if err == nil || err != ErrPartitionNumOutOfRange {
		t.Fatal("Should error because out of range, but it is not")
	}

	_, err = c.SubmitToPartition(2, 1)
	if err == nil || err != ErrPartitionNumOutOfRange {
		t.Fatal("Should error because out of range, but it is not")
	}
}

func TestClusterSubmit(t *testing.T) {
	c, err := NewCluster(
		2, func(arg interface{}) int {
			return arg.(int)
		},
		e.EngineConfig{
			NumOfWorker:  2,
			SoftLimit:    10,
			WaitDuration: time.Duration(5 * time.Millisecond)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		}, nil)
	if err != nil {
		t.Fatalf("Should not fail, but we got %v", err)
	}

	br, err := c.Submit(1)
	if err != nil {
		t.Fatalf("Should not error, but instead we got %v", err)
	}
	res, err := br.GetResult()
	if err != nil {
		t.Fatalf("Should not error, but instead we got %v", err)
	}
	if res.(int) != 1 {
		t.Fatalf("Should be equal to 1, but instead we got %d", res.(int))
	}

	br, err = c.SubmitToPartition(0, 2)
	if err != nil {
		t.Fatalf("Should not error, but instead we got %v", err)
	}
	res, err = br.GetResult()
	if err != nil {
		t.Fatalf("Should not error, but instead we got %v", err)
	}
	if res.(int) != 2 {
		t.Fatalf("Should be equal to 2, but instead we got %d", res.(int))
	}
}

func TestClusterSubmitMany(t *testing.T) {
	c, err := NewCluster(
		2, func(arg interface{}) int {
			return rand.Intn(1)
		},
		e.EngineConfig{
			NumOfWorker:  2,
			SoftLimit:    10,
			WaitDuration: time.Duration(5 * time.Millisecond)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		}, nil)
	if err != nil {
		t.Fatalf("Should not fail, but we got %v", err)
	}

	requests := make([]interface{}, 0, 64)
	for i := 0; i < 64; i++ {
		requests = append(requests, i)
	}
	brs, err := c.SubmitMany(requests)
	if err != nil {
		t.Fatalf("Should not error, but instead we got %v", err)
	}

	for i, br := range brs {
		r, err := br.GetResult()
		if err != nil {
			t.Fatalf("Should have no error here, but instead we got %v", err)
		}
		if i != r.(int) {
			t.Fatalf("Should be the same, but instead we got %d and %d", i, r.(int))
		}
	}
}

func TestClusterSubmitManyNilPartitioner(t *testing.T) {
	c, err := NewCluster(
		2, nil,
		e.EngineConfig{
			NumOfWorker:  2,
			SoftLimit:    10,
			WaitDuration: time.Duration(5 * time.Millisecond)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		}, nil)
	if err != nil {
		t.Fatalf("Should not fail, but we got %v", err)
	}

	requests := make([]interface{}, 0, 64)
	for i := 0; i < 64; i++ {
		requests = append(requests, i)
	}
	_, err = c.SubmitMany(requests)
	if err == nil || err != ErrNilPartitionerFunc {
		t.Fatal("Should error because nil partitioner func, but it is not")
	}
}
