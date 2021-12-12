package errors

import "fmt"

type NotFoundError struct {
	err string
}

func NotFound(err string) *NotFoundError {
	return &NotFoundError{err}
}

func NotFoundf(format string, v ...interface{}) *NotFoundError {
	msg := fmt.Sprintf(format, v)
	return &NotFoundError{msg}
}

func (e *NotFoundError) Error() string {
	return e.err
}
