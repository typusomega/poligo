package policy

import (
	"context"
	"time"
)

// RetryPolicy is a policy supporting retries
type RetryPolicy struct {
	BasePolicy

	ExpectedRetries       int
	SleepDurationProvider SleepDurationProvider
	Callback              OnRetryCallback
	Predicates            []RetryPredicate
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

			if !it.ShouldHandle(err) {
				return err
			}

			if !it.sleepIfRetryable(tryCount) {
				return err
			}

			it.Callback(err, tryCount)
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
				for _, pred := range it.Predicates {
					if !pred(val) {
						return val, nil
					}
				}
			}

			if !it.ShouldHandle(err) {
				return val, err
			}

			if !it.sleepIfRetryable(tryCount) {
				return val, err
			}

			it.Callback(err, tryCount)
		}
		tryCount++
	}
}

func (it *RetryPolicy) sleepIfRetryable(tryCount int) bool {
	sleepDuration, durationProvided := it.SleepDurationProvider(tryCount)
	canRetry := tryCount < it.ExpectedRetries || durationProvided
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

// RetryOption modifies the RetryPolicy
type RetryOption func(*RetryPolicy)
