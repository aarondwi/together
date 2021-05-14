package together

import (
	"context"
	"errors"
	"log"
	"testing"
)

/*
 * Note for this test,
 * we are NOT setting wg at all
 * cause wg.Wait() only wait for zero,
 * and we can assume wg.Done() is called by `Engine`
 */

var ErrTest = errors.New("Errors for test")

func TestBatchGetResult(t *testing.T) {
	b := &Batch{
		ID:      1,
		args:    map[uint64]interface{}{10: 2},
		argSize: 10,
		results: map[uint64]interface{}{10: 2},
	}

	br := BatchResult{id: 10, batch: b}
	res, err := br.GetResult()
	if err != nil {
		log.Fatalf("should not error because no error, but we got %v", err)
	}
	if res.(int) != 2 {
		log.Fatalf("We should receive 2, but instead we got %d", res.(int))
	}

	br = BatchResult{id: 11, batch: b}
	_, err = br.GetResult()
	if err == nil || err != ErrResultNotFound {
		log.Fatalf("err should be ErrResultNotFound, but instead we got %v", err)
	}

	b.err = ErrTest
	br = BatchResult{id: 11, batch: b}
	_, err = br.GetResult()
	if err == nil || err != ErrTest {
		log.Fatalf("err should be ErrTest, but instead we got %v", err)
	}
}

func TestBatchGetResultToChan(t *testing.T) {
	b := &Batch{
		ID:      1,
		args:    map[uint64]interface{}{10: 2},
		argSize: 10,
		results: map[uint64]interface{}{10: 2},
	}

	br := BatchResult{id: 10, batch: b}
	resCh1 := make(chan interface{}, 1)
	errCh1 := make(chan error, 1)
	br.GetResultToChan(resCh1, errCh1)
	select {
	case err := <-errCh1:
		log.Fatalf("should not error because no error, but we got %v", err)
	case res := <-resCh1:
		if res.(int) != 2 {
			log.Fatalf("We should receive 2, but instead we got %d", res.(int))
		}
	}

	br = BatchResult{id: 11, batch: b}
	resCh2 := make(chan interface{}, 1)
	errCh2 := make(chan error, 1)
	br.GetResultToChan(resCh2, errCh2)
	select {
	case err := <-errCh2:
		if err == nil || err != ErrResultNotFound {
			log.Fatalf("err should be ErrResultNotFound, but instead we got %v", err)
		}
	case res := <-resCh2:
		log.Fatalf("should receive nothing, but instead we got %v", res)
	}

	b.err = ErrTest
	br = BatchResult{id: 11, batch: b}
	resCh3 := make(chan interface{}, 1)
	errCh3 := make(chan error, 1)
	br.GetResultToChan(resCh3, errCh3)
	select {
	case err := <-errCh3:
		if err == nil || err != ErrTest {
			log.Fatalf("err should be ErrTest, but instead we got %v", err)
		}
	case res := <-resCh3:
		log.Fatalf("should receive nothing, but instead we got %v", res)
	}
}

func TestBatchGetResultWithCtx(t *testing.T) {
	b := &Batch{
		ID:      1,
		args:    map[uint64]interface{}{10: 2},
		argSize: 10,
		results: map[uint64]interface{}{10: 2},
	}

	br := BatchResult{id: 10, batch: b}
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, err := br.GetResultWithContext(ctx)
	if err == nil {
		log.Fatal("should return ctx.Err(), but it is not")
	}

	br = BatchResult{id: 10, batch: b}
	res, err := br.GetResultWithContext(context.Background())
	if err != nil {
		log.Fatalf("should not error because no error, but we got %v", err)
	}
	if res.(int) != 2 {
		log.Fatalf("We should receive 2, but instead we got %d", res.(int))
	}

	b.err = ErrTest
	br = BatchResult{id: 11, batch: b}
	_, err = br.GetResultWithContext(context.Background())
	if err == nil || err != ErrTest {
		log.Fatalf("err should be ErrTest, but instead we got %v", err)
	}
}

func TestBatchGetResultWithCtxToChan(t *testing.T) {
	b := &Batch{
		ID:      1,
		args:    map[uint64]interface{}{10: 2},
		argSize: 10,
		results: map[uint64]interface{}{10: 2},
	}

	br := BatchResult{id: 10, batch: b}
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	resCh1 := make(chan interface{}, 1)
	errCh1 := make(chan error, 1)
	br.GetResultWithContextToChan(ctx, resCh1, errCh1)
	select {
	case err := <-errCh1:
		if err == nil {
			log.Fatal("should return ctx.Err(), but it is not")
		}
	case res := <-resCh1:
		log.Fatalf("We should not get any result, cause ctx error, but we got %v", res)
	}

	br = BatchResult{id: 10, batch: b}
	resCh2 := make(chan interface{}, 1)
	errCh2 := make(chan error, 1)
	br.GetResultWithContextToChan(context.Background(), resCh2, errCh2)
	select {
	case err := <-errCh2:
		if err != nil {
			log.Fatalf("should not error because no error, but we got %v", err)
		}
	case res := <-resCh2:
		if res.(int) != 2 {
			log.Fatalf("We should receive 2, but instead we got %d", res.(int))
		}
	}

	b.err = ErrTest
	br = BatchResult{id: 11, batch: b}
	resCh3 := make(chan interface{}, 1)
	errCh3 := make(chan error, 1)
	br.GetResultWithContextToChan(context.Background(), resCh3, errCh3)
	select {
	case err := <-errCh3:
		if err == nil || err != ErrTest {
			log.Fatalf("err should be ErrTest, but instead we got %v", err)
		}
	case res := <-resCh3:
		log.Fatalf("should receive nothing, but instead we got %v", res)
	}
}
