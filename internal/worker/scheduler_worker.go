package worker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/HariPrasath-3/scheduler-service/internal/common"
	"github.com/HariPrasath-3/scheduler-service/internal/models"
	"github.com/HariPrasath-3/scheduler-service/internal/repository/dynamo"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SchedulerWorker struct {
	env  *env.Env
	repo dynamo.EventRepository

	workerID string
	sem      chan struct{} // global concurrency limiter
	wg       sync.WaitGroup
}

func NewSchedulerWorker(
	env *env.Env,
) *SchedulerWorker {
	return &SchedulerWorker{
		env:      env,
		repo:     dynamo.NewEventRepository(env),
		workerID: uuid.NewString(),
		sem:      make(chan struct{}, env.Config().SchedulerWorkerConfig.SemaphoreLimit),
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

	initialBackoff := time.Duration(w.env.Config().SchedulerWorkerConfig.BackoffMs) * time.Millisecond
	backoff := initialBackoff
	maxBackoff := time.Duration(w.env.Config().SchedulerWorkerConfig.MaxBackoffMs) * time.Millisecond

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
		backoff = initialBackoff
	}
}

func (w *SchedulerWorker) drainPartition(
	ctx context.Context,
	partition uint32,
) int {

	now := time.Now().Unix()
	bucketSize := int64(w.env.Config().SchedulerConfig.BucketSizeSec)
	currentBucket := now / bucketSize

	pastBucketCount := int64(w.env.Config().SchedulerWorkerConfig.PastBucketsCount)
	batchSize := w.env.Config().SchedulerWorkerConfig.BatchSize

	moved := 0

	// Acquire ONE semaphore token for ONE batch
	select {
	case w.sem <- struct{}{}:
		// token acquired
	default:
		return 0
	}

	// Batch can span multiple buckets (same partition)
	eventIDs := make([]string, 0, batchSize)

	// Scan buckets in time order: current → past
	for bucket := currentBucket; bucket >= currentBucket-pastBucketCount; bucket-- {

		if len(eventIDs) >= batchSize {
			break
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

		remaining := batchSize - len(eventIDs)

		// Pipeline LMOVE to reduce RTT
		pipe := w.env.Redis().Pipeline()
		cmds := make([]*redis.StringCmd, 0, remaining)

		for range remaining {
			cmds = append(
				cmds,
				pipe.LMove(ctx, scheduledKey, processingKey, "LEFT", "RIGHT"),
			)
		}

		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf(
				"pipeline failed partition=%d bucket=%d err=%v",
				partition, bucket, err,
			)
			break
		}

		for _, cmd := range cmds {
			eventID, err := cmd.Result()
			if err == redis.Nil {
				break // Source list exhausted — no further LMOVEs can succeed
			}
			if err != nil {
				// Command-specific failure; other LMOVEs may have succeeded
				log.Printf(
					"LMOVE failed partition=%d bucket=%d err=%v",
					partition, bucket, err,
				)
				continue
			}

			// This event was successfully moved to processing
			eventIDs = append(eventIDs, eventID)
		}
	}

	// No work claimed → release token
	if len(eventIDs) == 0 {
		<-w.sem
		return 0
	}

	moved = len(eventIDs)

	// Batch goroutine OWNS the semaphore token
	w.wg.Add(1)
	go func(batch []string) {
		defer w.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				log.Printf(
					"[PANIC] worker=%s partition=%d err=%v",
					w.workerID, partition, r,
				)
			}
			<-w.sem // release token when batch finishes
		}()

		// Sequential processing preserves simplicity & retry safety
		for _, eventID := range batch {
			w.processEvent(ctx, currentBucket, partition, eventID)
		}
	}(eventIDs)

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
		log.Printf("failed to get event %s: %v", eventID, err)
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
		log.Printf("failed to update status for event %s: %v", event.ID, err)
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
		common.RedisKeyFormatterProcessingEvents,
		bucket,
		partition,
	)

	err := w.env.Redis().LRem(ctx, processingKey, 1, eventID).Err()
	if err != nil {
		log.Printf("failed to ack event %s: %v", eventID, err)
	}
}
