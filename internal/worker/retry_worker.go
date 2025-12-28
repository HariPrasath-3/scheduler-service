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

type RetryWorker struct {
	env  *env.Env
	repo dynamo.EventRepository

	workerID string
}

func NewRetryWorker(
	env *env.Env,
) *RetryWorker {
	return &RetryWorker{
		env:      env,
		repo:     dynamo.NewEventRepository(env),
		workerID: uuid.NewString(),
	}
}

func (w *RetryWorker) Start(ctx context.Context) {
	log.Println("retry worker started")

	ticker := time.NewTicker(time.Duration(w.env.Config().RetryWorkerConfig.RetryScanIntervalMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("retry worker stopped")
			return
		case <-ticker.C:
			w.retryOnce(ctx)
		}
	}
}

func (w *RetryWorker) retryOnce(ctx context.Context) {
	now := time.Now().Unix()

	for partition := uint32(0); partition < uint32(w.env.Config().SchedulerConfig.TotalPartitions); partition++ {
		w.retryPartition(ctx, partition, now)
	}
}

func (w *RetryWorker) retryPartition(
	ctx context.Context,
	partition uint32,
	now int64,
) {
	currentBucket := now / int64(w.env.Config().SchedulerConfig.BucketSizeSec)
	pastBucketCount := int64(w.env.Config().RetryWorkerConfig.PastBucketsCount)

	for bucket := currentBucket; bucket >= currentBucket-pastBucketCount; bucket-- {
		processingKey := fmt.Sprintf(
			common.RedisKeyFormatterProcessingEvents,
			bucket,
			partition,
		)

		retryBatchSize := int64(w.env.Config().RetryWorkerConfig.RetryBatchSize)
		eventIDs, err := w.env.Redis().
			LRange(ctx, processingKey, 0, retryBatchSize-1).
			Result()
		if err != nil || len(eventIDs) == 0 {
			continue
		}

		for _, eventID := range eventIDs {
			w.retryEvent(ctx, bucket, partition, eventID)
		}
	}
}

func (w *RetryWorker) retryEvent(
	ctx context.Context,
	bucket int64,
	partition uint32,
	eventID string,
) {
	event, err := w.repo.Get(ctx, eventID)
	if err != nil {
		return
	}

	// If already processed, clean it up
	if event.Status != models.StatusScheduled {
		w.ackEvent(ctx, bucket, partition, eventID)
		return
	}

	// Move back to scheduled
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

	pipe := w.env.Redis().Pipeline()
	pipe.LRem(ctx, processingKey, 1, eventID)
	pipe.RPush(ctx, scheduledKey, eventID)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("failed to retry event %s: %v", eventID, err)
		return
	}

	log.Printf("retried event %s", eventID)
}

func (w *RetryWorker) ackEvent(
	ctx context.Context,
	bucket int64,
	partition uint32,
	eventID string,
) {
	processingKey := fmt.Sprintf(
		common.RedisKeyFormatterProcessingEvents,
		bucket,
		partition,
	)

	err := w.env.Redis().LRem(ctx, processingKey, 1, eventID).Err()
	if err != nil {
		log.Printf("failed to ack event %s: %v", eventID, err)
	}
}
