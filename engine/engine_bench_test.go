package engine

import (
	"math/rand"
	"runtime"
	"testing"

	tp "github.com/aarondwi/together/testparam"
	WP "github.com/aarondwi/together/workerpool"
)

func BenchmarkEngine_Parallel256(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
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
			_, err := br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEngine_Parallel1024(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
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
			_, err := br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEngine_Parallel4096(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
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
			_, err := br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEngineSubmitMany(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(256 * 8) // A batch of 256 ints

	ch := make(chan []interface{}, 16)
	res := make([]interface{}, 0, 256)
	for i := 0; i < 256; i++ {
		res = append(res, i)
	}
	go func() {
		for {
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs := e.SubmitMany(r)
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEngineSubmitManyInto(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(256 * 8) // A batch of 256 ints

	ch := make(chan []interface{}, 16)
	res := make([]interface{}, 0, 256)
	for i := 0; i < 256; i++ {
		res = append(res, i)
	}
	go func() {
		for {
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		brs := make([]BatchResult, 0, 256)
		for pb.Next() {
			r := <-ch
			brs = brs[:0]
			e.SubmitManyInto(r, &brs)
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEngineSubmitMany_TwiceCoreNum(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 2)
	b.ReportAllocs()
	b.SetBytes(256 * 8) // A batch of 256 ints

	ch := make(chan []interface{}, 16)
	res := make([]interface{}, 0, 256)
	for i := 0; i < 256; i++ {
		res = append(res, i)
	}
	go func() {
		for {
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs := e.SubmitMany(r)
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEngineSubmitManyInto_TwiceCoreNum(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 2)
	b.ReportAllocs()
	b.SetBytes(256 * 8) // A batch of 256 ints

	ch := make(chan []interface{}, 16)
	res := make([]interface{}, 0, 256)
	for i := 0; i < 256; i++ {
		res = append(res, i)
	}
	go func() {
		for {
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		brs := make([]BatchResult, 0, 256)
		for pb.Next() {
			r := <-ch
			brs = brs[:0]
			e.SubmitManyInto(r, &brs)
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEngineSubmitMany_FourTimesCoreNum(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 4)
	b.ReportAllocs()
	b.SetBytes(256 * 8) // A batch of 256 ints

	ch := make(chan []interface{}, 16)
	res := make([]interface{}, 0, 256)
	for i := 0; i < 256; i++ {
		res = append(res, i)
	}
	go func() {
		for {
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs := e.SubmitMany(r)
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEngineSubmitManyInto_FourTimesCoreNum(b *testing.B) {
	var wp_ebt = WP.GetDefaultWorkerPool()
	e, err := NewEngine(
		EngineConfig{tp.NUM_OF_WORKER, tp.NUM_OF_ARGS_TO_WAIT, tp.SLEEP_DURATION},
		tp.BatchFunc, wp_ebt)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 4)
	b.ReportAllocs()
	b.SetBytes(256 * 8) // A batch of 256 ints

	ch := make(chan []interface{}, 16)
	res := make([]interface{}, 0, 256)
	for i := 0; i < 256; i++ {
		res = append(res, i)
	}
	go func() {
		for {
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		brs := make([]BatchResult, 0, 256)
		for pb.Next() {
			r := <-ch
			brs = brs[:0]
			e.SubmitManyInto(r, &brs)
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
