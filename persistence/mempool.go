package persistence

import (
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
