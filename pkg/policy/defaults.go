package policy

import (
	"sync"
	"time"
)

// DefaultRetries is the default number of retries if not overriden by options
const DefaultRetries = 1

// DefaultBasePolicy is the base all policies come by default with
func DefaultBasePolicy() *BasePolicy {
	return &BasePolicy{ShouldHandle: func(_ error) bool { return true }}
}

// DefaultRetryPolicy is the default RetryPolicy
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		BasePolicy:            *DefaultBasePolicy(),
		ExpectedRetries:       DefaultRetries,
		SleepDurationProvider: func(int) (time.Duration, bool) { return time.Nanosecond, false },
		Predicates:            []RetryPredicate{},
		Callback:              func(err error, retryCount int) {},
	}
}

// DefaultCircuitBreakerPolicy is the default DefaultCircuitBreakerPolicy
func DefaultCircuitBreakerPolicy() *CircuitBreakerPolicy {
	return &CircuitBreakerPolicy{
		BasePolicy:        *DefaultBasePolicy(),
		MaxErrors:         DefaultRetries,
		BrokenForProvider: func(try int) (duration time.Duration, ok bool) { return time.Second * 2, true },
		OnBreak:           func(error, time.Duration) {},
		OnReset:           func() {},
		broken:            false,
		consecutiveErrors: 0,
		mux:               sync.Mutex{},
	}
}
