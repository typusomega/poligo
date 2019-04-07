package policy

import (
	"context"
	"time"
)

// Policy represents a specific policy
type Policy interface {
	// ExecuteVoid calls the given action and applies the policy
	ExecuteVoid(ctx context.Context, action func() error) error
	// Execute calls the given action and applies the policy
	Execute(ctx context.Context, action func() (interface{}, error)) (interface{}, error)
}

// BasePolicy is the base, all policy types have in common
type BasePolicy struct {
	ShouldHandle HandlePredicate
}

// SleepDurationProvider provides the next sleep duration for the given try
type SleepDurationProvider func(try int) (duration time.Duration, ok bool)
