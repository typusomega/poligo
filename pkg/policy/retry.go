package policy

import (
	"context"
	"time"
)

// RetryPolicy is a policy supporting retries
type RetryPolicy struct {
	policy

	expectedRetries       int
	sleepDurationProvider SleepDurationProvider
	callback              OnRetryCallback
	predicates            []RetryPredicate
}

// ExecuteVoid calls the given action and applies the policy
func (it *RetryPolicy) ExecuteVoid(ctx context.Context, action func() error) error {
	tryCount := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := action()
			if err == nil {
				return nil
			}

			if !it.shouldHandle(err) {
				return err
			}

			if !it.sleepIfRetryable(tryCount) {
				return err
			}

			it.callback(err, tryCount)
		}
		tryCount++
	}
}

// Execute calls the given action and applies the policy
func (it *RetryPolicy) Execute(ctx context.Context, action func() (interface{}, error)) (interface{}, error) {
	tryCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			val, err := action()

			if err == nil {
				for _, pred := range it.predicates {
					if pred(val) {
						continue
					} else {
						return val, nil
					}
				}
			}

			if !it.shouldHandle(err) {
				return val, err
			}

			if !it.sleepIfRetryable(tryCount) {
				return val, err
			}

			it.callback(err, tryCount)
		}
		tryCount++
	}
}

func (it *RetryPolicy) sleepIfRetryable(tryCount int) bool {
	sleepDuration, durationProvided := it.sleepDurationProvider(tryCount)
	canRetry := tryCount < it.expectedRetries || durationProvided
	if !canRetry {
		return false
	}

	time.Sleep(sleepDuration)

	return true
}

// OnRetryCallback is executed on every retry
type OnRetryCallback func(err error, retryCount int)

// RetryPredicate checks whether another retry is necessary
type RetryPredicate func(val interface{}) bool

// SleepDurationProvider provides the next sleep duration for the given try
type SleepDurationProvider func(try int) (duration time.Duration, ok bool)

// RetryOption options to enhance RetryPolicy
type RetryOption func(*RetryPolicy)

// WithDurations sets durations for the diverse retries
func WithDurations(durations ...time.Duration) RetryOption {
	return func(o *RetryPolicy) {
		o.sleepDurationProvider = func(try int) (duration time.Duration, ok bool) {
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
		o.sleepDurationProvider = provider
	}
}

// WithRetries sets retries
func WithRetries(retries int) RetryOption {
	return func(o *RetryPolicy) {
		o.expectedRetries = retries
	}
}

// WithCallback sets on retry callback
func WithCallback(callback OnRetryCallback) RetryOption {
	return func(o *RetryPolicy) {
		o.callback = callback
	}
}

// WithPredicates sets predicates checking for retry
func WithPredicates(predicates ...RetryPredicate) RetryOption {
	return func(o *RetryPolicy) {
		o.predicates = predicates
	}
}
