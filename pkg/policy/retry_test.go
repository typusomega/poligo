package policy_test

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/typusomega/poliGo/pkg/policy"
)

func (test *PolicySuite) TestExecuteCalled() {
	executeCalled := false

	policy.Handle(func(err error) bool { return false }).
		Retry().
		Execute(context.Background(), func() (interface{}, error) {
			executeCalled = true
			return nil, fmt.Errorf("fail")
		})

	assert.True(test.T(), executeCalled, "execute not called")
}

func (test *PolicySuite) TestHandledCalledOnError() {
	handleCalled := false

	policy.Handle(func(err error) bool {
		handleCalled = true
		return true
	}).Retry().Execute(context.Background(), func() (interface{}, error) { return nil, fmt.Errorf("fail") })

	assert.True(test.T(), handleCalled, "handle not called")
}

func (test *PolicySuite) TestRetryDefaultsToOneRetry() {
	callCount := 0

	policy.Handle(func(_ error) bool { return true }).
		Retry().
		Execute(context.Background(), func() (interface{}, error) {
			callCount++
			return nil, fmt.Errorf("fail")
		})

	assert.Equal(test.T(), 2, callCount, "execute not called twice")
}

func (test *PolicySuite) TestRetriesOnlyIfPredicatesAreMet() {
	callCount := 0
	builder := policy.Handle(func(_ error) bool { return true })

	builder.
		Retry(policy.WithPredicates(func(val interface{}) bool {
			return true
		})).
		Execute(context.Background(), func() (interface{}, error) {
			callCount++
			return nil, nil
		})

	assert.Equal(test.T(), 2, callCount, "does not retry even if predicates are met")

	callCount = 0
	builder.
		Retry(policy.WithPredicates(func(val interface{}) bool {
			return false
		})).
		Execute(context.Background(), func() (interface{}, error) {
			callCount++
			return nil, nil
		})

	assert.Equal(test.T(), 1, callCount, "execute not called twice")
}

func (test *PolicySuite) TestPredicatesReceiveCorrectInput() {
	expectedVal := "val"

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithPredicates(func(val interface{}) bool {
			assert.Equal(test.T(), expectedVal, val, "val does not match action's return value")
			return true
		})).
		Execute(context.Background(), func() (interface{}, error) { return expectedVal, nil })
}

func (test *PolicySuite) TestRetriesAsMuchAsConfigured() {
	expectedRetries := 5
	callCount := 0

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithRetries(expectedRetries)).
		Execute(context.Background(), func() (interface{}, error) {
			callCount++
			return nil, fmt.Errorf("fail")
		})

	assert.Equal(test.T(), expectedRetries+1, callCount, "execute not called as much as configured")
}

func (test *PolicySuite) TestCallbackIsExecutedOnEachRetry() {
	callbackCallCount := 0

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithCallback(func(err error, retryCount int) { callbackCallCount++ })).
		Execute(context.Background(), func() (interface{}, error) {
			return nil, fmt.Errorf("fail")
		})

	assert.Equal(test.T(), 1, callbackCallCount, "execute not called as much as configured")
}

func (test *PolicySuite) TestRetriesAreStoppedWhenContextCancelled() {
	expectedCalls := 3
	callCount := 0
	ctx, cancel := context.WithCancel(context.Background())

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithRetries(5)).
		Execute(ctx, func() (interface{}, error) {
			callCount++
			if callCount >= expectedCalls {
				cancel()
			}
			return nil, fmt.Errorf("fail")
		})

	assert.Equal(test.T(), expectedCalls, callCount, "context cancel did not stop retries")
}

func (test *PolicySuite) TestRetriesSleepForGivenDurations() {
	callCount := 0

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithDurations(time.Nanosecond, time.Nanosecond*2, time.Nanosecond*3)).
		Execute(context.Background(), func() (interface{}, error) {
			callCount++
			return nil, fmt.Errorf("fail")
		})

	assert.Equal(test.T(), 4, callCount, "was not retried as often as durations were given")
}

