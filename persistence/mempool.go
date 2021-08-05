package persistence

import (
	"fmt"
	"sync"
)

// A type alias that represents a type of pool state event
type PoolEvent string

// A set of constants that represent valids types of pool state events
const (
	// Pool is full
	POOLFULL PoolEvent = "pool is full"
	// Pool is empty
	POOLEMPTY PoolEvent = "pool is empty"
	// Pool has been reset
	POOLPURGE PoolEvent = "pool has been purged"
	// Pool has been resized
	POOLRESIZE PoolEvent = "pool has been resized"
)

// A struct that represents a memory pool of objects
type MemPool struct {
	// Represents the actual addressable pool of objects
	pool map[string]interface{}
	// Represents the syncrhonization lock for the pool
	mutex sync.Mutex

	// Represents the current number of object in the pool
	Count uint
	// Represents the maximum size of the pool
	size uint

	// Represents the event channel for pool state events
	eventchan chan PoolEvent
}

// A constructor function that creates a new memory pool of the given size and returns it.
// Requires a boolean that determines if the event handler for the pool should be initialized.
func NewMemPool(poolsize uint, enable_events bool) *MemPool {
	// Create a new MemPool object
	pool := &MemPool{
		pool:  make(map[string]interface{}),
		mutex: sync.Mutex{},
		Count: 0,
		size:  poolsize,
	}

	// Check if the event handler should be enabled
	if enable_events {
		// Initialize the event handler
		pool.eventchan = make(chan PoolEvent)
	}

	// Return the new pool
	return pool
}

// A method of MemPool that returns the event channel for the pool
func (pool *MemPool) GetEventChannel() chan PoolEvent {
	return pool.eventchan
}

// A method of MemPool that reports whether the pool is full
func (pool *MemPool) IsFull() bool {
	return pool.Count >= pool.size
}

// A method of MemPool that reports whether the pool is empty
func (pool *MemPool) IsEmpty() bool {
	return pool.Count == 0
}

// A method of MemPool that adds an object to the pool which is addressable by the given key.
// Returns an error if the pool is full. If an object exists for the key, it is overwritten.
func (pool *MemPool) Put(key string, object interface{}) error {
	// Acquire the lock on the pool
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	// Check if the pool is full
	if pool.IsFull() {
		// Return an error
		return fmt.Errorf("pool is full")
	}

	// Increment the pool count
	pool.Count++
	// Add the object to the pool
	pool.pool[key] = object

	// Check if the pool is full
	if pool.IsFull() {
		// Check if the event handler is initalized
		if pool.eventchan != nil {
			// Send a pool full event
			pool.eventchan <- POOLFULL
		}
	}

	// Return nil error
	return nil
}

// A method of MemPool that returns the object that is addressable by the given key.
// Returns the object and a boolean that indicates whether the object exists in the pool.
func (pool *MemPool) Get(key string) (interface{}, bool) {
	// Acquire the lock on the pool
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	// Retrieve the object from the pool
	object, ok := pool.pool[key]
	// Return the object and a boolean that indicates whether the object exists in the pool
	return object, ok
}

// A method of MemPool that removes the object that is addressable by the given key.
func (pool *MemPool) Remove(key string) {
	// Acquire the lock on the pool
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	// Remove the object from the pool
	delete(pool.pool, key)
	// Update the count of the pool
	pool.Count = uint(len(pool.pool))

	// Check if the pool is empty
	if pool.IsEmpty() {
		// Check if the event handler is initalized
		if pool.eventchan != nil {
			// Send a pool empty event
			pool.eventchan <- POOLEMPTY
		}
	}
}

// A method of MemPool that retrieves the object that is addressable by the given key and removes it from the pool.
// Returns the object and a boolean that indicates whether the object exists in the pool.
func (pool *MemPool) Pop(key string) (interface{}, bool) {
	// Retrieve the object from the pool and the boolean that indicates whether the object exists in the pool
	object, ok := pool.Get(key)
	// Check if the object exists in the pool
	if ok {
		// Remove the object from the pool
		pool.Remove(key)
	}

	// Return the object and a boolean that indicates whether the object exists in the pool
	return object, ok
}
