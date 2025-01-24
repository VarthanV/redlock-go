package lock

import "errors"

var (
	ErrUnableToAcquireLock = errors.New("unable to acquire lock")
)
