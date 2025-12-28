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
	env  *env.Env
	repo dynamo.EventRepository

	workerID string
	sem      chan struct{} // global concurrency limiter
}

func NewSchedulerWorker(
	env *env.Env,
) *SchedulerWorker {
	return &SchedulerWorker{
		env:      env,
		repo:     dynamo.NewEventRepository(env),
		workerID: uuid.NewString(),
		sem:      make(chan struct{}, env.Config().WorkerConfig.SemaphoreLimit),
	}
}

func (w *SchedulerWorker) Start(ctx context.Context) {
	log.Println("scheduler worker started")

	for p := uint32(0); p < uint32(w.env.Config().SchedulerConfig.TotalPartitions); p++ {
		go w.runPartition(ctx, p)
	}

	<-ctx.Done()
	log.Println("scheduler worker stopped")
}

func (w *SchedulerWorker) runPartition(
	ctx context.Context,
	partition uint32,
) {
	log.Printf("partition worker %d started", partition)

	backoff := 10 * time.Millisecond
	maxBackoff := 1 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		moved := w.drainPartition(ctx, partition)

		// no work found, exponential backoff
		if moved == 0 {
			time.Sleep(backoff)
			if backoff < maxBackoff {
				backoff *= 2
			}
			continue
		}

		// reset backoff when work was found
		backoff = 10 * time.Millisecond
	}
}

func (w *SchedulerWorker) drainPartition(
	ctx context.Context,
	partition uint32,
) int {

	now := time.Now().Unix()
	currentBucket := now / int64(w.env.Config().SchedulerConfig.BucketSizeSec)
	pastBucketCount := int64(w.env.Config().WorkerConfig.PastBucketsCount)

	moved := 0

	for bucket := currentBucket; bucket >= currentBucket-pastBucketCount; bucket-- {
		for {
			// 🔑 global concurrency guard
			select {
			case w.sem <- struct{}{}:
				// acquired
			default:
				return moved
			}

			scheduledKey := fmt.Sprintf(
				common.RedisKeyFormatterScheduledEvents,
				bucket,
				partition,
			)

			processingKey := fmt.Sprintf(
				common.RedisKeyFormatterProcessingEvents,
				bucket,
				partition,
			)

			eventID, err := w.env.Redis().
				LMove(ctx, scheduledKey, processingKey, "LEFT", "RIGHT").
				Result()

			if err != nil {
				<-w.sem // release slot
				return moved
			}

			moved++

			go func(eid string, b int64) {
				defer func() { <-w.sem }()
				w.processEvent(ctx, b, partition, eid)
			}(eventID, bucket)

			if moved >= w.env.Config().WorkerConfig.BatchSize {
				return moved
			}
		}
	}

	return moved
}

func (w *SchedulerWorker) processEvent(
	ctx context.Context,
	bucket int64,
	partition uint32,
	eventID string,
) {

	event, err := w.repo.Get(ctx, eventID)
	if err != nil {
		return
	}

	// idempotency check
	if event.Status != models.StatusScheduled {
		w.ackEvent(ctx, bucket, partition, eventID)
		return
	}

	// 🔥 execute (Kafka publish later)
	log.Printf("executing event %s", event.ID)

	// mark fired
	if err := w.repo.UpdateStatus(ctx, event.ID, models.StatusFired); err != nil {
		return
	}

	w.ackEvent(ctx, bucket, partition, eventID)
}

func (w *SchedulerWorker) ackEvent(
	ctx context.Context,
	bucket int64,
	partition uint32,
	eventID string,
) {
	processingKey := fmt.Sprintf(
		"scheduler:processing:%d:%d",
		bucket,
		partition,
	)

	_ = w.env.Redis().LRem(ctx, processingKey, 1, eventID).Err()
}
