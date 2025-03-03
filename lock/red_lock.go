package lock

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type redLock struct {
	clients      []*redis.Client
	requestSem   chan struct{}
	lockDuration time.Duration
}

func New(clients []*redis.Client, lockDuration time.Duration) ILock {
	maxConcurrencyAllowed := 10
	l := &redLock{
		clients:      clients,
		requestSem:   make(chan struct{}),
		lockDuration: lockDuration,
	}

	for i := 0; i < maxConcurrencyAllowed; i++ {
		l.requestSem <- struct{}{}
	}
	return l
}

type acquireActionOutcome struct {
	acquired bool
}

// Acquire implements ILock.
func (r *redLock) Acquire(ctx context.Context, key string) error {
	var (
		acquiredCount       = atomic.Int32{}
		quorum              = len(r.clients)/2 + 1
		acquireActionStream = make(chan acquireActionOutcome)
		done                = make(chan interface{})
		wg                  sync.WaitGroup
	)

	defer close(acquireActionStream)

	_, ok := ctx.Deadline()
	if !ok {
		logrus.Error(ErrContextWithDeadlineNeeded.Error())
		return ErrContextWithDeadlineNeeded
	}
	go func() {
		for val := range acquireActionStream {
			if val.acquired {
				newVal := acquiredCount.Add(1)
				if newVal == int32(quorum) {
					close(done)
					return
				}
			}
		}
	}()

	for _, c := range r.clients {
		wg.Add(1)
		// Practically finite number of redis clusters which is acceptable count to make concurrent requests
		// will only will be there, just adding semaphore to gracefully handle the situation
		go func(c *redis.Client) {
			defer wg.Done()
			<-r.requestSem
			// After the sem is acquired release it after the operation is done
			if acquiredCount.Load() == int32(quorum) {
				close(done)
				return
			}
			acquiredResult := c.SetNX(ctx, key, "1", r.lockDuration)
			acquireActionStream <- acquireActionOutcome{
				acquired: acquiredResult.Val(),
			}
			r.requestSem <- struct{}{}
		}(c)

	}

	go func() {
		wg.Wait()
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		default:
			close(done)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			if acquiredCount.Load() == int32(quorum) {
				return nil
			}
			return ErrUnableToAcquireLock
		}

	}
}

// Release implements ILock.
func (r *redLock) Release(ctx context.Context, key string) error {
	var (
		releasedCount = atomic.Int32{}
		quorum        = len(r.clients)/2 + 1
		done          = make(chan interface{})
		wg            sync.WaitGroup
	)

	_, ok := ctx.Deadline()
	if !ok {
		logrus.Error(ErrContextWithDeadlineNeeded.Error())
		return ErrContextWithDeadlineNeeded
	}

	for _, c := range r.clients {
		wg.Add(1)
		go func(client *redis.Client) {
			defer wg.Done()
			<-r.requestSem
			// After the sem is acquired release it after the operation is done
			if releasedCount.Load() == int32(quorum) {
				close(done)
				return
			}
			result := client.Del(ctx, key)
			if result.Val() == 1 {
				if releasedCount.Add(1) == int32(quorum) {
					close(done)
					return
				}
			}
			r.requestSem <- struct{}{}
		}(c)
	}

	go func() {
		wg.Wait()
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		default:
			close(done)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			if releasedCount.Load() == int32(quorum) {
				return nil
			}
			return ErrUnableToReleaseLock
		}
	}
}
