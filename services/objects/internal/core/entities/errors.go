package entities

import (
	"net/http"

	syserr "github.com/amir-mln/amdp-task/system/errors"
)

var (
	ErrReadingNilObject = syserr.NewSysError(
		http.StatusBadRequest,
		"objects.core.entities:nil-object-reader",
		"the underlying io data was nil",
		syserr.Application,
	)
	ErrInvalidObjectState = syserr.NewSysError(
		http.StatusInternalServerError,
		"objects.core.entities:invalid-object-state",
		"received an invalid object state of value %v",
		syserr.Internal,
	)
)
