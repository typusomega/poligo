package policy

import (
	"time"
)

func defaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		expectedRetries:       DefaultRetries,
		sleepDurationProvider: func(int) (time.Duration, bool) { return time.Nanosecond, false },
		predicates:            []RetryPredicate{},
		callback:              func(err error, retryCount int) {},
	}
}
