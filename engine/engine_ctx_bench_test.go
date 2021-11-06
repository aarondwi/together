package engine

import (
	"context"
	"math/rand"
	"testing"

	tp "github.com/aarondwi/together/testparam"
	WP "github.com/aarondwi/together/workerpool"
)

func BenchmarkEngineWithCtx_Parallel256(b *testing.B) {
	var wp_ecbt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ecbt)
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
	var wp_ecbt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ecbt)
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
	var wp_ecbt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ecbt)
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
