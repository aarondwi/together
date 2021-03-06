package engine

import (
	"context"
	"errors"
	"testing"

	WP "github.com/aarondwi/together/workerpool"
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
		wp:      nil,
	}

	br := BatchResult{id: 10, batch: b}
	res, err := br.GetResult()
	if err != nil {
		t.Fatalf("should not error because no error, but we got %v", err)
	}
	if res.(int) != 2 {
		t.Fatalf("We should receive 2, but instead we got %d", res.(int))
	}

	br = BatchResult{id: 11, batch: b}
	_, err = br.GetResult()
	if err == nil || err != ErrResultNotFound {
		t.Fatalf("err should be ErrResultNotFound, but instead we got %v", err)
	}

	b.err = ErrTest
	br = BatchResult{id: 11, batch: b}
	_, err = br.GetResult()
	if err == nil || err != ErrTest {
		t.Fatalf("err should be ErrTest, but instead we got %v", err)
	}
}

func TestBatchGetResultWithCtx(t *testing.T) {
	var wp, _ = WP.NewWorkerPool(2, 1, false)
	b := &Batch{
		ID:      1,
		args:    map[uint64]interface{}{10: 2},
		argSize: 5,
		results: map[uint64]interface{}{10: 2},
		wp:      wp,
	}

	br1 := BatchResult{id: 10, batch: b}
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, err := br1.GetResultWithContext(ctx)
	if err == nil {
		t.Fatal("should return ctx.Err(), but it is not")
	}

	br2 := BatchResult{id: 10, batch: b}
	res, err := br2.GetResultWithContext(context.Background())
	if err != nil {
		t.Fatalf("should not error because no error, but we got %v", err)
	}
	if res.(int) != 2 {
		t.Fatalf("We should receive 2, but instead we got %d", res.(int))
	}

	b.err = ErrTest
	br3 := BatchResult{id: 11, batch: b}
	_, err = br3.GetResultWithContext(context.Background())
	if err == nil || err != ErrTest {
		t.Fatalf("err should be ErrTest, but instead we got %v", err)
	}

	b.wp = nil
	b.err = nil
	br4 := BatchResult{id: 10, batch: b}
	res, err = br4.GetResultWithContext(context.Background())
	if err != nil {
		t.Fatalf("should not error because no error, but we got %v", err)
	}
	if res.(int) != 2 {
		t.Fatalf("We should receive 2, but instead we got %d", res.(int))
	}
}

func TestBatchResultIsError(t *testing.T) {
	errX := errors.New("error X is for test only")

	b := &Batch{
		ID:      1,
		args:    map[uint64]interface{}{10: errX},
		argSize: 10,
		results: map[uint64]interface{}{10: errX},
		wp:      nil,
	}

	br1 := BatchResult{id: 10, batch: b}
	_, err := br1.GetResult()
	if err == nil || err != errX {
		t.Fatalf("Should return `errX`, but we got %v", err)
	}

	br2 := BatchResult{id: 10, batch: b}
	_, err = br2.GetResultWithContext(context.Background())
	if err == nil || err != errX {
		t.Fatalf("Should return `errX`, but we got %v", err)
	}
}
