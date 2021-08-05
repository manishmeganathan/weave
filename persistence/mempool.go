package persistence

import "sync"

// A type alias that represents a type of pool state event
type PoolEvent uint

// A set of constants that represent valids types of pool state events
const (
	// Pool is full
	POOLFULL PoolEvent = 0
	// Pool is empty
	POOLEMPTY PoolEvent = 1
	// Pool has been reset
	POOLPURGE PoolEvent = 2
	// Pool has been resized
	POOLRESIZE PoolEvent = 3
)

// A struct that represents a memory pool of objects
type MemPool struct {
	// Represents the actual addressable pool of objects
	pool map[string]interface{}
	// Represents the syncrhonization lock for the pool
	mutex sync.Mutex

	// Represents the current size of the pool
	size uint
	// Represents the maximum size of the pool
	sizelimit uint

	// Represents the event handler chanel for pool state events
	eventchan chan PoolEvent
	// Represents whether the the event handler is enabled.
	events bool
}

// A constructor function that creates a new memory pool of the given size and returns it.
// Requires a boolean that determines if the event handler for the pool should be initialized.
func NewMemPool(poolsize uint, eventenable bool) *MemPool {
	// Create a new MemPool object
	pool := &MemPool{
		pool:      make(map[string]interface{}),
		mutex:     sync.Mutex{},
		size:      0,
		sizelimit: poolsize,
		events:    false,
	}

	// Check if the event handler should be enabled
	if eventenable {
		// Initialize the event handler
		pool.eventchan = make(chan PoolEvent)
		pool.events = true
	}

	// Return the new pool
	return pool
}
