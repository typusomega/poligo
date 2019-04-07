# PoliGo

PoliGo is a Go resilience and fault-handling library to help developers express policies such as Retry in a fluent manner.

## How to use

```go
import (
	"context"
	"fmt"
	"time"

	"github.com/typusomega/poliGo/pkg/policy"
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




PoliGo is strongly inspired by the awesome c# alternative [Polly](https://github.com/App-vNext/Polly)