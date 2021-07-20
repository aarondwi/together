package combiner

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	com "github.com/aarondwi/together/common"
	e "github.com/aarondwi/together/engine"
)

var wp, _ = com.NewWorkerPool(4, 10)
var ErrTest = errors.New("")

func TestAllSuccesses(t *testing.T) {
	c := NewCombiner(wp)

	res, err := c.AllSuccesses([]e.BatchResult{})
	if err != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, err)
	}

	e1, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {

			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, wp)

	res1 := e1.Submit(1)
	res2 := e1.Submit(2)
	res3 := e1.Submit(3)
	res4 := e1.Submit(4)

	resArr, err := c.AllSuccesses([]e.BatchResult{res1, res2, res3, res4})
	if err != nil {
		log.Fatalf("It should not error, cause all ok, but we got %v", err)
	}
	if len(resArr) != 4 {
		log.Fatalf("It should only be 4, cause we submit 4 times, but instead we got %v", resArr...)
	}
	if resArr[0].(int) != 2 || resArr[1].(int) != 4 ||
		resArr[2].(int) != 6 || resArr[3].(int) != 8 {
		log.Fatalf("It should be 2, 4, 6, 8, but instead we got %v", resArr...)
	}

	e2, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)
	res5 := e1.Submit(3)
	res6 := e2.Submit(4)

	_, err = c.AllSuccesses([]e.BatchResult{res5, res6})
	if err == nil || err != ErrTest {
		log.Fatalf("It should return error ErrTest, but instead we got %v", err)
	}
}

func TestAllSuccessesNilWorkerPool(t *testing.T) {
	c := NewCombiner(nil)

	res, err := c.AllSuccesses([]e.BatchResult{})
	if err != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, err)
	}

	e1, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {

			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, nil)

	res1 := e1.Submit(1)
	res2 := e1.Submit(2)
	res3 := e1.Submit(3)
	res4 := e1.Submit(4)

	resArr, err := c.AllSuccesses([]e.BatchResult{res1, res2, res3, res4})
	if err != nil {
		log.Fatalf("It should not error, cause all ok, but we got %v", err)
	}
	if len(resArr) != 4 {
		log.Fatalf("It should only be 4, cause we submit 4 times, but instead we got %v", resArr...)
	}
	if resArr[0].(int) != 2 || resArr[1].(int) != 4 ||
		resArr[2].(int) != 6 || resArr[3].(int) != 8 {
		log.Fatalf("It should be 2, 4, 6, 8, but instead we got %v", resArr...)
	}

	e2, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)
	res5 := e1.Submit(3)
	res6 := e2.Submit(4)

	_, err = c.AllSuccesses([]e.BatchResult{res5, res6})
	if err == nil || err != ErrTest {
		log.Fatalf("It should return error ErrTest, but instead we got %v", err)
	}
}

func TestAllSuccessesWithCtx(t *testing.T) {
	c := NewCombiner(wp)

	res, err := c.AllSuccessesWithContext(context.Background(), []e.BatchResult{})
	if err != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, err = c.AllSuccessesWithContext(ctx, []e.BatchResult{})
	if err == nil {
		log.Fatal("Should return Err Cancelled, but instead we got nil")
	}

	e1, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, wp)

	res1 := e1.Submit(1)
	res2 := e1.Submit(2)
	res3 := e1.Submit(3)
	res4 := e1.Submit(4)

	resArr, err := c.AllSuccessesWithContext(
		context.Background(), []e.BatchResult{res1, res2, res3, res4})
	if err != nil {
		log.Fatalf("It should not error, cause all ok, but we got %v", err)
	}
	if len(resArr) != 4 {
		log.Fatalf("It should only be 4, cause we submit 4 times, but instead we got %v", resArr...)
	}
	if resArr[0].(int) != 2 || resArr[1].(int) != 4 ||
		resArr[2].(int) != 6 || resArr[3].(int) != 8 {
		log.Fatalf("It should be 2 and 4, but instead we got %v", resArr...)
	}

	e2, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)
	res5 := e1.Submit(3)
	res6 := e2.Submit(4)

	_, err = c.AllSuccessesWithContext(context.Background(), []e.BatchResult{res5, res6})
	if err == nil || err != ErrTest {
		log.Fatalf("It should return error ErrTest, but instead we got %v", err)
	}
}

