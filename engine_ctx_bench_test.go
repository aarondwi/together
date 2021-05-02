package together

import (
	"context"
	"log"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func batchFuncCtx(m map[uint64]interface{}) (map[uint64]interface{}, error) {
	// simulate slow-enough network call
	time.Sleep(5 * time.Millisecond)
	results := make(map[uint64]interface{}, 5000)
	for k, _ := range m {
		// make static to reduce allocation
		// I am benchmarking this library, not business code
		results[k] = 128
	}
	return results, nil
}

func BenchmarkEngineWithCtx_SingleCore_Parallel256(b *testing.B) {
	runtime.GOMAXPROCS(1)
	e, err := NewEngine(5, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_SingleCore_Parallel1024(b *testing.B) {
	runtime.GOMAXPROCS(1)
	e, err := NewEngine(5, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_SingleCore_Parallel4096(b *testing.B) {
	runtime.GOMAXPROCS(1)
	e, err := NewEngine(50, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_SingleCore_Parallel8192(b *testing.B) {
	runtime.GOMAXPROCS(1)
	e, err := NewEngine(50, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(8192)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_SingleCore_Parallel16384(b *testing.B) {
	runtime.GOMAXPROCS(1)
	e, err := NewEngine(50, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(16384)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel256(b *testing.B) {
	e, err := NewEngine(5, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel1024(b *testing.B) {
	e, err := NewEngine(5, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel4096(b *testing.B) {
	e, err := NewEngine(50, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel8192(b *testing.B) {
	e, err := NewEngine(50, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(8192)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel16384(b *testing.B) {
	e, err := NewEngine(50, 5000, time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(16384)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.SubmitWithContext(context.Background(), i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}
