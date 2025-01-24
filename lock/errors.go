package lock

import "errors"

var (
	ErrUnableToAcquireLock       = errors.New("unable to acquire lock")
	ErrUnableToReleaseLock       = errors.New("unable to release lock")
	ErrContextWithDeadlineNeeded = errors.New("context with deadline needed,Refer more https://tinyurl.com/58ccxyey")
)
