package cluster

import (
	"log"
	"testing"
	"time"

	com "github.com/aarondwi/together/common"
	e "github.com/aarondwi/together/engine"
)

func TestClusterValidation(t *testing.T) {
	var wp, _ = com.NewWorkerPool(4, 10)
	_, err := NewCluster(
		1, nil,
		e.EngineConfig{
			NumOfWorker:  2,
			ArgSizeLimit: 10,
			WaitDuration: time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != com.ErrPartitionNumberTooLow {
		log.Fatal("Should return ErrPartitionNumberTooLow cause only given 1, but it is not")
	}

	_, err = NewCluster(
		2, nil,
		e.EngineConfig{
			NumOfWorker:  2,
			ArgSizeLimit: 10,
			WaitDuration: time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != e.ErrNilWorkerFn {
		log.Fatal("Should return ErrNilWorkerFn cause given nil, but it is not")
	}
}

func TestClusterSubmitNilPartitionerFn(t *testing.T) {
	c, err := NewCluster(
		2, nil,
		e.EngineConfig{
			NumOfWorker:  2,
			ArgSizeLimit: 10,
			WaitDuration: time.Duration(time.Second)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		}, nil)
	if err != nil {
		log.Fatalf("Should not fail, but we got %v", err)
	}

	_, err = c.Submit(1)
	if err == nil || err != com.ErrNilPartitionerFunc {
		log.Fatal("Should error because nil partitioner func, but it is not")
	}
}

func TestClusterSubmitOutsideRange(t *testing.T) {
	c, err := NewCluster(
		2, func(arg interface{}) int {
			return 2
		},
		e.EngineConfig{
			NumOfWorker:  2,
			ArgSizeLimit: 10,
			WaitDuration: time.Duration(time.Second)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		}, nil)
	if err != nil {
		log.Fatalf("Should not fail, but we got %v", err)
	}

	_, err = c.Submit(1)
	if err == nil || err != com.ErrPartitionNumOutOfRange {
		log.Fatal("Should error because out of range, but it is not")
	}

	_, err = c.SubmitToPartition(2, 1)
	if err == nil || err != com.ErrPartitionNumOutOfRange {
		log.Fatal("Should error because out of range, but it is not")
	}
}

func TestClusterSubmit(t *testing.T) {
	c, err := NewCluster(
		2, func(arg interface{}) int {
			return arg.(int)
		},
		e.EngineConfig{
			NumOfWorker:  2,
			ArgSizeLimit: 10,
			WaitDuration: time.Duration(5 * time.Millisecond)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return m, nil
		}, nil)
	if err != nil {
		log.Fatalf("Should not fail, but we got %v", err)
	}

	br, err := c.Submit(1)
	if err != nil {
		log.Fatalf("Should not error, but instead we got %v", err)
	}
	res, err := br.GetResult()
	if err != nil {
		log.Fatalf("Should not error, but instead we got %v", err)
	}
	if res.(int) != 1 {
		log.Fatalf("Should be equal to 1, but instead we got %d", res.(int))
	}

	br, err = c.SubmitToPartition(0, 2)
	if err != nil {
		log.Fatalf("Should not error, but instead we got %v", err)
	}
	res, err = br.GetResult()
	if err != nil {
		log.Fatalf("Should not error, but instead we got %v", err)
	}
	if res.(int) != 2 {
		log.Fatalf("Should be equal to 2, but instead we got %d", res.(int))
	}
}
