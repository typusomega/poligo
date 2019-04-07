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
		handlePredicate: func(_ error) bool { return true },
	}
}

// Builder is used to build complex policies
type Builder interface {
	Retry(opts ...RetryOption) *RetryPolicy
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
