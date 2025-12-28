package main

import (
	"log"

	"github.com/HariPrasath-3/scheduler-service/internal/worker"
)

func startWorker(app *application) {
	log.Println("starting scheduler worker")

	schedulerWorker := worker.NewSchedulerWorker(app.env)
	go schedulerWorker.Start(app.confCtx)
	app.AddShutdownCallback(func() {
		log.Println("stopping scheduler worker")
	})

	retryWorker := worker.NewRetryWorker(app.env)
	go retryWorker.Start(app.confCtx)
	app.AddShutdownCallback(func() {
		log.Println("stopping retry worker")
	})

	log.Println("Scheduler worker started")
}
