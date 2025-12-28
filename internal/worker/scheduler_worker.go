package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/HariPrasath-3/scheduler-service/internal/common"
	"github.com/HariPrasath-3/scheduler-service/internal/models"
	"github.com/HariPrasath-3/scheduler-service/internal/repository/dynamo"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	"github.com/google/uuid"
)

type SchedulerWorker struct {
	env      *env.Env
	repo     dynamo.EventRepository
	workerID string
}

func NewSchedulerWorker(env *env.Env) *SchedulerWorker {
	return &SchedulerWorker{
		env:      env,
		repo:     dynamo.NewEventRepository(env),
		workerID: uuid.NewString(),
	}
}

func (w *SchedulerWorker) Start(ctx context.Context) {
	log.Printf("starting scheduler worker %s", w.workerID)

	cfg := w.env.Config().WorkerConfig

	ticker := time.NewTicker(time.Duration(cfg.PollIntervalMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %s shutting down", w.workerID)
			return

		case <-ticker.C:
			w.processTick(ctx)
		}
	}
}

func (w *SchedulerWorker) processTick(ctx context.Context) {
	schedulerConfig := w.env.Config().SchedulerConfig
	workerConfig := w.env.Config().WorkerConfig

	now := time.Now().Unix()
	currentBucket := now / int64(schedulerConfig.BucketSizeSec)

	for bucket := currentBucket; bucket >= currentBucket-int64(workerConfig.PastBucketsCount); bucket-- {
		for partition := uint32(0); partition < uint32(schedulerConfig.TotalPartitions); partition++ {
			w.processPartition(ctx, bucket, partition)
		}
	}
}

func (w *SchedulerWorker) processPartition(
	ctx context.Context,
	bucket int64,
	partition uint32,
) {

	redisKey := fmt.Sprintf(
		common.RedisKeyFormatterScheduledEvents,
		bucket,
		partition,
	)

	eventID, err := w.env.Redis().LPop(ctx, redisKey).Result()
	if err != nil {
		return // empty or transient error
	}

	w.processEvent(ctx, eventID)
}

func (w *SchedulerWorker) processEvent(
	ctx context.Context,
	eventID string,
) {

	event, err := w.repo.Get(ctx, eventID)
	if err != nil {
		log.Printf("event %s not found, skipping", eventID)
		return
	}

	// status check
	if event.Status != models.StatusScheduled {
		return
	}

	// execute_at safety
	if event.ExecuteAt > time.Now().Unix() {
		// not due yet → skip (or requeue later)
		return
	}

	// publish to Kafka (log-only for now)
	log.Printf(
		"worker %s executing event %s topic=%s",
		w.workerID,
		event.ID,
		event.Topic,
	)

	// TODO: kafkaProducer.Publish(event.Topic, event.Payload)

	// mark as FIRED
	if err := w.repo.UpdateStatus(
		ctx,
		event.ID,
		models.StatusFired,
	); err != nil {
		log.Printf("failed to mark event %s as FIRED", event.ID)
	}
}
