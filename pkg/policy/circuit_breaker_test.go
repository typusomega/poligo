package policy_test

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/typusomega/poliGo/pkg/policy"
)

func (test *PolicySuite) TestCBExecuteCalled() {
	executeCalled := false

	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.Execute(context.Background(), func() (interface{}, error) {
		executeCalled = true
		return nil, fmt.Errorf("fail")
	})

	assert.True(test.T(), executeCalled, "execute not called")
}

func (test *PolicySuite) TestCBHandleCalledOnError() {
	handleCalled := false

	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.BasePolicy.ShouldHandle = func(err error) bool {
		handleCalled = true
		return true
	}

	circuitBreaker.Execute(context.Background(), defaultFailingAction)

	assert.True(test.T(), handleCalled, "handle not called")
}

func (test *PolicySuite) TestCBOnlyAppliedIfShouldHandle() {
	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.MaxErrors = 0
	circuitBreaker.BasePolicy.ShouldHandle = func(err error) bool {
		return false
	}

	expErr := fmt.Errorf("fail")
	_, err := circuitBreaker.Execute(context.Background(), func() (interface{}, error) { return nil, expErr })
	assert.Equal(test.T(), expErr, err, "execute not called as often as expected")

	circuitBreaker.BasePolicy.ShouldHandle = func(err error) bool {
		return true
	}
	_, _ = circuitBreaker.Execute(context.Background(), defaultFailingAction)
	_, err = circuitBreaker.Execute(context.Background(), defaultFailingAction)
	assert.IsType(test.T(), policy.CircuitBrokenError{}, err)

}

func (test *PolicySuite) TestExecuteNotCalledWhenMoreThanMaxErrors() {
	executeCalled := 0

	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.MaxErrors = 1

	execute := func() (interface{}, error) {
		executeCalled++
		return nil, fmt.Errorf("fail")
	}

	circuitBreaker.Execute(context.Background(), execute)
	circuitBreaker.Execute(context.Background(), execute)
	circuitBreaker.Execute(context.Background(), execute)

	assert.Equal(test.T(), 1, executeCalled, "execute not called as often as expected")
}

func (test *PolicySuite) TestReturnsIfThereIsNoError() {
	expectedVal := "test"

	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	execute := func() (interface{}, error) {
		return expectedVal, nil
	}

	val, _ := circuitBreaker.Execute(context.Background(), execute)

	assert.Equal(test.T(), expectedVal, val, "execute not called as often as expected")
}

// Void

func (test *PolicySuite) TestVoidCBExecuteCalled() {
	executeCalled := false

	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.ExecuteVoid(context.Background(), func() error {
		executeCalled = true
		return fmt.Errorf("fail")
	})

	assert.True(test.T(), executeCalled, "execute not called")
}

func (test *PolicySuite) TestVoidCBHandleCalledOnError() {
	handleCalled := false

	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.BasePolicy.ShouldHandle = func(err error) bool {
		handleCalled = true
		return true
	}

	circuitBreaker.ExecuteVoid(context.Background(), defaultFailingVoidAction)

	assert.True(test.T(), handleCalled, "handle not called")
}

func (test *PolicySuite) TestVoidCBOnlyAppliedIfShouldHandle() {
	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.MaxErrors = 0
	circuitBreaker.BasePolicy.ShouldHandle = func(err error) bool {
		return false
	}

	expErr := fmt.Errorf("fail")
	err := circuitBreaker.ExecuteVoid(context.Background(), func() error { return expErr })
	assert.Equal(test.T(), expErr, err, "execute not called as often as expected")

	circuitBreaker.BasePolicy.ShouldHandle = func(err error) bool {
		return true
	}
	_ = circuitBreaker.ExecuteVoid(context.Background(), defaultFailingVoidAction)
	err = circuitBreaker.ExecuteVoid(context.Background(), defaultFailingVoidAction)
	assert.IsType(test.T(), policy.CircuitBrokenError{}, err)

}

func (test *PolicySuite) TestVoidExecuteNotCalledWhenMoreThanMaxErrors() {
	executeCalled := 0

	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.MaxErrors = 1

	execute := func() error {
		executeCalled++
		return fmt.Errorf("fail")
	}

	circuitBreaker.ExecuteVoid(context.Background(), execute)
	circuitBreaker.ExecuteVoid(context.Background(), execute)
	circuitBreaker.ExecuteVoid(context.Background(), execute)

	assert.Equal(test.T(), 1, executeCalled, "execute not called as often as expected")
}

func (test *PolicySuite) TestVoidReturnsIfThereIsNoError() {
	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	execute := func() error { return nil }

	err := circuitBreaker.ExecuteVoid(context.Background(), execute)

	assert.Nil(test.T(), err, "execute not called as often as expected")
}

// common

func (test *PolicySuite) TestCircuitBrokenError() {
	err := policy.CircuitBrokenError{}

	assert.Equal(test.T(), "circuit broken", err.Error())
}

func (test *PolicySuite) TestCircuitOpenedAfterBrokenForDuration() {
	expectedErr := fmt.Errorf("fail")
	action := func() (interface{}, error) { return nil, expectedErr }
	brokenTime := time.Millisecond * 5
	circuitBreaker := policy.DefaultCircuitBreakerPolicy()
	circuitBreaker.MaxErrors = 0
	circuitBreaker.BrokenForProvider = func(try int) (duration time.Duration, ok bool) {
		return brokenTime, true
	}

	_, _ = circuitBreaker.Execute(context.Background(), action)
	_, err := circuitBreaker.Execute(context.Background(), action)
	assert.IsType(test.T(), policy.CircuitBrokenError{}, err)

	sleepTime := brokenTime + time.Millisecond
	time.Sleep(sleepTime)
	_, err = circuitBreaker.Execute(context.Background(), action)
	assert.Equal(test.T(), expectedErr, err)
}

var defaultFailingAction = func() (interface{}, error) { return nil, fmt.Errorf("fail") }
var defaultFailingVoidAction = func() error { return fmt.Errorf("fail") }
