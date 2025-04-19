package messaging

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessagePayload interface {
	MessageTitle() string
	MessageType() MessageType
}

type Message struct {
	ID          uuid.UUID   `json:"transaction_id"`
	Entity      *string     `json:"entity"`
	EntityID    *int64      `json:"entity_id"`
	UserId      *int64      `json:"user_id"`
	Title       string      `json:"title"`
	Type        MessageType `json:"message_type"`
	CreatedAt   time.Time   `json:"created_at"`
	PublishedAt time.Time   `json:"published_at"`
	Body        []byte      `json:"-"`
	// the below fields are for database use only
	publishAt *time.Time
}

// this will panic if body argument can not be JSON encoded
func NewMessage(p MessagePayload, opts ...MessageOption) *Message {
	b, err := json.Marshal(any(p))
	if err != nil {
		panic(err)
	}

	ob := &Message{
		ID:        uuid.New(),
		Title:     p.MessageTitle(),
		Type:      p.MessageType(),
		Body:      b,
		CreatedAt: time.Now(),
	}
	for _, opt := range opts {
		opt(ob)
	}

	return ob
}
