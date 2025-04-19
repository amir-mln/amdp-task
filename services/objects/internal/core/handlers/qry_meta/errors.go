package qry_meta

import (
	"net/http"

	syserr "github.com/amir-mln/amdp-task/system/errors"
)

var (
	ErrInvalidRequestObjectID = syserr.NewSysError(
		http.StatusBadRequest,
		"objects.core.qry_meta:invalid-uuid",
		"received an invalid uuid of %q from request",
		syserr.Application,
	)
)
