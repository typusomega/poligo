package policy_test

import (
	"fmt"
	"reflect"
	"time"

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

func (test *PolicySuite) TestHandleSetsBasePolicy() {
	var expectedFunc policy.HandlePredicate = (func(error) bool { return false })

	plcy := policy.Handle(expectedFunc).Retry()

	assert.Equal(test.T(), reflect.ValueOf(expectedFunc), reflect.ValueOf(plcy.BasePolicy.ShouldHandle), "policy's ShouldHandle not set correctly")
}

func (test *PolicySuite) TestOnlyGivenErrorsAreHandledInHandleType() {
	fn := policy.HandleType(CustomError{}).
		Retry().BasePolicy.ShouldHandle

	assert.False(test.T(), fn(fmt.Errorf("")), "ShouldHandle returned true but wrong error type was given")
	assert.True(test.T(), fn(CustomError{}), "ShouldHandle returned false but correct error type was given")
}

func (test *PolicySuite) TestAllGivenErrorsAreHandledWithOrCascade() {
	fn := policy.HandleType(CustomError{}).
		Or(AnotherCustomError{}).
		Or(fmt.Errorf("")).
		Retry().BasePolicy.ShouldHandle

	assert.True(test.T(), fn(CustomError{}), "ShouldHandle returned false but correct error type was given")
	assert.True(test.T(), fn(AnotherCustomError{}), "ShouldHandle returned false but correct error type was given")
	assert.True(test.T(), fn(fmt.Errorf("test")), "ShouldHandle returned false but correct error type was given")
}

func (test *PolicySuite) TestHandleAllHandlesAllKindsOfErrors() {
	fn := policy.HandleAll().
		Retry().BasePolicy.ShouldHandle

	assert.True(test.T(), fn(CustomError{}), "ShouldHandle returned false but correct error type was given")
	assert.True(test.T(), fn(AnotherCustomError{}), "ShouldHandle returned false but correct error type was given")
	assert.True(test.T(), fn(fmt.Errorf("test")), "ShouldHandle returned false but correct error type was given")
}

func (test *PolicySuite) TestRetryWithDurationsSetsSleepProviderAccordingly() {
	expectedDurations := []time.Duration{time.Nanosecond, time.Nanosecond * 2, time.Nanosecond * 3}

	sleepProvider := policy.HandleAll().
		Retry(policy.WithDurations(expectedDurations...)).SleepDurationProvider

	for i, expectedDuration := range expectedDurations {
		duration, _ := sleepProvider(i)
		assert.Equal(test.T(), expectedDuration, duration, "sleepProvider duration does not match given duration")
	}

	dur, ok := sleepProvider(123)
	assert.Equal(test.T(), expectedDurations[len(expectedDurations)-1], dur, "sleepProvider duration does not match given duration")
	assert.False(test.T(), ok, "sleepProvider returned ok without having duration configured")
}

func (test *PolicySuite) TestWithSleepDurationProviderSetsCorrectProvider() {
	var expectedFunc policy.SleepDurationProvider = func(try int) (duration time.Duration, ok bool) { return time.Second, false }

	plcy := policy.HandleAll().Retry(policy.WithSleepDurationProvider(expectedFunc))

	assert.Equal(test.T(), reflect.ValueOf(expectedFunc), reflect.ValueOf(plcy.SleepDurationProvider), "policy's SleepDurationProvider not set correctly")
}

func (test *PolicySuite) TestWithRetriesSetsRetriesCorrectly() {
	expectedRetries := 4

	plcy := policy.HandleAll().Retry(policy.WithRetries(expectedRetries))

	assert.Equal(test.T(), expectedRetries, plcy.ExpectedRetries, "policy's ExpectedRetries not set correctly")
}

func (test *PolicySuite) TestWithCallbackSetsCallback() {
	var expectedFunc policy.OnRetryCallback = func(err error, retryCount int) {}

	plcy := policy.HandleAll().Retry(policy.WithCallback(expectedFunc))

	assert.Equal(test.T(), reflect.ValueOf(expectedFunc), reflect.ValueOf(plcy.Callback), "policy's Callback not set correctly")
}

func (test *PolicySuite) TestWithPredicatesSetsPredicates() {
	var pred1 policy.RetryPredicate = func(val interface{}) bool { return true }
	var pred2 policy.RetryPredicate = func(val interface{}) bool { return false }
	expectedPredicates := []policy.RetryPredicate{pred1, pred2}

	plcy := policy.HandleAll().Retry(policy.WithPredicates(expectedPredicates...))

	assert.Equal(test.T(), expectedPredicates, plcy.Predicates, "policy's Predicates not set correctly")
}
