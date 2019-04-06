package main

import (
	"context"
	"fmt"
	"time"

	"github.com/typusomega/poliGo/pkg/policy"
)

func main() {
	policy.Handle(func(_ error) bool { return true }).
		Retry(policy.WithDurations(time.Second, time.Second, time.Second)).
		Execute(context.Background(), func() (interface{}, error) {
			println("executed")
			return nil, fmt.Errorf("fail")
		})
}
