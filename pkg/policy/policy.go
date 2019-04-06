package policy

import "context"

// Policy represents a specific policy
type Policy interface {
	// ExecuteVoid calls the given action and applies the policy
	ExecuteVoid(context context.Context, action func() error) error
	// Execute calls the given action and applies the policy
	Execute(context context.Context, action func() (interface{}, error)) (interface{}, error)
}

type policy struct {
	shouldHandle HandlePredicate
}
