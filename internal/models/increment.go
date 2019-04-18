package models

import "sync"

const (
	// MaxValue for this compuster is shifted & inverted uint(0)
	MaxValue int32 = 1<<31 - 1
	// MinValue value for incrementor
	MinValue int32 = 0
	// IncremetStep default step for increment
	IncremetStep int32 = 1
)

// Increment base struct for incrementor model
type Increment struct {
	Obj
	Username string
	// Current incrementor value
	Number int32
	// Max value incrementing to
	MaxValue int32
	// Incrementing step
	Step int32
	// sync mutex
	mutex *sync.Mutex
}
