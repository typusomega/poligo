package policy

import (
	"time"
)

// DefaultRetries is the default number of retries if not overriden by options
const DefaultRetries = 1

func defaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		expectedRetries:       DefaultRetries,
		sleepDurationProvider: func(int) (time.Duration, bool) { return time.Nanosecond, false },
		predicates:            []RetryPredicate{},
		callback:              func(err error, retryCount int) {},
	}
}
