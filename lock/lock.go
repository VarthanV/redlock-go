package lock

import "context"

type ILock interface {
	Acquire(ctx context.Context, key string) error
	Release(ctx context.Context, key string) error
}
