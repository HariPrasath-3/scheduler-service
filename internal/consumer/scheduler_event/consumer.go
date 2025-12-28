package scheduler_event

import (
	"context"
	"log"

	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	"github.com/HariPrasath-3/scheduler-service/pkg/kafka"
)

type SchedulerConsumer struct {
	consumer *kafka.Consumer
}

func NewSchedulerConsumer(
	ctx context.Context,
	env *env.Env,
	brokers []string,
	groupID string,
) (*SchedulerConsumer, error) {

	handler := NewScheduleEventHandler(env)

	consumer, err := kafka.NewConsumer(
		brokers,
		groupID,
		[]string{"schedule_events"},
		handler,
	)
	if err != nil {
		return nil, err
	}

	return &SchedulerConsumer{
		consumer: consumer,
	}, nil
}

func (s *SchedulerConsumer) Start(ctx context.Context) {
	log.Println("starting scheduler consumer")
	s.consumer.Start(ctx)
}

func (s *SchedulerConsumer) Stop() error {
	log.Println("stopping scheduler consumer")
	return s.consumer.Close()
}