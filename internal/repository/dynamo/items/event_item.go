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

	ID          string `dynamodbav:"id"`
	ReferenceID string `dynamodbav:"reference_id"`
	Topic       string `dynamodbav:"topic"`

	// Execution timestamp (Unix seconds)
	ExecuteAt int64 `dynamodbav:"execute_at"`

	Payload   []byte `dynamodbav:"payload"`
	Status    string `dynamodbav:"status"`
	CreatedAt int64  `dynamodbav:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"`
}

func ToDynamoEventItem(e *models.Event) *DynamoEventItem {
	return &DynamoEventItem{
		PK:          fmt.Sprintf(EventPkFormat, e.ID),
		SK:          EventSk,
		ID:          e.ID,
		ReferenceID: e.ReferenceID,
		Topic:       e.Topic,
		ExecuteAt:   e.ExecuteAt,
		Payload:     e.Payload,
		Status:      string(e.Status),
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func FromDynamoEventItem(item *DynamoEventItem) *models.Event {
	return &models.Event{
		ID:          item.ID,
		ReferenceID: item.ReferenceID,
		Topic:       item.Topic,
		ExecuteAt:   item.ExecuteAt,
		Payload:     item.Payload,
		Status:      models.EventStatus(item.Status),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