func TestAllSuccessesWithCtxNilWorkerPool(t *testing.T) {
	c := NewCombiner(nil)

	res, err := c.AllSuccessesWithContext(context.Background(), []e.BatchResult{})
	if err != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, err = c.AllSuccessesWithContext(ctx, []e.BatchResult{})
	if err == nil {
		log.Fatal("Should return Err Cancelled, but instead we got nil")
	}

	e1, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, nil)

	res1 := e1.Submit(1)
	res2 := e1.Submit(2)
	res3 := e1.Submit(3)
	res4 := e1.Submit(4)

	resArr, err := c.AllSuccessesWithContext(
		context.Background(), []e.BatchResult{res1, res2, res3, res4})
	if err != nil {
		log.Fatalf("It should not error, cause all ok, but we got %v", err)
	}
	if len(resArr) != 4 {
		log.Fatalf("It should only be 4, cause we submit 4 times, but instead we got %v", resArr...)
	}
	if resArr[0].(int) != 2 || resArr[1].(int) != 4 ||
		resArr[2].(int) != 6 || resArr[3].(int) != 8 {
		log.Fatalf("It should be 2 and 4, but instead we got %v", resArr...)
	}

	e2, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)
	res5 := e1.Submit(3)
	res6 := e2.Submit(4)

	_, err = c.AllSuccessesWithContext(context.Background(), []e.BatchResult{res5, res6})
	if err == nil || err != ErrTest {
		log.Fatalf("It should return error ErrTest, but instead we got %v", err)
	}
}

func TestRace(t *testing.T) {
	c := NewCombiner(wp)

	res, errs := c.Race([]e.BatchResult{})
	if errs != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, errs)
	}
	e1, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, wp)
	e2, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)
	e3, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 3
			}
			return res, nil
		}, wp)

	ErrTest2 := errors.New("")
	e4, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(5 * time.Millisecond)
			return nil, ErrTest2
		}, wp)

	br1 := e1.Submit(3)
	br2 := e2.Submit(4)
	br3 := e3.Submit(5)

	res, errs = c.Race([]e.BatchResult{br1, br2, br3})
	if errs != nil {
		log.Fatalf("errs should be nil, but instead we got %v", errs)
	}
	if res.(int) != 15 {
		log.Fatalf("We should receive result for `br3` which is 15, but instead we got %v", res)
	}

	br4 := e4.Submit(7)
	br5 := e2.Submit(8)
	_, errs = c.Race([]e.BatchResult{br4, br5})
	if len(errs) != 2 {
		log.Fatalf("errs should have len 2, but instead we got %v", errs)
	}
	if !(errs[0] == ErrTest2 && errs[1] == ErrTest) {
		log.Fatalf("Receive wrong errs: %v", errs)
	}
}

func TestRaceNilWorkerPool(t *testing.T) {
	c := NewCombiner(nil)

	res, errs := c.Race([]e.BatchResult{})
	if errs != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, errs)
	}
	e1, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, nil)
	e2, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, nil)
	e3, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 3
			}
			return res, nil
		}, nil)

	ErrTest2 := errors.New("")
	e4, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(5 * time.Millisecond)
			return nil, ErrTest2
		}, nil)

	br1 := e1.Submit(3)
	br2 := e2.Submit(4)
	br3 := e3.Submit(5)

	res, errs = c.Race([]e.BatchResult{br1, br2, br3})
	if errs != nil {
		log.Fatalf("errs should be nil, but instead we got %v", errs)
	}
	if res.(int) != 15 {
		log.Fatalf("We should receive result for `br3` which is 15, but instead we got %v", res)
	}

	br4 := e4.Submit(7)
	br5 := e2.Submit(8)
	_, errs = c.Race([]e.BatchResult{br4, br5})
	if len(errs) != 2 {
		log.Fatalf("errs should have len 2, but instead we got %v", errs)
	}
	if !(errs[0] == ErrTest2 && errs[1] == ErrTest) {
		log.Fatalf("Receive wrong errs: %v", errs)
	}
}

