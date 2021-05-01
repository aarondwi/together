package together

import (
	"errors"
	"log"
	"sync"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	valShouldFail := 13
	globalCount := 0
	e, err := NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		// notes that in real usage
		// usually you won't just doing in-memory operation
		// but rather, doing a network call
		// and network call is much more expensive than just locking + memory ops
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			globalCount++
			res := make(map[uint64]interface{})
			log.Printf("===== %d =====", globalCount)
			for k, v := range m {
				log.Printf("%d: %d", k, v)
			}
			for k, v := range m {
				if v.(int) != valShouldFail {
					res[k] = v.(int) * 2
				} else {
					log.Printf("Now, we get %d", v)
				}
			}
			return res, nil
		})
	if err != nil {
		log.Fatalf("It should not error, cause all correct, but got %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(24)
	for i := 0; i < 24; i++ {
		go func(j int) {
			res, err := e.Submit(j)
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
	ErrTest := errors.New("")
	e, err := NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		// notes that in real usage
		// usually you won't just doing in-memory operation
		// but rather, doing a network call
		// and network call is much more expensive than just locking + memory ops
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		})
	if err != nil {
		log.Fatalf("It should not error, cause all correct, but got %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(24)
	for i := 0; i < 24; i++ {
		go func(j int) {
			_, err := e.Submit(j)
			if err == nil || err != ErrTest {
				log.Fatal("Should receive ErrTest, but it is not")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestEngineValidation(t *testing.T) {
	_, err := NewEngine(-1, 10, time.Duration(time.Second), nil)
	if err == nil || err != ErrNumOfWorkerLessThanEqualZero {
		log.Fatal("Should fail cause numOfWorker <= 0, but it is not")
	}

	_, err = NewEngine(1, -1, time.Duration(time.Second), nil)
	if err == nil || err != ErrArgSizeLimitLessThanEqualOne {
		log.Fatal("Should fail cause argSizeLimit <= 1, but it is not")
	}
}