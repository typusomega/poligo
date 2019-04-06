package policy

// Handle is the entrypoint to build policies
func Handle(predicate HandlePredicate) *Builder {
	return &Builder{
		handlePredicate: predicate,
	}
}

// Builder is used to build complex policies
type Builder struct {
	handlePredicate HandlePredicate
}

// HandlePredicate is used in the Handle function
type HandlePredicate func(err error) bool

// Retry creates a RetryPolicy
func (it *Builder) Retry(opts ...RetryOption) *RetryPolicy {
	policy := newRetryPolicy(policy{shouldHandle: it.handlePredicate})

	for _, opt := range opts {
		opt(policy)
	}

	return policy
}
