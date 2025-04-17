package messaging

import "database/sql/driver"

//go:generate stringer -type=MessageType -output=msg_type_string.go
type MessageType uint

const (
	Event MessageType = iota
	Command
	Query
)

func (mt MessageType) MarshalText() (text []byte, err error) {
	if !mt.Valid() {
		return nil, ErrInvalidMessageType.WithArgs(mt)
	}

	return []byte(mt.String()), nil
}

func (mt *MessageType) UnmarshalText(text []byte) error {
	switch text := string(text); text {
	case Event.String():
		*mt = Event
		return nil
	case Command.String():
		*mt = Command
		return nil
	case Query.String():
		*mt = Query
		return nil
	default:
		return ErrInvalidMessageType.WithArgs(text)
	}
}

func (mt *MessageType) Scan(v any) error {
	switch v := v.(type) {
	case []byte:
		return mt.UnmarshalText(v)
	case string:
		return mt.UnmarshalText([]byte(v))
	default:
		return ErrInvalidMessageType.WithArgs(v)
	}
}

func (mt MessageType) Value() (driver.Value, error) {
	s := mt.String()
	if !mt.Valid() {
		return nil, ErrInvalidMessageType.WithArgs(s)
	}

	return s, nil
}

func (mt MessageType) Valid() bool {
	return Event <= mt && mt <= Query
}
