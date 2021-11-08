package combiner

import (
	"math/rand"
	"runtime"
	"testing"

	e "github.com/aarondwi/together/engine"
	WP "github.com/aarondwi/together/workerpool"
)

func BenchmarkAllSuccess_4Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllSuccess_4Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllSuccess_4Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllSuccess_8Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllSuccess_8Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllSuccess_8Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllSuccess_16Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllSuccess_16Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAllSuccess_16Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.AllSuccess(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_4Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_4Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_4Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_8Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_8Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_8Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_16Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_16Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkRace_16Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, err := c.Race(brs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkEvery_4Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEvery_4Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEvery_4Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_4, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_4 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_4)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_4; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEvery_8Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEvery_8Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEvery_8Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_8, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_8 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_8)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_8; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEvery_16Services_Parallel256(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(256 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEvery_16Services_Parallel1024(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(1024 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkEvery_16Services_Parallel4096(b *testing.B) {
	wpb := WP.GetDefaultWorkerPool()
	c := NewCombiner(wpb)
	engines := e.GetEnginesForBenchmarks(e.PARTITION_16, wpb)

	b.SetParallelism(4096 / runtime.NumCPU())
	b.ReportAllocs()
	b.SetBytes(int64(e.PARTITION_16 * 8))
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		brs := make([]e.BatchResult, 0, e.PARTITION_16)
		for pb.Next() {
			brs = brs[:0]
			for j := 0; j < e.PARTITION_16; j++ {
				br := engines[j].Submit(i)
				i++
				brs = append(brs, br)
			}
			_, errs := c.Every(brs)
			for _, err := range errs {
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
