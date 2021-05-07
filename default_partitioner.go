package together

import (
	"math/rand"
	"sync"
)

var globalRng = rand.New(rand.NewSource(rand.Int63()))

var randPool = sync.Pool{
	New: func() interface{} {
		return rand.New(
			rand.NewSource(globalRng.Int63() - int64(globalRng.Int31())))
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
