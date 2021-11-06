package cluster

import (
	"math/rand"
	"runtime"
	"testing"

	e "github.com/aarondwi/together/engine"
	tp "github.com/aarondwi/together/testparam"
	WP "github.com/aarondwi/together/workerpool"
)

func BenchmarkCluster_Partition4_Parallel256(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_4, WP.GetDefaultPartitioner(tp.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_4,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition4_Parallel1024(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_4, WP.GetDefaultPartitioner(tp.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_4,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition4_Parallel4096(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_4, WP.GetDefaultPartitioner(tp.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_4,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition8_Parallel256(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_8, WP.GetDefaultPartitioner(tp.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_8,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition8_Parallel1024(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_8, WP.GetDefaultPartitioner(tp.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_8,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition8_Parallel4096(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_8, WP.GetDefaultPartitioner(tp.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_8,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition16_Parallel256(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_16, WP.GetDefaultPartitioner(tp.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_16,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition16_Parallel1024(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_16, WP.GetDefaultPartitioner(tp.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_16,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition16_Parallel4096(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_16, WP.GetDefaultPartitioner(tp.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_16,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
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
			br, err := c.Submit(i)
			if err != nil {
				b.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition4(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_4, WP.GetDefaultPartitioner(tp.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_4,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition4_TwiceCoreNum(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_4, WP.GetDefaultPartitioner(tp.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_4,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 2)
	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition4_FourTimesCoreNum(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_4, WP.GetDefaultPartitioner(tp.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_4,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 4)
	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition8(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_8, WP.GetDefaultPartitioner(tp.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_8,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition8_TwiceCoreNum(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_8, WP.GetDefaultPartitioner(tp.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_8,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 2)
	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition8_FourTimesCoreNum(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_8, WP.GetDefaultPartitioner(tp.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_8,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 4)
	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition16(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_16, WP.GetDefaultPartitioner(tp.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_16,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition16_TwiceCoreNum(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_16, WP.GetDefaultPartitioner(tp.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_16,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 2)
	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkCluster_SubmitMany_Partition16_FourTimesCoreNum(b *testing.B) {
	var wp_cbt = WP.GetDefaultWorkerPool()
	c, err := NewCluster(
		tp.PARTITION_16, WP.GetDefaultPartitioner(tp.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  tp.NUM_OF_WORKER / tp.PARTITION_16,
			ArgSizeLimit: tp.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: tp.SLEEP_DURATION},
		tp.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		b.Fatal(err)
	}

	b.SetParallelism(runtime.NumCPU() * 4)
	b.ReportAllocs()
	b.SetBytes(256 * 8)

	ch := make(chan []interface{}, 256)
	go func() {
		for {
			res := make([]interface{}, 0, 256)
			for i := 0; i < 256; i++ {
				res = append(res, i)
			}
			ch <- res
		}
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := <-ch
			brs, err := c.SubmitMany(r)
			if err != nil {
				b.Fatal(err)
			}
			for _, br := range brs {
				_, err := br.GetResult()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
