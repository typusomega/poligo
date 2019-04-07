package policy_test

import (
	"context"
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

func (test *PolicySuite) TestOnlyGivenErrorsAreHandledInHandleErrorType() {
	callCount := 0
	policy.HandleErrorType(CustomError{}).
		Retry().
		ExecuteVoid(context.Background(), func() error {
			callCount++
			return fmt.Errorf("")
		})
	assert.Equal(test.T(), 1, callCount, "retried but different error was thrown")

	callCount = 0
	policy.HandleErrorType(CustomError{}).
		Retry().
		ExecuteVoid(context.Background(), func() error {
			callCount++
			return CustomError{}
		})
	assert.Equal(test.T(), 2, callCount, "retried but different error was thrown")
}

func (test *PolicySuite) TestAllGivenErrorsAreHandledWithOrCascade() {
	plcy := policy.HandleErrorType(CustomError{}).
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
