package together

import (
	"log"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func batchFunc(m map[uint64]interface{}) (map[uint64]interface{}, error) {
	// simulate network call
	time.Sleep(2 * time.Millisecond)
	results := make(map[uint64]interface{}, 5000)
	for k, v := range m {
		results[k] = v.(int64) + 1
	}
	return results, nil
}

func BenchmarkEngine_SingleCore_Parallel256(b *testing.B) {
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
			_, err := e.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngine_SingleCore_Parallel1024(b *testing.B) {
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
			_, err := e.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkEngine_Parallel256(b *testing.B) {
	e, err := NewEngine(5, 5000, time.Duration(5)*time.Millisecond, batchFunc)
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
	e, err := NewEngine(5, 5000, time.Duration(5)*time.Millisecond, batchFunc)
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
