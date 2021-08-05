package persistence

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
