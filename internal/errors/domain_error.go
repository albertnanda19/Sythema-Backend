package errors

type DomainError struct {
	status  int
	message string
	err     error
}

func New(status int, message string) *DomainError {
	return &DomainError{status: status, message: message}
}

func Wrap(status int, message string, err error) *DomainError {
	return &DomainError{status: status, message: message, err: err}
}

func (e *DomainError) Error() string {
	if e == nil {
		return ""
	}
	if e.err != nil {
		return e.message
	}
	return e.message
}

func (e *DomainError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *DomainError) Status() int {
	if e == nil {
		return 500
	}
	if e.status == 0 {
		return 500
	}
	return e.status
}

func (e *DomainError) Message() string {
	if e == nil {
		return "Internal server error"
	}
	if e.message == "" {
		return "Internal server error"
	}
	return e.message
}
