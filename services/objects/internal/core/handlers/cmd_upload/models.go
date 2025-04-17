package cmd_upload

import (
	"io"

	"github.com/google/uuid"
)

type Command struct {
	// Should be parsed from JWT and read from HTTP Request's context
	// It's not supported in the current version
	UserID uint64
	Mime   string
	Name   string
	Object io.Reader
}

type Response struct {
	OID   uuid.UUID `json:"object_id"`
	State string    `json:"state"`
}
