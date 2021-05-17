package cluster

import (
	"log"
	"math/rand"
	"testing"

	c "github.com/aarondwi/together/common"
)

var wp_cbt = c.GetDefaultWorkerPool()

func BenchmarkCluster_Partition4_Parallel256(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		c.PARTITION_4, c.GetDefaultPartitioner(c.PARTITION_4),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_4, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
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
	c, err := NewCluster(
		// cluster params
		c.PARTITION_4, c.GetDefaultPartitioner(c.PARTITION_4),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_4, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
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
	c, err := NewCluster(
		// cluster params
		c.PARTITION_4, c.GetDefaultPartitioner(c.PARTITION_4),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_4, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
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
	c, err := NewCluster(
		// cluster params
		c.PARTITION_8, c.GetDefaultPartitioner(c.PARTITION_8),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_8, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
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
	c, err := NewCluster(
		// cluster params
		c.PARTITION_8, c.GetDefaultPartitioner(c.PARTITION_8),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_8, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
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
	c, err := NewCluster(
		// cluster params
		c.PARTITION_8, c.GetDefaultPartitioner(c.PARTITION_8),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_8, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
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
	c, err := NewCluster(
		// cluster params
		c.PARTITION_16, c.GetDefaultPartitioner(c.PARTITION_16),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_16, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
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
	c, err := NewCluster(
		// cluster params
		c.PARTITION_16, c.GetDefaultPartitioner(c.PARTITION_16),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_16, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
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
	c, err := NewCluster(
		// cluster params
		c.PARTITION_16, c.GetDefaultPartitioner(c.PARTITION_16),
		// per-engine param
		c.NUM_OF_WORKER/c.PARTITION_16, c.NUM_OF_ARGS_TO_WAIT,
		c.SLEEP_DURATION,
		c.BatchFunc,
		wp_cbt,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
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
