package validator

import (
	"fmt"
	"time"

	schedulerV1 "github.com/HariPrasath-3/scheduler-service/proto/scheduler/v1"
	"github.com/google/uuid"
)

type SchedulerValidator struct{}

func NewSchedulerValidator() *SchedulerValidator {
	return &SchedulerValidator{}
}

func (v *SchedulerValidator) ValidateScheduleRequest(req *schedulerV1.ScheduleRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate event_id (must be valid UUID)
	if req.GetEventId() == "" {
		return fmt.Errorf("event_id is required")
	}
	if _, err := uuid.Parse(req.GetEventId()); err != nil {
		return fmt.Errorf("event_id must be a valid UUID")
	}

	// Validate execute_at (must be in the future)
	if req.GetExecuteAt() <= 0 {
		return fmt.Errorf("execute_at is required and must be positive")
	}
	if req.GetExecuteAt() <= time.Now().Unix() {
		return fmt.Errorf("execute_at must be in the future")
	}

	// Validate topic
	if req.GetTopic() == "" {
		return fmt.Errorf("topic is required")
	}

	// Validate payload
	if len(req.GetPayload()) == 0 {
		return fmt.Errorf("payload is required")
	}

	return nil
}

func (v *SchedulerValidator) ValidateCancelRequest(req *schedulerV1.CancelRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate event_id
	if req.GetEventId() == "" {
		return fmt.Errorf("event_id is required")
	}
	if _, err := uuid.Parse(req.GetEventId()); err != nil {
		return fmt.Errorf("event_id must be a valid UUID")
	}

	return nil
}
