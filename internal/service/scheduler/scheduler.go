package scheduler

import (
	"context"
	"encoding/json"
	"log"

	schedulev1 "github.com/HariPrasath-3/scheduler-service/client/golang/proto/github.com/HariPrasath-3/scheduler-service/proto/scheduler/v1"
	"github.com/HariPrasath-3/scheduler-service/internal/models"
	"github.com/HariPrasath-3/scheduler-service/internal/repository/dynamo"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
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
	req *schedulev1.ScheduleRequest,
) error {
	log.Printf("Received Schedule request: %+v", req)

	event := &models.Event{
		ReferenceID: req.GetReferenceId(),
		Topic:       req.GetTopic(),
		ExecuteAt:   req.GetExecuteAt(),
		Payload:     req.GetPayload(),
	}
	event.GenerateId()
	err := s.repo.Save(ctx, event)
	if err != nil {
		log.Printf("failed to save event: %v", err)
		return err
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return err
	}
	err = s.env.Producer().Send(ctx, "schedule_events", event.ID, eventBytes)
	if err != nil {
		log.Printf("failed to send event to kafka: %v", err)
		return err
	}

	log.Printf("Scheduled event successfully: %+v", event)
	return nil
}
