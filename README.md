# 🔒 Distributed Lock using Redis (redlock-go) 🔒

This project implements a distributed locking mechanism using Redis, inspired by the **RedLock algorithm**. The goal is to provide a reliable way to manage locks across multiple Redis instances, ensuring consistency and fault tolerance. 🚀

---

## ✨ Features

- **Distributed Locking**: Acquire and release locks across multiple Redis instances. 🔑
- **Quorum-based Decision Making**: Ensures that a lock is acquired or released only if a majority of Redis instances agree. ✅
- **Concurrency Control**: Uses a semaphore to limit the number of concurrent requests to Redis instances. 🚦
- **Context Support**: Supports context-based timeout and cancellation for lock operations. ⏳
- **Error Handling**: Provides clear error messages for common failure scenarios. ❌

---

## 🛠️ Installation

To use this package, you need to have **Go** installed on your machine. You can install the package using:

```bash
go get github.com/VarthanV/redlock-go
```
---

## 🚀 Usage

### Importing the Package

```go
import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/VarthanV/redlock-go/lock"
)
```

### Creating a New RedLock Instance

You need to provide a list of Redis clients and the duration for which the lock should be held.

```go
clients := []*redis.Client{
    redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
    redis.NewClient(&redis.Options{Addr: "localhost:6380"}),
    redis.NewClient(&redis.Options{Addr: "localhost:6381"}),
}

lockDuration := 10 * time.Second
l := lock.New(clients, lockDuration)
```

### Acquiring a Lock 🔐

To acquire a lock, use the `Acquire` method with a context and a unique key.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

key := "my-resource-key"
err := l.Acquire(ctx, key)
if err != nil {
    logrus.Errorf("Failed to acquire lock: %v", err)
    return
}

// 🚨 Critical section
// Perform operations that require the lock

// Release the lock when done
err = l.Release(ctx, key)
if err != nil {
    logrus.Errorf("Failed to release lock: %v", err)
}
```

### Releasing a Lock 🔓

To release a lock, use the `Release` method with the same context and key used to acquire the lock.

```go
err := l.Release(ctx, key)
if err != nil {
    logrus.Errorf("Failed to release lock: %v", err)
}
```

---

## ⚙️ Configuration

- **Max Concurrency**: The maximum number of concurrent requests allowed to Redis instances is set to **10** by default. You can adjust this by modifying the `maxConcurrencyAllowed` variable in the `New` function.
- **Lock Duration**: The duration for which the lock is held is specified when creating a new `redLock` instance.

---

## 🚨 Error Handling

The package defines several error scenarios:

- `ErrContextWithDeadlineNeeded`: The context provided must have a deadline. ⏰
- `ErrUnableToAcquireLock`: The lock could not be acquired on a majority of Redis instances. ❌
- `ErrUnableToReleaseLock`: The lock could not be released on a majority of Redis instances. ❌

---

## 📦 Dependencies

- [go-redis](https://github.com/redis/go-redis): Redis client for Go. 🛠️
- [logrus](https://github.com/sirupsen/logrus): Structured logger for Go. 📝

---


# 🎯 Use Cases for Distributed Lock using Redis (RedLock)

This package provides a reliable distributed locking mechanism using Redis. Here are some key use cases:

---

## 1. **Distributed Resource Management** 🌐  
   Ensure exclusive access to shared resources (e.g., databases, files, APIs) across multiple services or instances.

## 2. **Preventing Race Conditions** 🏁  
   Avoid conflicts when multiple processes or threads modify shared data simultaneously.

## 3. **Scheduled Task Coordination** ⏰  
   Ensure scheduled tasks (e.g., cron jobs) run only once across multiple instances.

## 4. **Leader Election** 👑  
   Elect a leader in a distributed system to handle specific tasks (e.g., background jobs).

## 5. **Inventory Management** 📦  
   Prevent overselling by locking inventory items during checkout in e-commerce systems.

## 6. **Rate Limiting** 🚦  
   Control the number of requests or operations across distributed services.

## 7. **Distributed Caching** 🗄️  
   Manage cache updates or invalidations to avoid stale data.

## 8. **Batch Processing** 🔄  
   Coordinate batch jobs to prevent duplicate processing.

## 9. **Critical Section Protection** 🔒  
   Protect code sections that should not run concurrently in distributed environments.

## 10. **Distributed Transactions** 💳  
   Coordinate multi-step workflows (e.g., payment processing) for consistency.

## 11. **Preventing Duplicate Events** 🔄  
   Ensure messages or events are processed only once in event-driven systems.

## 12. **Maintenance Operations** 🛠️  
   Coordinate tasks like backups or cleanup across distributed systems.

## 13. **Game Development** 🎮  
   Manage player actions or game state updates in multiplayer games.

## 14. **Financial Systems** 💰  
   Ensure atomicity in transactions (e.g., account transfers) to prevent double-spending.

## 15. **Distributed File Systems** 📂  
   Manage file locks to prevent conflicts in distributed file systems.

---


## 🙏 Acknowledgments

- Inspired by the [RedLock algorithm](https://redis.io/topics/distlock) by Redis. 🎯

---
