package engine

import (
	"context"
	"math/rand"
	"testing"
)

func BenchmarkEngineWithCtx_Parallel256(b *testing.B) {
	once.Do(initEngineWP)
	e, err := NewEngine(
		EngineConfig{NUM_OF_WORKER, NUM_OF_ARGS_TO_WAIT, NUM_OF_ARGS_TO_WAIT * 2, SLEEP_DURATION},
		BatchFunc, wpb)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(256)
	b.ReportAllocs()
	b.SetBytes(8)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			br := e.Submit(i)
			_, err := br.GetResultWithContext(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel1024(b *testing.B) {
	once.Do(initEngineWP)
	e, err := NewEngine(
		EngineConfig{NUM_OF_WORKER, NUM_OF_ARGS_TO_WAIT, NUM_OF_ARGS_TO_WAIT * 2, SLEEP_DURATION},
		BatchFunc, wpb)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(1024)
	b.ReportAllocs()
	b.SetBytes(8)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			br := e.Submit(i)
			_, err := br.GetResultWithContext(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel4096(b *testing.B) {
	once.Do(initEngineWP)
	e, err := NewEngine(
		EngineConfig{NUM_OF_WORKER, NUM_OF_ARGS_TO_WAIT, NUM_OF_ARGS_TO_WAIT * 2, SLEEP_DURATION},
		BatchFunc, wpb)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(4096)
	b.ReportAllocs()
	b.SetBytes(8)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			br := e.Submit(i)
			_, err := br.GetResultWithContext(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
