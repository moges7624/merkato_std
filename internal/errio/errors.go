package errio

// TODO: revisit this
type DomainError struct {
	Code string
	Err  error
}

func (e DomainError) Error() string {
	return e.Err.Error()
}

func (e DomainError) NotFound(err error) *DomainError {
	return &DomainError{
		Code: "not_found",
		Err:  err,
	}
}
