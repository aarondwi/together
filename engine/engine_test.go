package engine

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	WP "github.com/aarondwi/together/workerpool"
)

func TestEngine(t *testing.T) {
	var wp, _ = WP.NewWorkerPool(4, 10, false)
	valShouldFail := 13
	var globalCount uint32
	e, err := NewEngine(
		EngineConfig{1, 10, 20, time.Duration(5 * time.Millisecond)},
		// notes that in real usage
		// usually you won't just doing in-memory operation
		// but rather, doing a network call
		// and network call is much more expensive than just locking + memory ops
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			atomic.AddUint32(&globalCount, 1)
			res := make(map[uint64]interface{})
			for k, v := range m {
				if v.(int) != valShouldFail {
					res[k] = v.(int) * 2
				}
			}
			return res, nil
		}, wp)
	if err != nil {
		t.Fatalf("It should not error, cause all correct, but got %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(24)
	for i := 0; i < 24; i++ {
		go func(j int) {
			br := e.Submit(j)
			res, err := br.GetResult()
			if j == valShouldFail {
				if err == nil || err != ErrResultNotFound {
					t.Fatalf("Submit with arg %d should fail, but we got %v, with error %v", valShouldFail, res, err)
				}
			} else {
				if err != nil {
					t.Fatalf(
						"Call with arg %d should not fail, but we got %v",
						j, err)
				}
				if res.(int) != j*2 {
					t.Fatalf(
						"Call with arg %d should return %d, but we got %d",
						j, j*2, res.(int))
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	// depends when workers took from batch
	val := atomic.LoadUint32(&globalCount)
	if !(val == 2 || val == 3) {
		t.Fatalf("batch should be called 2 or 3 times, depending on worker timing, but we got %d", atomic.LoadUint32(&globalCount))
	}
}

func TestEngineReturnsError(t *testing.T) {
	var wp, _ = WP.NewWorkerPool(4, 10, false)
	ErrTest := errors.New("")
	e, err := NewEngine(
		EngineConfig{1, 10, 20, time.Duration(5 * time.Millisecond)},
		// notes that in real usage
		// usually you won't just doing in-memory operation
		// but rather, doing a network call
		// and network call is much more expensive than just locking + memory ops
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)
	if err != nil {
		t.Fatalf("It should not error, cause all correct, but got %v", err)
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
	var wp, _ = WP.NewWorkerPool(4, 10, false)
	_, err := NewEngine(
		EngineConfig{-1, 10, 20, time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != WP.ErrNumberOfWorkerLessThanEqualZero {
		log.Fatal("Should fail cause numOfWorker <= 0, but it is not")
	}

	_, err = NewEngine(
		EngineConfig{1, 2, 20, time.Duration(time.Second)},
		nil, wp)
	if err == nil || err != ErrNilWorkerFn {
		log.Fatal("Should fail cause nil workerFn, but it is not")
	}

	placeholderFunc := func(m map[uint64]interface{}) (map[uint64]interface{}, error) { return m, nil }

	_, err = NewEngine(
		EngineConfig{1, 20, 10, time.Duration(time.Second)},
		placeholderFunc, wp)
	if err == nil || err != ErrEngineBrokenSizes {
		t.Fatalf("Should fail cause softLimit > hardLimit, but it is not, with error %v", err)
	}

	_, err = NewEngine(
		EngineConfig{1, -1, 10, time.Duration(time.Second)},
		placeholderFunc, wp)
	if err == nil || err != ErrEngineBrokenSizes {
		t.Fatalf("Should fail cause softLimit <= 1, but it is not, with error %v", err)
	}

	_, err = NewEngine(
		EngineConfig{1, 0, 10, time.Duration(time.Second)},
		placeholderFunc, wp)
	if err != nil {
		t.Fatalf("Should be okay cause should be softLimit=5 and hardLimit=10, but instead we got %v", err)
	}

	_, err = NewEngine(
		EngineConfig{1, 10, 20, time.Duration(time.Second)},
		placeholderFunc, wp)
	if err != nil {
		t.Fatalf("Should be okay cause all correct, but instead we got %v", err)
	}
}

func TestEngineSubmitManyPlusDeadlock(t *testing.T) {
	// these tests are combined cause both test same behavior
	// the deadlock case happen when the number of batch pass the number of channel buffer
	errX := errors.New("errX is just a test error")
	valShouldFail := 8
	e, _ := NewEngine(
		EngineConfig{1, 10, 20, time.Duration(5 * time.Millisecond)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			res := make(map[uint64]interface{})
			for k, v := range m {
				if v.(int) != valShouldFail {
					res[k] = v.(int) * 2
				} else {
					res[k] = errX
				}
			}
			return res, nil
		}, nil)

	requests := make([]interface{}, 0, 100)
	for i := 0; i < 100; i++ {
		requests = append(requests, i)
	}
	result := e.SubmitMany(requests)
	for i, res := range result {
		r, err := res.GetResult()
		if i != valShouldFail && err != nil {
			t.Fatalf("It should be only returning error when i == %d with err %v", valShouldFail, err)
		}
		if i != valShouldFail && r.(int) != i*2 {
			t.Fatalf("It should return the same, but instead we got %d and %d", r.(int), i)
		}
	}
}

func TestEngineNoWorkerAvailableWaitHardLimit(t *testing.T) {
	var globalCount uint32
	e, _ := NewEngine(
		EngineConfig{1, 10, 20, time.Duration(5 * time.Millisecond)},
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			atomic.AddUint32(&globalCount, 1)
			time.Sleep(20 * time.Millisecond)
			return m, nil
		}, nil)

	requests := make([]interface{}, 0, 14)
	for i := 0; i < 14; i++ {
		requests = append(requests, i)
	}
	e.SubmitMany(requests)

	requests2 := make([]interface{}, 0, 9)
	for i := 0; i < 9; i++ {
		requests2 = append(requests2, i)
	}
	e.SubmitMany(requests2)

	requests3 := make([]interface{}, 0, 9)
	for i := 0; i < 9; i++ {
		requests3 = append(requests3, i)
	}
	result := e.SubmitMany(requests3)
	for _, res := range result {
		_, err := res.GetResult()
		if err != nil {
			t.Fatalf("Should have no error, but instead we got %v", err)
		}
	}

	if atomic.LoadUint32(&globalCount) != 2 {
		t.Fatalf("It should only be called 2 times, cause only 1 worker, but instead we got %d",
			atomic.LoadUint32(&globalCount))
	}
}
