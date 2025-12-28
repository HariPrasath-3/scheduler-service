package scheduler_event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"

	"github.com/HariPrasath-3/scheduler-service/internal/common"
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
	cfg := h.env.Config().SchedulerConfig

	var event models.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// basic validation
	if event.ID == "" {
		return fmt.Errorf("missing event id")
	}
	if event.ExecuteAt <= 0 {
		return fmt.Errorf("invalid execute_at")
	}

	// compute bucket (minute-level)
	bucket := event.ExecuteAt / int64(cfg.BucketSizeSec)

	// compute partition (stable hash)
	partition := computePartition(event.ID, cfg.TotalPartitions)

	// redis key
	redisKey := fmt.Sprintf(
		common.RedisKeyFormatterScheduledEvents,
		bucket,
		partition,
	)

	// enqueue ONLY event_id
	if err := h.env.Redis().
		RPush(ctx, redisKey, event.ID).
		Err(); err != nil {
		return fmt.Errorf("redis enqueue failed: %w", err)
	}

	log.Print("enqueued event ", event.ID,
		" to bucket ", bucket,
		" partition ", partition,
		" redis_key ", redisKey,
	)
	return nil
}

func computePartition(eventID string, totalPartitions int) uint32 {
	var hash uint32
	for i := 0; i < len(eventID); i++ {
		hash = hash*31 + uint32(eventID[i])
	}
	return hash % uint32(totalPartitions)
}
