package items

import (
	"fmt"

	"github.com/HariPrasath-3/scheduler-service/internal/models"
)

const (
	EventPkFormat = "event#%s"
	EventSk       = "meta"
)

type DynamoEventItem struct {
	PK string `dynamodbav:"pk"`
	SK string `dynamodbav:"sk"`

	ID        string `dynamodbav:"id"`
	Topic     string `dynamodbav:"topic"`
	ExecuteAt int64  `dynamodbav:"execute_at"` // Execution timestamp (Unix seconds)
	Payload   []byte `dynamodbav:"payload"`
	Status    string `dynamodbav:"status"`
	CreatedAt int64  `dynamodbav:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}

func ToDynamoEventItem(e *models.Event) *DynamoEventItem {
	return &DynamoEventItem{
		PK:        fmt.Sprintf(EventPkFormat, e.ID),
		SK:        EventSk,
		ID:        e.ID,
		Topic:     e.Topic,
		ExecuteAt: e.ExecuteAt,
		Payload:   e.Payload,
		Status:    string(e.Status),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func FromDynamoEventItem(item *DynamoEventItem) *models.Event {
	return &models.Event{
		ID:        item.ID,
		Topic:     item.Topic,
		ExecuteAt: item.ExecuteAt,
		Payload:   item.Payload,
		Status:    models.EventStatus(item.Status),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
