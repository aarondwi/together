package together

import "errors"

// ErrNumOfWorkerLessThanEqualZero is returned
// when given numOfWorker is <= zero
var ErrNumOfWorkerLessThanEqualZero = errors.New(
	"numOfWorker is expected to be > 0")

// ErrArgSizeLimitLessThanEqualOne is returned
// when given argSizeLimit <= 1
var ErrArgSizeLimitLessThanEqualOne = errors.New(
	"argSizeLimit is expected to be > 1")

// ErrResultNotFound is returned
// when an ID is not found in the results map
var ErrResultNotFound = errors.New(
	"A result is not found on the resulting map. " +
		"Please check your code to ensure all ids are matched with corresponding result.")
