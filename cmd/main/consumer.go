package main

import (
	"log"

	"github.com/HariPrasath-3/scheduler-service/internal/consumer/scheduler_event"
)

func startConsumer(app *application) {
	ctx := app.confCtx
	cfg := app.appConfig

	// Create scheduler consumer service
	svc, err := scheduler_event.NewSchedulerConsumer(
		ctx,
		app.env,
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
	)
	if err != nil {
		log.Fatalf("failed to create scheduler consumer: %v", err)
	}
	svc.Start(ctx)
	app.AddShutdownCallback(func() {
		_ = svc.Stop()
	})

	log.Println("Scheduler consumer started")
}
