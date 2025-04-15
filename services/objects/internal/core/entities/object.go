package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"time"

	"github.com/google/uuid"
)

//go:generate stringer -type=ObjectState -output=obj_state_string.go
type ObjectState uint

const (
	Initial ObjectState = iota
	Complete
	Failed
)

func (o *ObjectState) Scan(v any) error {
	var err error
	if s, ok := v.(string); ok {
		switch s {
		case Initial.String():
			*o = Initial
		case Complete.String():
			*o = Complete
		case Failed.String():
			*o = Failed
		default:
			err = fmt.Errorf("")
		}
	} else {
		err = fmt.Errorf("")
	}
	return err
}

type Object struct {
	ID        uint64
	OID       uuid.UUID
	UserID    uint64
	Name      string
	Mime      string
	Size      uint64
	Hash      string
	State     ObjectState
	CreatedAt time.Time
	r         io.Reader
	h         hash.Hash
}

func NewObject(uid uint64, name, mime string, r io.Reader) *Object {
	return &Object{
		OID:       uuid.New(),
		UserID:    uid,
		Name:      name,
		Mime:      mime,
		Size:      0,
		Hash:      "",
		State:     Initial,
		CreatedAt: time.Now(),
		r:         r,
		h:         sha256.New(),
	}
}

// the object storage can use this method to store the object
// asynchronously in streams. It also behaves similar to [io.TeeReader]
// and will write data to underlying [Object] hash.
func (o *Object) Read(p []byte) (n int, err error) {
	if o.r == nil {
		return 0, fmt.Errorf("") // TODO: Custom Error
	}
	if o.h == nil {
		o.h = sha256.New()
	}

	n, err = o.r.Read(p)
	if n > 0 {
		// the [Write] method of [hash.Hash] never returns error
		// so it's safe to ignore its returned values
		_, _ = o.h.Write(p[:n])
		o.Size += uint64(n)
	}
	if err == io.EOF {
		o.Hash = hex.EncodeToString(o.h.Sum(nil))
	}

	return
}
