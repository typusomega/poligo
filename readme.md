# PoliGo

PoliGo is a Go resilience and fault-handling library to help developers express policies such as Retry in a fluent manner.

## Installation

__go.mod__
`require github.com/typusomega/poligo`

__go get__
`go get github.com/typusomega/poligo`

## How to use

```go
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
```

### Handle, HandleErrorType

Very often we want to have all kinds of errors handled no matter their reason.

```go
policy.HandleAll().
	Retry().
```

Sometimes checking errors is as trivial as just switching its type. 

```go
policy.HandleErrorType(MyCustomError{}).
	Or(AnotherCustomError{}).
```

But in some cases we have special needs and need special policies for specific states of a given error.
This is when `policy.Handle` comes into play.

```go
	// only handle errors with `lenghtNegative` with this policy
	pol := policy.Handle(func(err error) bool {
		if err != nil {
			if err, ok := err.(*areaError); ok {
				if err.lengthNegative() {
					return true
				}
			}
		}
		return false
	}).Retry()
```



PoliGo is strongly inspired by the awesome c# alternative [Polly](https://github.com/App-vNext/Polly)