package entities

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"hash"
	"io"
	"time"

	"github.com/google/uuid"
)

//go:generate stringer -type=ObjectState -output=obj_state_string.go
type ObjectState uint

const (
	Initial ObjectState = iota
	Completed
	Failed
)

func (os ObjectState) MarshalText() (text []byte, err error) {
	s := os.String()
	if !os.Valid() {
		return nil, ErrInvalidObjectState.WithArgs(s)
	}

	return []byte(s), nil
}

func (os *ObjectState) UnmarshalText(text []byte) error {
	switch text := string(text); text {
	case Initial.String():
		*os = Initial
		return nil
	case Completed.String():
		*os = Completed
		return nil
	case Failed.String():
		*os = Failed
		return nil
	default:
		return ErrInvalidObjectState.WithArgs(os)
	}
}

func (os *ObjectState) Scan(v any) error {
	switch v := v.(type) {
	case []byte:
		return os.UnmarshalText(v)
	case string:
		return os.UnmarshalText([]byte(v))
	default:
		return ErrInvalidObjectState.WithArgs(v)
	}
}

func (os ObjectState) Value() (driver.Value, error) {
	s := os.String()
	if !os.Valid() {
		return nil, ErrInvalidObjectState.WithArgs(s)
	}

	return s, nil
}

func (os ObjectState) Valid() bool {
	return Initial <= os && os <= Completed
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
		return 0, ErrReadingNilObject
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

func (o *Object) EntityID() uint64 {
	return o.ID
}

func (o *Object) EntityName() string {
	return "Object"
}
