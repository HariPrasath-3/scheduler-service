package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/HariPrasath-3/scheduler-service/internal/models"
	"github.com/HariPrasath-3/scheduler-service/internal/repository/dynamo"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	schedulerV1 "github.com/HariPrasath-3/scheduler-service/proto/scheduler/v1"
)

type SchedulerService struct {
	env  *env.Env
	repo dynamo.EventRepository
}

func NewSchedulerService(env *env.Env) *SchedulerService {
	return &SchedulerService{
		env:  env,
		repo: dynamo.NewEventRepository(env),
	}
}

func (s *SchedulerService) Schedule(
	ctx context.Context,
	req *schedulerV1.ScheduleRequest,
) error {
	log.Printf("Received Schedule request: %+v", req)

	event, err := s.repo.Get(ctx, req.GetEventId())
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return fmt.Errorf("failed to get event with ID %s: %v", req.GetEventId(), err)
	}
	if event != nil {
		log.Printf("event already exists: %+v", event)
		return fmt.Errorf("event with ID %s already exists", req.GetEventId())
	}

	event = &models.Event{
		ID:        req.GetEventId(),
		Topic:     req.GetTopic(),
		ExecuteAt: req.GetExecuteAt(),
		Payload:   req.GetPayload(),
	}
	err = s.repo.Save(ctx, event)
	if err != nil {
		log.Printf("failed to save event: %v", err)
		return fmt.Errorf("failed to save event with ID %s: %v", req.GetEventId(), err)
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return err
	}
	err = s.env.Producer().Send(ctx, "schedule_events", event.ID, eventBytes)
	if err != nil {
		log.Printf("failed to send event to kafka: %v", err)
		return fmt.Errorf("failed to enqueue event with ID %s: %v", event.ID, err)
	}

	log.Printf("Scheduled event successfully: %+v", event)
	return nil
}

func (s *SchedulerService) Cancel(
	ctx context.Context,
	req *schedulerV1.CancelRequest,
) error {
	log.Printf("Received Cancel request: %+v", req)

	event, err := s.repo.Get(ctx, req.GetEventId())
	if err != nil {
		log.Printf("failed to get event: %v", err)
		return fmt.Errorf("failed to get event with ID %s: %v", req.GetEventId(), err)
	}
	if event == nil {
		log.Printf("event not found: %s", req.GetEventId())
		return fmt.Errorf("event with ID %s not found", req.GetEventId())
	}

	if event.Status != models.StatusScheduled {
		log.Printf("event not in scheduled state: %+v", event)
		return fmt.Errorf("event with ID %s not in scheduled state", req.GetEventId())
	}

	err = s.repo.UpdateStatus(ctx, event.ID, models.StatusCancelled)
	if err != nil {
		log.Printf("failed to update event status: %v", err)
		return fmt.Errorf("failed to cancel event with ID %s: %v", req.GetEventId(), err)
	}

	log.Printf("Cancelled event successfully: %+v", event)
	return nil
}
