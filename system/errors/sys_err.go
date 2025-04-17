package errors

import (
	"fmt"
	"runtime/debug"
)

type ErrorType uint

const (
	Application ErrorType = iota
	Internal
)

type SystemError struct {
	code    int
	typ     ErrorType
	key     string
	message string
	err     error
	stack   []byte
	args    []any
}

func NewSysError(code int, key, msg string, t ErrorType) SystemError {
	return SystemError{
		code:    code,
		key:     key,
		message: msg,
		typ:     t,
	}
}

func (e SystemError) Error() string {
	return fmt.Sprintf(e.message, e.args...)
}

func (e SystemError) Unwrap() error {
	return e.err
}

func (e SystemError) Is(target error) bool {
	if target == nil {
		return false
	}

	if e2, ok := target.(SystemError); ok {
		return e.code == e2.code &&
			e.key == e2.key &&
			e.message == e2.message &&
			e.typ == e2.typ
	}

	return false
}

func (e SystemError) WithStack() SystemError {
	e.stack = debug.Stack()
	return e
}

func (e SystemError) WithArgs(args ...any) SystemError {
	e.args = args

	return e
}

func (e SystemError) WithError(err error) SystemError {
	e.err = err

	return e
}
