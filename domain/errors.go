package domain

type CauseCycleError struct {
}

func (e CauseCycleError) Error() string {
	return "cycle detected"
}

type RequestBodyError struct {
}

func (e RequestBodyError) Error() string {
	return "body attribute error"
}
