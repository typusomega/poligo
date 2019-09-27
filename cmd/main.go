package main

import (
	"context"
	"fmt"
	"time"

	"github.com/typusomega/poligo/pkg/policy"
)

func main() {
	// select kind of errors to handle
	result, err := policy.Handle(func(_ error) bool { return true }).
		// tell you want to retry and how often
		Retry(policy.WithDurations(time.Second, time.Second, time.Second),
			// tell what to do before the next retry
			policy.WithCallback(log)).
		// execute the given action with the created policy
		Execute(context.Background(), doAwesomeStuff)

	fmt.Printf("executed with policy result: '%v', err: '%v'\n", result, err)
}

func log(err error, retryCount int) {
	fmt.Printf("execution failed: '%v' (retry %v)\n", err, retryCount)
}

func doAwesomeStuff() (interface{}, error) {
	return 42, fmt.Errorf("fail")
}
