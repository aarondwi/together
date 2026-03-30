package workerpool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolValidation(t *testing.T) {
	_, err := NewWorkerPool(1, 2, false)
	if err != ErrNumberOfWorkerPartitionTooLow {
		t.Fatalf("It should return `ErrNumberOfWorkerPartitionTooLow`, but instead we got %v", err)
	}

	_, err = NewWorkerPool(2, 0, false)
	if err != ErrNumberOfWorkerLessThanEqualZero {
		t.Fatalf("It should return `ErrNumberOfWorkerLessThanEqualZero`, but instead we got %v", err)
	}
}

func TestWorkerPoolNoRuntimeCreation(t *testing.T) {
	wp, _ := NewWorkerPool(2, 4, false)

	var a int32
	var wg sync.WaitGroup

	wg.Add(64)
	for i := 0; i < 64; i++ {
		wp.Submit(func() {
			time.Sleep(1 * time.Millisecond)
			atomic.AddInt32(&a, 1)
			wg.Done()
		})
	}

	wg.Wait()
	if atomic.LoadInt32(&a) != 64 {
		t.Fatalf("It should already be 64, but instead we got %d", a)
	}
}

func TestWorkerPoolAllowRuntimeCreation(t *testing.T) {
	wp, _ := NewWorkerPool(2, 16, true)

	var a int32
	var wg sync.WaitGroup

	wg.Add(64)
	for i := 0; i < 64; i++ {
		wp.Submit(func() {
			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&a, 1)
			wg.Done()
		})
	}

	wg.Wait()
	if atomic.LoadInt32(&a) != 64 {
		t.Fatalf("It should already be 66, but instead we got %d", a)
	}
}
