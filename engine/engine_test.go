package engine

import (
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	com "github.com/aarondwi/together/common"
)

func TestEngine(t *testing.T) {
	var wp, _ = com.NewWorkerPool(4, 10)
	valShouldFail := 13
	globalCount := 0
	e, err := NewEngine(
		EngineConfig{1, 10, time.Duration(5 * time.Millisecond)},
		// notes that in real usage
		// usually you won't just doing in-memory operation
		// but rather, doing a network call
		// and network call is much more expensive than just locking + memory ops
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			globalCount++
			res := make(map[uint64]interface{})
			for k, v := range m {
				if v.(int) != valShouldFail {
					res[k] = v.(int) * 2
				}
			}
			return res, nil
		}, wp)
	if err != nil {
		log.Fatalf("It should not error, cause all correct, but got %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(24)
	for i := 0; i < 24; i++ {
		go func(j int) {
			br := e.Submit(j)
			res, err := br.GetResult()
			if j == valShouldFail {
				if err == nil || err != ErrResultNotFound {
					log.Fatalf("Submit with arg %d should fail, but we got %v, with error %v", valShouldFail, res, err)
				}
			} else {
				if err != nil {
					log.Fatalf(
						"Call with arg %d should not fail, but we got %v",
						j, err)
				}
				if res.(int) != j*2 {
					log.Fatalf(
						"Call with arg %d should return %d, but we got %d",
						j, j*2, res.(int))
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if globalCount != 3 {
		log.Fatalf("batch should be called 3 times, but we got %d", globalCount)
	}
}

func TestEngineReturnsError(t *testing.T) {
	var wp, _ = com.NewWorkerPool(4, 10)
	ErrTest := errors.New("")
	e, err := NewEngine(
		EngineConfig{1, 10, time.Duration(5 * time.Millisecond)},
		// notes that in real usage
		// usually you won't just doing in-memory operation
		// but rather, doing a network call
		// and network call is much more expensive than just locking + memory ops
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)
	if err != nil {
		log.Fatalf("It should not error, cause all correct, but got %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(24)
	for i := 0; i < 24; i++ {
		go func(j int) {
			br := e.Submit(j)
			_, err := br.GetResult()
			if err == nil || err != ErrTest {
				log.Fatal("Should receive ErrTest, but it is not")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestEngineValidation(t *testing.T) {
	var wp, _ = com.NewWorkerPool(4, 10)
	_, err := NewEngine(
		EngineConfig{-1, 10, time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != com.ErrNumOfWorkerLessThanEqualZero {
		log.Fatal("Should fail cause numOfWorker <= 0, but it is not")
	}

	_, err = NewEngine(
		EngineConfig{1, -1, time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != ErrArgSizeLimitLessThanEqualOne {
		log.Fatal("Should fail cause argSizeLimit <= 1, but it is not")
	}

	_, err = NewEngine(
		EngineConfig{1, 2, time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != ErrNilWorkerFn {
		log.Fatal("Should fail cause nil workerFn, but it is not")
	}
}
