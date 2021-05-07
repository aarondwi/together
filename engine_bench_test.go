package together

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkEngine_Parallel256(b *testing.B) {
	e, err := NewEngine(
		NUM_OF_WORKER, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngine_Parallel1024(b *testing.B) {
	e, err := NewEngine(
		NUM_OF_WORKER, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngine_Parallel4096(b *testing.B) {
	e, err := NewEngine(
		NUM_OF_WORKER, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngine_Parallel16384(b *testing.B) {
	e, err := NewEngine(
		NUM_OF_WORKER, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5)*time.Millisecond, batchFunc)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(16384)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := e.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}
