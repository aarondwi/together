package together

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"
)

func TestEngineWithCtx(t *testing.T) {
	valShouldFail := 18
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
			res, err := e.SubmitWithContext(context.Background(), j)
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

func TestEngineCtxReturnsError(t *testing.T) {
	ErrTest := errors.New("")
	e, err := NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			// gives us time to cancel the ctx first
			time.Sleep(10 * time.Millisecond)
			return nil, ErrTest
		})
	if err != nil {
		log.Fatalf("It should not error, cause all correct, but got %v", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, err = e.SubmitWithContext(ctx, 10)
	// should not be ErrTest, because ctx is cancelled first
	if err == nil || err == ErrTest {
		log.Fatal("Should not receive ErrTest, but it is not")
	}

	ctx, cancelFunc = context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Millisecond)
		cancelFunc()
	}()
	_, err = e.SubmitWithContext(ctx, 10)
	// should not be ErrTest, because ctx is cancelled first
	if err == nil || err == ErrTest {
		log.Fatal("Should not receive ErrTest, but it is not")
	}

	_, err = e.SubmitWithContext(context.Background(), 10)
	if err == nil || err != ErrTest {
		log.Fatal("Should receive ErrTest, but it is not")
	}
}
