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
	retry := policy.DefaultRetryPolicy()

	retry.Execute(context.Background(), func() (interface{}, error) {
		executeCalled = true
		return nil, fmt.Errorf("fail")
	})

	assert.True(test.T(), executeCalled, "execute not called")
}

func (test *PolicySuite) TestHandledCalledOnError() {
	handleCalled := false
	retry := policy.DefaultRetryPolicy()
	retry.BasePolicy.ShouldHandle = func(err error) bool {
		handleCalled = true
		return true
	}

	retry.Execute(context.Background(), func() (interface{}, error) { return nil, fmt.Errorf("fail") })

	assert.True(test.T(), handleCalled, "handle not called")
}

func (test *PolicySuite) TestRetriesOnlyIfHandleIsTrue() {
	callCount := 0
	retry := policy.DefaultRetryPolicy()
	retry.BasePolicy.ShouldHandle = func(err error) bool {
		return true
	}
	retry.Execute(context.Background(), func() (interface{}, error) {
		callCount++
		return nil, nil
	})
	assert.Equal(test.T(), 2, callCount, "does not retry even if predicates are met")

	callCount = 0
	retry.BasePolicy.ShouldHandle = func(err error) bool {
		return false
	}
	retry.Execute(context.Background(), func() (interface{}, error) {
		callCount++
		return nil, nil
	})
	assert.Equal(test.T(), 1, callCount, "execute not called twice")
}

func (test *PolicySuite) TestRetriesOnlyIfPredicatesAreMet() {
	callCount := 0
	retry := policy.DefaultRetryPolicy()
	retry.Predicates = []policy.RetryPredicate{
		func(val interface{}) bool {
			return true
		},
	}
	retry.Execute(context.Background(), func() (interface{}, error) {
		callCount++
		return nil, nil
	})
	assert.Equal(test.T(), 2, callCount, "does not retry even if predicates are met")

	callCount = 0
	retry.Predicates = []policy.RetryPredicate{
		func(val interface{}) bool {
			return false
		},
	}
	retry.Execute(context.Background(), func() (interface{}, error) {
		callCount++
		return nil, nil
	})
	assert.Equal(test.T(), 1, callCount, "execute not called twice")
}

func (test *PolicySuite) TestPredicatesReceiveCorrectInput() {
	expectedVal := "val"
	retry := policy.DefaultRetryPolicy()

	retry.Predicates = []policy.RetryPredicate{
		func(val interface{}) bool {
			assert.Equal(test.T(), expectedVal, val, "val does not match action's return value")
			return true
		},
	}
	retry.Execute(context.Background(), func() (interface{}, error) { return expectedVal, nil })
}

func (test *PolicySuite) TestRetriesAsMuchAsConfigured() {
	expectedRetries := 5
	callCount := 0
	retry := policy.DefaultRetryPolicy()

	retry.ExpectedRetries = expectedRetries
	retry.Execute(context.Background(), func() (interface{}, error) {
		callCount++
		return nil, fmt.Errorf("fail")
	})

	assert.Equal(test.T(), expectedRetries+1, callCount, "execute not called as much as configured")
}

func (test *PolicySuite) TestCallbackIsExecutedOnEachRetry() {
	callbackCallCount := 0
	retry := policy.DefaultRetryPolicy()

	retry.Callback = func(err error, retryCount int) { callbackCallCount++ }
	retry.Execute(context.Background(), func() (interface{}, error) {
		return nil, fmt.Errorf("fail")
	})

	assert.Equal(test.T(), 1, callbackCallCount, "execute not called as much as configured")
}

func (test *PolicySuite) TestRetriesAreStoppedWhenContextCancelled() {
	expectedCalls := 3
	callCount := 0
	ctx, cancel := context.WithCancel(context.Background())
	retry := policy.DefaultRetryPolicy()

	retry.ExpectedRetries = 5
	retry.Execute(ctx, func() (interface{}, error) {
		callCount++
		if callCount >= expectedCalls {
			cancel()
		}
		return nil, fmt.Errorf("fail")
	})

	assert.Equal(test.T(), expectedCalls, callCount, "context cancel did not stop retries")
}

