package interfaces

// Incrementor is the interface that wrap basic incrementor method
type Incrementor interface {
	// Must return current number
	GetNumber() int64
	// Increment current number by
	// After each call must return current number more on one
	IncrementNumber()
	// Lower boundry must not be less then zero
	// Must set upper boundry to incrementing number
	// Must set incrementor step to incrementing number
	// On call currentNumber must by reseted to lower boundary if  currentNumber greater then new maximum value
	SetSettings(step uint64, max uint64)
}
