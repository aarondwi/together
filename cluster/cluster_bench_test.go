package cluster

import (
	"log"
	"math/rand"
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
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
		log.Fatal(err)
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
				log.Fatal(err)
			}
			_, err = br.GetResult()
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}
