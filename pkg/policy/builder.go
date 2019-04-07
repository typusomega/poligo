package policy

import (
	"reflect"
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
	plcy := defaultRetryPolicy()
	plcy.policy = policy{shouldHandle: it.handlePredicate}

	for _, opt := range opts {
		opt(plcy)
	}

	return plcy
}
