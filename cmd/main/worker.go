package main

import (
	"log"

	"github.com/HariPrasath-3/scheduler-service/internal/worker"
)

func startWorker(app *application) {
	log.Println("starting scheduler worker")

	w := worker.NewSchedulerWorker(app.env)
	go w.Start(app.confCtx)
	app.AddShutdownCallback(func() {
		log.Println("stopping scheduler worker")
	})

	log.Println("Scheduler worker started")
}
