package messaging

import (
	"time"

	"github.com/google/uuid"
)

type MessageOption func(*Message)

func WithTxID(ud uuid.UUID) MessageOption {
	return func(msg *Message) {
		msg.Header.TxID = ud
	}
}

func WithUserID(id uint64) MessageOption {
	return func(msg *Message) {
		msg.Header.UserId = &id
	}
}

func WithCreateTime(t time.Time) MessageOption {
	return func(o *Message) {
		o.Header.CreatedAt = t
	}
}

func WithPublishTime(t *time.Time) MessageOption {
	return func(o *Message) {
		o.publishAt = t
	}
}

type MessageEntity interface {
	EntityName() string
	EntityID() uint64
}

func WithEntity(ent MessageEntity) MessageOption {
	return func(m *Message) {
		m.Header.Entity = ent.EntityName()
		m.Header.EntityID = ent.EntityID()
	}
}