func TestRaceWithCtx(t *testing.T) {
	c := NewCombiner(wp)

	res, errs := c.RaceWithContext(
		context.Background(), []e.BatchResult{})
	if errs != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, errs)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, errs = c.RaceWithContext(ctx, []e.BatchResult{})
	if errs == nil || len(errs) != 0 {
		log.Fatalf("errs have len zero, but instead we got %v", errs)
	}

	e1, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, wp)
	e2, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)
	e3, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 3
			}
			return res, nil
		}, wp)

	ErrTest2 := errors.New("")
	e4, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(5 * time.Millisecond)
			return nil, ErrTest2
		}, wp)

	br1 := e1.Submit(3)
	br2 := e2.Submit(4)
	br3 := e3.Submit(5)

	res, errs = c.RaceWithContext(
		context.Background(), []e.BatchResult{br1, br2, br3})
	if errs != nil {
		log.Fatalf("errs should be nil, but instead we got %v", errs)
	}
	if res.(int) != 15 {
		log.Fatalf("We should receive result for `br3` which is 15, but instead we got %v", res)
	}

	br4 := e4.Submit(7)
	br5 := e2.Submit(8)
	_, errs = c.RaceWithContext(
		context.Background(), []e.BatchResult{br4, br5})
	if len(errs) != 2 {
		log.Fatalf("errs should have len 2, but instead we got %v", errs)
	}
	if !(errs[0] == ErrTest2 && errs[1] == ErrTest) {
		log.Fatalf("Receive wrong errs: %v", errs)
	}
}

func TestRaceWithCtxNilWorkerPool(t *testing.T) {
	c := NewCombiner(nil)

	res, errs := c.RaceWithContext(
		context.Background(), []e.BatchResult{})
	if errs != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, errs)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, errs = c.RaceWithContext(ctx, []e.BatchResult{})
	if errs == nil || len(errs) != 0 {
		log.Fatalf("errs have len zero, but instead we got %v", errs)
	}

	e1, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(10 * time.Millisecond)
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, nil)
	e2, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, nil)
	e3, _ := e.NewEngine(
		1, 10, time.Duration(5*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(1 * time.Millisecond)
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 3
			}
			return res, nil
		}, nil)

	ErrTest2 := errors.New("")
	e4, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			time.Sleep(5 * time.Millisecond)
			return nil, ErrTest2
		}, nil)

	br1 := e1.Submit(3)
	br2 := e2.Submit(4)
	br3 := e3.Submit(5)

	res, errs = c.RaceWithContext(
		context.Background(), []e.BatchResult{br1, br2, br3})
	if errs != nil {
		log.Fatalf("errs should be nil, but instead we got %v", errs)
	}
	if res.(int) != 15 {
		log.Fatalf("We should receive result for `br3` which is 15, but instead we got %v", res)
	}

	br4 := e4.Submit(7)
	br5 := e2.Submit(8)
	_, errs = c.RaceWithContext(
		context.Background(), []e.BatchResult{br4, br5})
	if len(errs) != 2 {
		log.Fatalf("errs should have len 2, but instead we got %v", errs)
	}
	if !(errs[0] == ErrTest2 && errs[1] == ErrTest) {
		log.Fatalf("Receive wrong errs: %v", errs)
	}
}

func TestEvery(t *testing.T) {
	c := NewCombiner(wp)

	res, errs := c.Every([]e.BatchResult{})
	if errs != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, errs)
	}
	e1, _ := e.NewEngine(
		1, 10, time.Duration(2*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, wp)
	e2, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)

	br1 := e1.Submit(3)
	br2 := e2.Submit(4)
	br3 := e1.Submit(5)
	br4 := e2.Submit(6)
	br5 := e1.Submit(7)

	res, errs = c.Every([]e.BatchResult{br1, br2, br3, br4, br5})
	if len(res) != 5 || len(errs) != 5 {
		log.Fatalf("errs should be nil, but instead we got %v", errs)
	}
	if !(errs[0] == nil && errs[1] == ErrTest &&
		errs[2] == nil && errs[3] == ErrTest && errs[4] == nil) {
		log.Fatalf("Receive wrong errs: %v", errs)
	}
	if !(res[0].(int) == 6 && res[1] == nil &&
		res[2].(int) == 10 && res[3] == nil && res[4].(int) == 14) {
		log.Fatalf("Receive wrong res: %v", res)
	}
}

