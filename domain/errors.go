package domain

type ErrGraphCycle struct{}

func (e ErrGraphCycle) Error() string {
	return "cycle detected"
}

type ErrRequestBody struct{}

func (e ErrRequestBody) Error() string {
	return "body attribute error"
}

type ErrNotImplemented struct{}

func (e ErrNotImplemented) Error() string {
	return "not implemented"
}
