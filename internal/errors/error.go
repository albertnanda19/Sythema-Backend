package errors

import (
	"fmt"
)

type Error interface {
	error
	Code() string
	Message() string
	Status() int
	Unwrap() error
}

type domainError struct {
	code    string
	message string
	status  int
	err     error
}

func (e *domainError) Error() string {
	if e == nil {
		return ""
	}
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.code, e.message)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *domainError) Code() string {
	if e == nil {
		return CodeInternal
	}
	if e.code == "" {
		return CodeInternal
	}
	return e.code
}

func (e *domainError) Message() string {
	if e == nil {
		return MsgInternal
	}
	if e.message == "" {
		return MsgInternal
	}
	return e.message
}

func (e *domainError) Status() int {
	if e == nil {
		return 500
	}
	if e.status == 0 {
		return 500
	}
	return e.status
}

func (e *domainError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func New(code string, status int, message string) Error {
	return &domainError{code: code, status: status, message: message}
}

func Wrap(code string, status int, message string, err error) Error {
	return &domainError{code: code, status: status, message: message, err: err}
}
