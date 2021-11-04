package workerpool

import (
	"math/rand"
	"sync"
)

// Please look at https://github.com/golang/go/issues/21393
// for why we need mutex
var globalRng = rand.New(rand.NewSource(rand.Int63()))
var mu = sync.Mutex{}

var randPool = sync.Pool{
	New: func() interface{} {
		mu.Lock()
		r := rand.New(
			rand.NewSource(globalRng.Int63()))
		mu.Unlock()
		return r
	},
}

func GetDefaultPartitioner(x int) func(arg interface{}) int {
	return func(arg interface{}) int {
		r := randPool.Get().(*rand.Rand)
		result := r.Intn(x)
		randPool.Put(r)
		return result
	}
}

func GetDefaultWorkerPool() *WorkerPool {
	var WP, _ = NewWorkerPool(8, 512, false)
	return WP
}
