package cluster

import (
	"log"
	"math/rand"
	"testing"

	c "github.com/aarondwi/together/common"
	e "github.com/aarondwi/together/engine"
)

func BenchmarkCluster_Partition4_Parallel256(b *testing.B) {
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_4, c.GetDefaultPartitioner(c.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_4,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_4, c.GetDefaultPartitioner(c.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_4,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_4, c.GetDefaultPartitioner(c.PARTITION_4),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_4,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_8, c.GetDefaultPartitioner(c.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_8,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_8, c.GetDefaultPartitioner(c.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_8,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_8, c.GetDefaultPartitioner(c.PARTITION_8),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_8,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_16, c.GetDefaultPartitioner(c.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_16,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_16, c.GetDefaultPartitioner(c.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_16,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
	var wp_cbt = c.GetDefaultWorkerPool()
	c, err := NewCluster(
		c.PARTITION_16, c.GetDefaultPartitioner(c.PARTITION_16),
		e.EngineConfig{
			NumOfWorker:  c.NUM_OF_WORKER / c.PARTITION_16,
			ArgSizeLimit: c.NUM_OF_ARGS_TO_WAIT,
			WaitDuration: c.SLEEP_DURATION},
		c.BatchFunc,
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
