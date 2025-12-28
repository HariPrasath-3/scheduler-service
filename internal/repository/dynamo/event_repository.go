package dynamo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/HariPrasath-3/scheduler-service/internal/models"
	"github.com/HariPrasath-3/scheduler-service/internal/repository/dynamo/items"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type EventRepository interface {
	Save(ctx context.Context, event *models.Event) error
	Get(ctx context.Context, eventID string) (*models.Event, error)
	UpdateStatus(ctx context.Context, eventID string, newStatus models.EventStatus) error
}

type EventRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
}

func NewEventRepository(
	environment *env.Env,
) EventRepository {
	return &EventRepositoryImpl{
		client:    environment.Dynamo(),
		tableName: environment.Config().Dynamo.TableName,
	}
}

// Save saves a new event to DynamoDB. If an event with the same ID already exists, it ignores the error.
func (r *EventRepositoryImpl) Save(
	ctx context.Context,
	event *models.Event,
) error {

	now := time.Now().Unix()
	event.CreatedAt = now
	event.UpdatedAt = now
	event.Status = models.StatusScheduled

	item := items.ToDynamoEventItem(event)

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           &r.tableName,
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(pk)"),
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			log.Printf("item already exists: %v", err)
			return nil
		}
	}
	return err
}

func (r *EventRepositoryImpl) Get(
	ctx context.Context,
	eventID string,
) (*models.Event, error) {
	pk := fmt.Sprintf(items.EventPkFormat, eventID)
	sk := items.EventSk

	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	var item items.DynamoEventItem
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	event := items.FromDynamoEventItem(&item)
	return event, nil
}

func (r *EventRepositoryImpl) UpdateStatus(
	ctx context.Context,
	eventID string,
	newStatus models.EventStatus,
) error {
	pk := fmt.Sprintf(items.EventPkFormat, eventID)
	sk := items.EventSk

	now := time.Now().Unix()

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("SET #status = :newStatus, updated_at = :updatedAt"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newStatus": &types.AttributeValueMemberS{Value: string(newStatus)},
			":updatedAt": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", now)},
		},
	})
	return err
}
