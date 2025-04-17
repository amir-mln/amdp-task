package messaging

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type MessageHeader struct {
	TxID        uuid.UUID   `json:"transaction_id"`
	Entity      string      `json:"entity"`
	EntityID    uint64      `json:"entity_id"`
	UserId      *uint64     `json:"user_id"`
	Title       string      `json:"title"`
	Type        MessageType `json:"message_type"`
	CreatedAt   time.Time   `json:"created_at"`
	PublishedAt time.Time   `json:"published_at"`
}

type MessageBody interface {
	MessageTitle() string
	MessageType() MessageType
}

type Message struct {
	Header MessageHeader
	Body   []byte
	// the below fields are for database use only
	id        uint64
	publishAt *time.Time
}

func NewMessage(body MessageBody, opts ...MessageOption) *Message {
	b, err := json.Marshal(any(body))
	if err != nil {
		panic(err)
	}

	ob := &Message{
		Header: MessageHeader{
			TxID:      uuid.New(),
			Entity:    "",
			EntityID:  0,
			Title:     body.MessageTitle(),
			Type:      body.MessageType(),
			CreatedAt: time.Now(),
		},
		Body: b,
	}
	for _, opt := range opts {
		opt(ob)
	}

	return ob
}
