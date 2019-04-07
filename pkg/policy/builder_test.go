package policy_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/typusomega/poliGo/pkg/policy"
)

type CustomError struct {
}

func (CustomError) Error() string {
	return ""
}

type AnotherCustomError struct {
}

func (AnotherCustomError) Error() string {
	return ""
}

func (test *PolicySuite) TestOnlyGivenErrorsAreHandledInHandleType() {
	callCount := 0
	policy.HandleType(CustomError{}).
		Retry().
		ExecuteVoid(context.Background(), func() error {
			callCount++
			return fmt.Errorf("")
		})
	assert.Equal(test.T(), 1, callCount, "retried but different error was thrown")

	callCount = 0
	policy.HandleType(CustomError{}).
		Retry().
		ExecuteVoid(context.Background(), func() error {
			callCount++
			return CustomError{}
		})
	assert.Equal(test.T(), 2, callCount, "retried but different error was thrown")
}

func (test *PolicySuite) TestAllGivenErrorsAreHandledWithOrCascade() {
	plcy := policy.HandleType(CustomError{}).
		Or(AnotherCustomError{}).
		Retry()

	callCount := 0
	plcy.ExecuteVoid(context.Background(), func() error {
		callCount++
		return CustomError{}
	})
	assert.Equal(test.T(), 2, callCount, "retried but different error was thrown")

	callCount = 0
	plcy.ExecuteVoid(context.Background(), func() error {
		callCount++
		return AnotherCustomError{}
	})
	assert.Equal(test.T(), 2, callCount, "retried but different error was thrown")
}

func (test *PolicySuite) TestHandleAllHandlesAllKindsOfErrors() {
	errs := []error{CustomError{}, AnotherCustomError{}, fmt.Errorf("fail"), errors.New("")}
	expectedRetries := len(errs) - 1
	callCount := 0

	policy.HandleAll().
		Retry(policy.WithRetries(expectedRetries)).
		ExecuteVoid(context.Background(), func() error {
			err := errs[callCount]
			callCount++
			return err
		})

	assert.Equal(test.T(), expectedRetries+1, callCount, "not all kinds of errors were handled")
}