func (test *PolicySuite) TestSleepDurationProviderIsUsedOnEachRetry() {
	callCount := 0
	expectedCalls := 3

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithSleepDurationProvider(func(try int) (duration time.Duration, ok bool) {
			if callCount >= expectedCalls {
				return time.Nanosecond, false
			}
			return time.Nanosecond, true
		})).
		Execute(context.Background(), func() (interface{}, error) {
			callCount++
			return nil, fmt.Errorf("fail")
		})

	assert.Equal(test.T(), expectedCalls, callCount, "was not called like configured in sleepDurationProvider")
}

// VOID

func (test *PolicySuite) TestVoidExecuteCalled() {
	executeCalled := false

	policy.Handle(func(err error) bool { return false }).
		Retry().
		ExecuteVoid(context.Background(), func() error {
			executeCalled = true
			return fmt.Errorf("fail")
		})

	assert.True(test.T(), executeCalled, "execute not called")
}

func (test *PolicySuite) TestVoidHandledCalledOnError() {
	handleCalled := false

	policy.Handle(func(err error) bool {
		handleCalled = true
		return true
	}).Retry().ExecuteVoid(context.Background(), func() error { return fmt.Errorf("fail") })

	assert.True(test.T(), handleCalled, "handle not called")
}

func (test *PolicySuite) TestVoidRetryDefaultsToOneRetry() {
	callCount := 0

	policy.Handle(func(_ error) bool { return true }).
		Retry().
		ExecuteVoid(context.Background(), func() error {
			callCount++
			return fmt.Errorf("fail")
		})

	assert.Equal(test.T(), 2, callCount, "execute not called twice")
}

func (test *PolicySuite) TestVoidIgnoresPredicates() {
	callCount := 0
	builder := policy.Handle(func(_ error) bool { return true })

	builder.
		Retry(policy.WithPredicates(func(val interface{}) bool {
			callCount++
			return true
		})).
		ExecuteVoid(context.Background(), func() error {
			return nil
		})

	assert.Equal(test.T(), 0, callCount, "predicate called")
}

func (test *PolicySuite) TestVoidRetriesAsMuchAsConfigured() {
	expectedRetries := 5
	callCount := 0

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithRetries(expectedRetries)).
		ExecuteVoid(context.Background(), func() error {
			callCount++
			return fmt.Errorf("fail")
		})

	assert.Equal(test.T(), expectedRetries+1, callCount, "execute not called as much as configured")
}

func (test *PolicySuite) TestVoidCallbackIsExecutedOnEachRetry() {
	callbackCallCount := 0

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithCallback(func(err error, retryCount int) { callbackCallCount++ })).
		ExecuteVoid(context.Background(), func() error {
			return fmt.Errorf("fail")
		})

	assert.Equal(test.T(), 1, callbackCallCount, "execute not called as much as configured")
}

func (test *PolicySuite) TestVoidRetriesAreStoppedWhenContextCancelled() {
	expectedCalls := 3
	callCount := 0
	ctx, cancel := context.WithCancel(context.Background())

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithRetries(5)).
		ExecuteVoid(ctx, func() error {
			callCount++
			if callCount >= expectedCalls {
				cancel()
			}
			return fmt.Errorf("fail")
		})

	assert.Equal(test.T(), expectedCalls, callCount, "context cancel did not stop retries")
}

func (test *PolicySuite) TestVoidRetriesSleepForGivenDurations() {
	callCount := 0

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithDurations(time.Nanosecond, time.Nanosecond*2, time.Nanosecond*3)).
		ExecuteVoid(context.Background(), func() error {
			callCount++
			return fmt.Errorf("fail")
		})

	assert.Equal(test.T(), 4, callCount, "was not retried as often as durations were given")
}

func (test *PolicySuite) TestVoidSleepDurationProviderIsUsedOnEachRetry() {
	callCount := 0
	expectedCalls := 3

	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithSleepDurationProvider(func(try int) (duration time.Duration, ok bool) {
			if callCount >= expectedCalls {
				return time.Nanosecond, false
			}
			return time.Nanosecond, true
		})).
		ExecuteVoid(context.Background(), func() error {
			callCount++
			return fmt.Errorf("fail")
		})

	assert.Equal(test.T(), expectedCalls, callCount, "was not called like configured in sleepDurationProvider")
}
