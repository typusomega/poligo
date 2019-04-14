package policy

import (
	"context"
	"sync"
	"time"
)

// CircuitBreakerPolicy is a policy offering circuit breaker capabillities
type CircuitBreakerPolicy struct {
	BasePolicy

	MaxErrors int
	//TODO error rate
	BrokenForProvider SleepDurationProvider
	OnBreak           OnBreakCallback
	OnReset           func()

	consecutiveErrors int
	mux               sync.Mutex
	broken            bool
}

// ExecuteVoid calls the given action and applies the policy
func (it *CircuitBreakerPolicy) ExecuteVoid(ctx context.Context, action func() error) error {
	if it.broken {
		return CircuitBrokenError{}
	}

	err := action()
	if err == nil {
		it.mux.Lock()
		it.consecutiveErrors = 0
		it.mux.Unlock()
		return err
	}

	if !it.ShouldHandle(err) {
		return err
	}

	it.breakIfNecessary(err)

	return err
}

// Execute calls the given action and applies the policy
func (it *CircuitBreakerPolicy) Execute(ctx context.Context, action func() (interface{}, error)) (interface{}, error) {
	if it.broken {
		return nil, CircuitBrokenError{}
	}

	outcome, err := action()
	if err == nil {
		it.mux.Lock()
		it.consecutiveErrors = 0
		it.mux.Unlock()
		return outcome, err
	}

	if !it.ShouldHandle(err) {
		return outcome, err
	}

	it.breakIfNecessary(err)

	return outcome, err
}

func (it *CircuitBreakerPolicy) resetAfter(duration time.Duration) {
	time.Sleep(duration)

	it.mux.Lock()
	it.broken = false
	it.consecutiveErrors = 0
	it.mux.Unlock()

	it.OnReset()
}

func (it *CircuitBreakerPolicy) breakCircuit(err error) {
	it.broken = true

	dur, _ := it.BrokenForProvider(it.consecutiveErrors)
	it.OnBreak(err, dur)

	go it.resetAfter(dur)
}

func (it *CircuitBreakerPolicy) breakIfNecessary(err error) {
	it.mux.Lock()
	it.consecutiveErrors++
	if it.consecutiveErrors >= it.MaxErrors {
		it.breakCircuit(err)
	}
	it.mux.Unlock()
}

// CircuitBrokenError signalizes that the circuit is currently broken
type CircuitBrokenError struct {
}

func (CircuitBrokenError) Error() string {
	return "circuit broken"
}
