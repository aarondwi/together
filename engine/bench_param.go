package engine

import (
	"math/rand"
	"time"

	WP "github.com/aarondwi/together/workerpool"
)

var (
	NUM_OF_WORKER       = 256
	NUM_OF_ARGS_TO_WAIT = 1024
)

var (
	PARTITION_4    = 4
	PARTITION_8    = 8
	PARTITION_16   = 16
	SLEEP_DURATION = time.Duration(5 * time.Millisecond)
)

func BatchFunc(m map[uint64]interface{}) (
	map[uint64]interface{}, error) {
	time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Microsecond)
	return m, nil
}

func GetEngineConfigForBenchmarks(n int) EngineConfig {
	return EngineConfig{
		NumOfWorker:  NUM_OF_WORKER / n,
		ArgSizeLimit: NUM_OF_ARGS_TO_WAIT / n,
		WaitDuration: SLEEP_DURATION,
	}
}

func GetEngineForBenchmarks(n int, wpb *WP.WorkerPool) *Engine {
	engine, _ := NewEngine(
		GetEngineConfigForBenchmarks(n),
		BatchFunc, wpb)
	return engine
}

func GetEnginesForBenchmarks(n int, wpb *WP.WorkerPool) []*Engine {
	engines := make([]*Engine, 0, n)
	for i := 0; i < n; i++ {
		engines = append(engines, GetEngineForBenchmarks(n, wpb))
	}
	return engines
}
