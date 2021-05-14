package cluster

import "time"

var (
	NUM_OF_WORKER       = 128
	NUM_OF_ARGS_TO_WAIT = 1000
)

var (
	PARTITION_4  = 4
	PARTITION_8  = 8
	PARTITION_16 = 16
)

func BatchFunc(m map[uint64]interface{}) (map[uint64]interface{}, error) {
	// simulate fairly fast network call
	time.Sleep(2 * time.Millisecond)
	return m, nil
}
