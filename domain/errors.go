package domain

type CauseCycleError struct {
}

func (e CauseCycleError) Error() string {
	return "cycle detected"
}
