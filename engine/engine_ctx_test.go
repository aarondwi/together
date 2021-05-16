package engine

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	com "github.com/aarondwi/together/common"
)

func TestEngineWithCtx(t *testing.T) {
	var wp, _ = com.NewWorkerPool(4, 10)
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
			res, err := br.GetResultWithContext(
				context.Background())
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
	/*
	 * Notes because how go handles local variable,
	 * if you are using multiple `GetResultWithContext` in single function,
	 * be sure to assign it to different local variable.
	 *
	 * See https://stackoverflow.com/questions/25919213/why-does-go-handle-closures-differently-in-goroutines
	 * for details
	 */
	var wp, _ = com.NewWorkerPool(4, 10)
	ErrTest := errors.New("")
	e, err := NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			// gives us time to cancel the ctx first
			time.Sleep(10 * time.Millisecond)
			return nil, ErrTest
		}, wp)
	if err != nil {
		log.Fatalf("It should not error, cause all correct, but got %v", err)
	}

	// we can assign to br, but for consistency later on,
	// we start with 1
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	br1 := e.Submit(10)
	_, err = br1.GetResultWithContext(ctx)
	// should not be ErrTest, because ctx is cancelled first
	if err == nil || err == ErrTest {
		log.Fatal("Should not receive ErrTest, cause ctx already cancelled, but it is")
	}

	// can't assign to br, cause go will consider it to point to same object, causing race condition
	// so, just continue the numbering
	ctx, cancelFunc = context.WithCancel(context.Background())
	go func() {
		time.Sleep(2 * time.Millisecond)
		cancelFunc()
	}()
	br2 := e.Submit(10)
	_, err = br2.GetResultWithContext(ctx)
	// should not be ErrTest, because ctx is cancelled first
	if err == nil || err == ErrTest {
		log.Fatal("Should not receive ErrTest, cause ctx got cancelled first, but it is")
	}

	// same as before, now at 3
	br3 := e.Submit(10)
	_, err = br3.GetResultWithContext(context.Background())
	if err == nil || err != ErrTest {
		log.Fatalf("Should receive ErrTest, but instead we got %v", err)
	}
}
