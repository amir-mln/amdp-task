package messaging

import (
	"time"

	"github.com/google/uuid"
)

type MessageOption func(*Message)

func WithTxID(ud uuid.UUID) MessageOption {
	return func(msg *Message) {
		msg.ID = ud
	}
}

func WithUserID(id int64) MessageOption {
	return func(msg *Message) {
		msg.UserId = &id
	}
}

func WithCreateTime(t time.Time) MessageOption {
	return func(msg *Message) {
		msg.CreatedAt = t
	}
}

func WithPublishTime(t *time.Time) MessageOption {
	return func(msg *Message) {
		msg.publishAt = t
	}
}

type MessageEntity interface {
	EntityName() string
	EntityID() int64
}

func WithEntity(ent MessageEntity) MessageOption {
	return func(m *Message) {
		n, id := ent.EntityName(), ent.EntityID()
		m.Entity = &n
		m.EntityID = &id
	}
}
