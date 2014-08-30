package command

import (
	"errors"
)

type errorTracker struct {
	throw bool
	fifo  bool
	msg   string
}

// Create an error instance from the logs.
// If no error has been logged yet, nil is returned.
func (e *errorTracker) Error() error {
	if len(e.msg) > 0 {
		return errors.New(e.msg)
	}

	return nil
}

// Check if the tracker instance can hold any more errors.
// If the tracker is configured to drop old errors in favor
// of new ones, this method always returns false.
func (e *errorTracker) Saturated() bool {
	// FIFO is saturated on first entry
	// LIFO is never saturated
	return e.fifo && len(e.msg) > 0
}

// Store the given error message is not empty. If the tracker
// is configured to panic upon receiving errors, this method
// panics if the provided error is not an empty string.
func (e *errorTracker) Store(err string) {
	if len(err) > 0 {
		if false == e.Saturated() {
			e.msg = err
		}

		if e.throw {
			panic(err)
		}
	}
}

// Factory method for error trackers.
func newErrorTracker(panicOnError bool, storeFirstError bool) *errorTracker {
	return &errorTracker{panicOnError, storeFirstError, ""}
}
