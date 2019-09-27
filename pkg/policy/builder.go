package policy

import (
	"reflect"
	"time"
)

// Handle is the entrypoint to build complex policies
// The given predicate is evaluated after each execution to check whether to apply the policy or not
func Handle(predicate HandlePredicate) Builder {
	return &builder{
		handlePredicate: predicate,
	}
}

// HandleType is the entrypoint to build policies handling errors of the type of the given errorObj
func HandleType(errorObj interface{}) ErrorBuilder {
	return &builder{
		handlePredicate: func(err error) bool { return reflect.TypeOf(err) == reflect.TypeOf(errorObj) },
	}
}

// HandleAll handles all kinds of errors
func HandleAll() Builder {
	return &builder{
		handlePredicate: func(err error) bool { return err != nil },
	}
}

// Builder is used to build complex policies
type Builder interface {
	Retry(opts ...RetryOption) *RetryPolicy
	WithCircuitBreaker(opts ...CircuitBreakerOption) *CircuitBreakerPolicy
}

// ErrorBuilder is used to build complex error policies
type ErrorBuilder interface {
	Builder
	Or(errorObj interface{}) ErrorBuilder
}

type builder struct {
	handlePredicate HandlePredicate
}

// HandlePredicate is used in the Handle function
type HandlePredicate func(err error) bool

// Or adds additional error types to the handled errors of the policy
func (it *builder) Or(errorObj interface{}) ErrorBuilder {
	pred := it.handlePredicate
	return &builder{
		handlePredicate: func(err error) bool {
			if reflect.TypeOf(err) == reflect.TypeOf(errorObj) {
				return true
			}
			return pred(err)
		},
	}
}

// Retry creates a RetryPolicy
func (it *builder) Retry(opts ...RetryOption) *RetryPolicy {
	plcy := DefaultRetryPolicy()
	plcy.BasePolicy = BasePolicy{ShouldHandle: it.handlePredicate}

	for _, opt := range opts {
		opt(plcy)
	}

	return plcy
}

// WithDurations sets durations for the diverse retries
func WithDurations(durations ...time.Duration) RetryOption {
	return func(o *RetryPolicy) {
		o.SleepDurationProvider = func(try int) (duration time.Duration, ok bool) {
			length := len(durations)
			if try < length {
				return durations[try], true
			}
			return durations[length-1], false
		}
	}
}

// WithSleepDurationProvider sets the SleepDurationProvider
func WithSleepDurationProvider(provider SleepDurationProvider) RetryOption {
	return func(o *RetryPolicy) {
		o.SleepDurationProvider = provider
	}
}

// WithRetries sets retries
func WithRetries(retries int) RetryOption {
	return func(o *RetryPolicy) {
		o.ExpectedRetries = retries
	}
}

// WithCallback sets on retry callback
func WithCallback(callback OnRetryCallback) RetryOption {
	return func(o *RetryPolicy) {
		o.Callback = callback
	}
}

// WithPredicates sets predicates checking for retry
func WithPredicates(predicates ...RetryPredicate) RetryOption {
	return func(o *RetryPolicy) {
		o.Predicates = predicates
	}
}

// WithCircuitBreaker creates a CircuitBreakerPolicy
func (it *builder) WithCircuitBreaker(opts ...CircuitBreakerOption) *CircuitBreakerPolicy {
	plcy := DefaultCircuitBreakerPolicy()
	plcy.BasePolicy = BasePolicy{ShouldHandle: it.handlePredicate}

	for _, opt := range opts {
		opt(plcy)
	}

	return plcy
}

// WithBrokenForProvider sets the SleepDurationProvider telling how long to keep the circuit broken for
func WithBrokenForProvider(provider SleepDurationProvider) CircuitBreakerOption {
	return func(o *CircuitBreakerPolicy) {
		o.BrokenForProvider = provider
	}
}

// WithOnBreakCallback sets the callback to be called whenever the circuit is broken
func WithOnBreakCallback(callback OnBreakCallback) CircuitBreakerOption {
	return func(o *CircuitBreakerPolicy) {
		o.OnBreak = callback
	}
}

// WithOnResetCallback sets the callback to be called whenever the circuit is reset
func WithOnResetCallback(callback func()) CircuitBreakerOption {
	return func(o *CircuitBreakerPolicy) {
		o.OnReset = callback
	}
}

// OnBreakCallback is the callback to be called whenever the circuit is broken
type OnBreakCallback func(error, time.Duration)

// CircuitBreakerOption modifies the CircuitBreakerPolicy
type CircuitBreakerOption func(*CircuitBreakerPolicy)