func TestEveryNilWorkerPool(t *testing.T) {
	c := NewCombiner(nil)

	res, errs := c.Every([]e.BatchResult{})
	if errs != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, errs)
	}
	e1, _ := e.NewEngine(
		1, 10, time.Duration(2*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, nil)
	e2, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, nil)

	br1 := e1.Submit(3)
	br2 := e2.Submit(4)
	br3 := e1.Submit(5)
	br4 := e2.Submit(6)
	br5 := e1.Submit(7)

	res, errs = c.Every([]e.BatchResult{br1, br2, br3, br4, br5})
	if len(res) != 5 || len(errs) != 5 {
		log.Fatalf("errs should be nil, but instead we got %v", errs)
	}
	if !(errs[0] == nil && errs[1] == ErrTest &&
		errs[2] == nil && errs[3] == ErrTest && errs[4] == nil) {
		log.Fatalf("Receive wrong errs: %v", errs)
	}
	if !(res[0].(int) == 6 && res[1] == nil &&
		res[2].(int) == 10 && res[3] == nil && res[4].(int) == 14) {
		log.Fatalf("Receive wrong res: %v", res)
	}
}

func TestEveryWithCtx(t *testing.T) {
	c := NewCombiner(wp)

	res, errs := c.EveryWithContext(
		context.Background(), []e.BatchResult{})
	if errs != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, errs)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, errs = c.EveryWithContext(ctx, []e.BatchResult{})
	if errs == nil || len(errs) != 0 {
		log.Fatalf("errs have len zero, but instead we got %v", errs)
	}

	e1, _ := e.NewEngine(
		1, 10, time.Duration(2*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, wp)
	e2, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, wp)

	br1 := e1.Submit(3)
	br2 := e2.Submit(4)
	br3 := e1.Submit(5)
	br4 := e2.Submit(6)
	br5 := e1.Submit(7)

	res, errs = c.EveryWithContext(context.Background(), []e.BatchResult{br1, br2, br3, br4, br5})
	if len(res) != 5 || len(errs) != 5 {
		log.Fatalf("errs should be nil, but instead we got %v", errs)
	}
	if !(errs[0] == nil && errs[1] == ErrTest &&
		errs[2] == nil && errs[3] == ErrTest && errs[4] == nil) {
		log.Fatalf("Receive wrong errs: %v", errs)
	}
	if !(res[0].(int) == 6 && res[1] == nil &&
		res[2].(int) == 10 && res[3] == nil && res[4].(int) == 14) {
		log.Fatalf("Receive wrong res: %v", res)
	}
}

func TestEveryWithCtxNilWorkerPool(t *testing.T) {
	c := NewCombiner(nil)

	res, errs := c.EveryWithContext(
		context.Background(), []e.BatchResult{})
	if errs != nil || res != nil {
		log.Fatalf("Both should be nil, but instead we got %v and %v", res, errs)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, errs = c.EveryWithContext(ctx, []e.BatchResult{})
	if errs == nil || len(errs) != 0 {
		log.Fatalf("errs have len zero, but instead we got %v", errs)
	}

	e1, _ := e.NewEngine(
		1, 10, time.Duration(2*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			res := make(map[uint64]interface{})
			for k, v := range m {
				res[k] = v.(int) * 2
			}
			return res, nil
		}, nil)
	e2, _ := e.NewEngine(
		1, 10, time.Duration(1*time.Millisecond),
		func(m map[uint64]interface{}) (
			map[uint64]interface{}, error) {
			return nil, ErrTest
		}, nil)

	br1 := e1.Submit(3)
	br2 := e2.Submit(4)
	br3 := e1.Submit(5)
	br4 := e2.Submit(6)
	br5 := e1.Submit(7)

	res, errs = c.EveryWithContext(context.Background(), []e.BatchResult{br1, br2, br3, br4, br5})
	if len(res) != 5 || len(errs) != 5 {
		log.Fatalf("errs should be nil, but instead we got %v", errs)
	}
	if !(errs[0] == nil && errs[1] == ErrTest &&
		errs[2] == nil && errs[3] == ErrTest && errs[4] == nil) {
		log.Fatalf("Receive wrong errs: %v", errs)
	}
	if !(res[0].(int) == 6 && res[1] == nil &&
		res[2].(int) == 10 && res[3] == nil && res[4].(int) == 14) {
		log.Fatalf("Receive wrong res: %v", res)
	}
}
