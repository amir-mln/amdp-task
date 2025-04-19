package errors

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	internalErrorMsg = "key: %q, trace id: %q"
)

var (
	unexpectedError = NewSysError(
		http.StatusInternalServerError,
		"system.errors:unexpected-internal-error",
		internalErrorMsg,
		Internal,
	)
)

func wrapForLogging(err error) *SystemError {
	if err == nil {
		return nil
	}

	var se *SystemError
	trace := uuid.New().String()
	if !errors.As(err, &se) {
		se = unexpectedError.WithCause(err).WithArgs(internalErrorMsg, trace)
		se.traceId = trace
	} else if se.typ == Internal {
		newErr := clone(se)
		newErr.message = internalErrorMsg
		newErr.args = []any{internalErrorMsg, trace}
		newErr.traceId = trace
		newErr.cause = se
		se.stack = nil
		se = newErr
	}

	return se
}

func HandleHTTPError(w http.ResponseWriter, err error, log *zap.Logger) {
	se := wrapForLogging(err)
	if se.typ == Internal {
		LogSystemError(se, log)
	}

	http.Error(w, se.Error(), se.code)
}

func LogSystemError(se *SystemError, log *zap.Logger) {
	if se == nil {
		return
	}

	fields := []zap.Field{
		zap.Stringer("Type", se.typ),
		zap.String("Key", se.key),
		zap.String("Trace ID", se.traceId),
		zap.NamedError("Cause", se.cause),
		zap.ByteString("Stack", se.stack),
	}
	log.Error("Received an error", fields...)
}

func LogError(err error, log *zap.Logger) {
	se := wrapForLogging(err)
	if se == nil {
		return
	}
	LogSystemError(se, log)
}
