package errors

import (
	"fmt"
	"runtime/debug"
)

//go:generate stringer -type=ErrorType -output=error_type_string.go
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
	traceId string
	cause   error
	stack   []byte
	args    []any
}

func NewSysError(code int, key, msg string, t ErrorType) *SystemError {
	return &SystemError{
		code:    code,
		key:     key,
		message: msg,
		typ:     t,
		args:    []any{},
	}
}

func clone(src *SystemError) *SystemError {
	se := SystemError{
		code:    src.code,
		key:     src.key,
		message: src.message,
		typ:     src.typ,
	}
	se.stack = make([]byte, len(src.stack))
	copy(se.stack, se.stack)
	se.args = make([]any, len(src.args))
	copy(se.args, se.args)

	return &se
}

func (e *SystemError) Error() string {
	return fmt.Sprintf(e.message, e.args...)
}

func (e *SystemError) Unwrap() error {
	return e.cause
}

func (e *SystemError) Is(target error) bool {
	if target == nil {
		return false
	}

	if e2, ok := target.(*SystemError); ok {
		return e.code == e2.code &&
			e.key == e2.key &&
			e.message == e2.message &&
			e.typ == e2.typ
	}

	return false
}

func (e *SystemError) WithStack() *SystemError {
	se := clone(e)
	se.stack = debug.Stack()
	return se
}

func (e *SystemError) WithArgs(args ...any) *SystemError {
	se := clone(e)
	se.args = append(se.args, args...)
	return se
}

func (e *SystemError) WithCause(err error) *SystemError {
	se := clone(e)
	se.cause = err
	return se
}
