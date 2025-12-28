package scheduler_event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"

	"github.com/HariPrasath-3/scheduler-service/internal/models"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	"github.com/HariPrasath-3/scheduler-service/pkg/kafka"
)

type ScheduleEventHandler struct {
	env *env.Env
}

func NewScheduleEventHandler(env *env.Env) kafka.MessageHandler {
	return &ScheduleEventHandler{
		env: env,
	}
}

func (h *ScheduleEventHandler) HandleMessage(
	ctx context.Context,
	msg *sarama.ConsumerMessage,
) error {

	var event models.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Received schedule event: %+v", event)

	return nil
}
