package initializers

import "sync"

// Vars for initializers
var (
	initOnce  sync.Once
	initError error
	mutex     sync.RWMutex
)
