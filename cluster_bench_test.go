package together

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkCluster_Partition4_Parallel256(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_4, GetDefaultPartitioner(PARTITION_4),
		// per-engine param
		NUM_OF_WORKER/PARTITION_4, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition4_Parallel1024(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_4, GetDefaultPartitioner(PARTITION_4),
		// per-engine param
		NUM_OF_WORKER/PARTITION_4, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition4_Parallel4096(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_4, GetDefaultPartitioner(PARTITION_4),
		// per-engine param
		NUM_OF_WORKER/PARTITION_4, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition4_Parallel16384(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_4, GetDefaultPartitioner(PARTITION_4),
		// per-engine param
		NUM_OF_WORKER/PARTITION_4, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(16384)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition8_Parallel256(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_8, GetDefaultPartitioner(PARTITION_8),
		// per-engine param
		NUM_OF_WORKER/PARTITION_8, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition8_Parallel1024(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_8, GetDefaultPartitioner(PARTITION_8),
		// per-engine param
		NUM_OF_WORKER/PARTITION_8, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition8_Parallel4096(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_8, GetDefaultPartitioner(PARTITION_8),
		// per-engine param
		NUM_OF_WORKER/PARTITION_8, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition8_Parallel16384(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_8, GetDefaultPartitioner(PARTITION_8),
		// per-engine param
		NUM_OF_WORKER/PARTITION_8, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(16384)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition16_Parallel256(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_16, GetDefaultPartitioner(PARTITION_16),
		// per-engine param
		NUM_OF_WORKER/PARTITION_16, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(256)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition16_Parallel1024(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_16, GetDefaultPartitioner(PARTITION_16),
		// per-engine param
		NUM_OF_WORKER/PARTITION_16, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(1024)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition16_Parallel4096(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_16, GetDefaultPartitioner(PARTITION_16),
		// per-engine param
		NUM_OF_WORKER/PARTITION_16, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(4096)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func BenchmarkCluster_Partition16_Parallel16384(b *testing.B) {
	c, err := NewCluster(
		// cluster params
		PARTITION_16, GetDefaultPartitioner(PARTITION_16),
		// per-engine param
		NUM_OF_WORKER/PARTITION_16, NUM_OF_ARGS_TO_WAIT,
		time.Duration(5*time.Millisecond),
		batchFunc,
	)
	if err != nil {
		log.Fatal(err)
	}

	b.SetParallelism(16384)
	b.RunParallel(func(pb *testing.PB) {
		i := rand.Int63n(1000000)
		for pb.Next() {
			i++
			_, err := c.Submit(i)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}
