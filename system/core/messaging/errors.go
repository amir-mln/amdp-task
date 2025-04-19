package messaging

import (
	"net/http"

	syserr "github.com/amir-mln/amdp-task/system/errors"
)

var (
	ErrInvalidMessageType = syserr.NewSysError(
		http.StatusInternalServerError,
		"system.core.messaging:invalid message record type",
		"received a message with invalid type of %v",
		syserr.Internal,
	)
)
