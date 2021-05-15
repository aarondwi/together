package engine

import (
	"context"
	"log"
	"math/rand"
	"testing"

	c "github.com/aarondwi/together/common"
)

func BenchmarkEngineWithCtx_Parallel256(b *testing.B) {
	e, err := NewEngine(
		c.NUM_OF_WORKER, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION, c.BatchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			br := e.Submit(i)
			_, err := br.GetResultWithContext(context.Background())
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel1024(b *testing.B) {
	e, err := NewEngine(
		c.NUM_OF_WORKER, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION, c.BatchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			br := e.Submit(i)
			_, err := br.GetResultWithContext(context.Background())
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel4096(b *testing.B) {
	e, err := NewEngine(
		c.NUM_OF_WORKER, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION, c.BatchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			br := e.Submit(i)
			_, err := br.GetResultWithContext(context.Background())
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineWithCtx_Parallel16384(b *testing.B) {
	e, err := NewEngine(
		c.NUM_OF_WORKER, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION, c.BatchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(16384)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			br := e.Submit(i)
			_, err := br.GetResultWithContext(context.Background())
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}