func (test *PolicySuite) TestSleepDurationProviderIsUsedOnEachRetry() {
	callCount := 0
	expectedCalls := 3
	retry := policy.DefaultRetryPolicy()

	retry.SleepDurationProvider = func(try int) (duration time.Duration, ok bool) {
		if callCount >= expectedCalls {
			return time.Nanosecond, false
		}
		return time.Nanosecond, true
	}
	retry.Execute(context.Background(), func() (interface{}, error) {
		callCount++
		return nil, fmt.Errorf("fail")
	})

	assert.Equal(test.T(), expectedCalls, callCount, "was not called like configured in sleepDurationProvider")
}

// VOID

func (test *PolicySuite) TestVoidExecuteCalled() {
	executeCalled := false
	retry := policy.DefaultRetryPolicy()

	retry.ExecuteVoid(context.Background(), func() error {
		executeCalled = true
		return fmt.Errorf("fail")
	})

	assert.True(test.T(), executeCalled, "execute not called")
}

func (test *PolicySuite) TestVoidHandleCalledOnError() {
	handleCalled := false
	retry := policy.DefaultRetryPolicy()

	retry.BasePolicy.ShouldHandle = func(err error) bool {
		handleCalled = true
		return true
	}
	retry.ExecuteVoid(context.Background(), func() error { return fmt.Errorf("fail") })

	assert.True(test.T(), handleCalled, "handle not called")
}

func (test *PolicySuite) TestVoidRetriesOnlyIfHandleIsTrue() {
	callCount := 0
	retry := policy.DefaultRetryPolicy()
	retry.ExecuteVoid(context.Background(), func() error {
		callCount++
		return fmt.Errorf("")
	})
	assert.Equal(test.T(), 2, callCount, "does not retry even if handle is true")

	callCount = 0
	retry.BasePolicy.ShouldHandle = func(err error) bool {
		return false
	}
	retry.ExecuteVoid(context.Background(), func() error {
		callCount++
		return fmt.Errorf("")
	})
	assert.Equal(test.T(), 1, callCount, "execute not called twice")
}

func (test *PolicySuite) TestVoidIgnoresPredicates() {
	callCount := 0
	retry := policy.DefaultRetryPolicy()

	retry.Predicates = []policy.RetryPredicate{
		func(val interface{}) bool {
			callCount++
			return true
		},
	}
	retry.ExecuteVoid(context.Background(), func() error {
		return nil
	})

	assert.Equal(test.T(), 0, callCount, "predicate called")
}

func (test *PolicySuite) TestVoidRetriesAsMuchAsConfigured() {
	expectedRetries := 5
	callCount := 0
	retry := policy.DefaultRetryPolicy()

	retry.ExpectedRetries = expectedRetries
	retry.ExecuteVoid(context.Background(), func() error {
		callCount++
		return fmt.Errorf("fail")
	})

	assert.Equal(test.T(), expectedRetries+1, callCount, "execute not called as much as configured")
}

func (test *PolicySuite) TestVoidCallbackIsExecutedOnEachRetry() {
	callbackCallCount := 0
	retry := policy.DefaultRetryPolicy()

	retry.Callback = func(err error, retryCount int) { callbackCallCount++ }
	retry.ExecuteVoid(context.Background(), func() error {
		return fmt.Errorf("fail")

	})

	assert.Equal(test.T(), 1, callbackCallCount, "execute not called as much as configured")
}

func (test *PolicySuite) TestVoidRetriesAreStoppedWhenContextCancelled() {
	expectedCalls := 3
	callCount := 0
	ctx, cancel := context.WithCancel(context.Background())
	retry := policy.DefaultRetryPolicy()

	retry.ExpectedRetries = 5
	retry.ExecuteVoid(ctx, func() error {
		callCount++
		if callCount >= expectedCalls {
			cancel()
		}
		return fmt.Errorf("fail")
	})

	assert.Equal(test.T(), expectedCalls, callCount, "context cancel did not stop retries")
}

func (test *PolicySuite) TestVoidSleepDurationProviderIsUsedOnEachRetry() {
	callCount := 0
	expectedCalls := 3
	retry := policy.DefaultRetryPolicy()

	retry.SleepDurationProvider = func(try int) (duration time.Duration, ok bool) {
		if callCount >= expectedCalls {
			return time.Nanosecond, false
		}
		return time.Nanosecond, true
	}

	retry.ExecuteVoid(context.Background(), func() error {
		callCount++
		return fmt.Errorf("fail")
	})

	assert.Equal(test.T(), expectedCalls, callCount, "was not called like configured in sleepDurationProvider")
}
