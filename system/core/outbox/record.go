package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Record struct {
	id          uint64
	UUID        uuid.UUID
	Entity      string
	EntityID    uint64
	Title       string
	Target      string
	Payload     any
	CreatedAt   time.Time
	ProcessedAt *time.Time
}
