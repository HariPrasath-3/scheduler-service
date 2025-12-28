package main

import (
	"log"
	"os"
)

const (
	modeApi      = "api"
	modeConsumer = "consumer"
	modeWorker   = "worker"
)

func main() {
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		log.Fatal("APP_MODE not set (api | consumer | worker)")
	}
	app := &application{
		closed: make(chan struct{}),
	}
	defer app.Close()
	initialize(app)

	switch mode {
	case modeApi:
		log.Println("Starting Scheduler API")
		startAPI(app)

	case modeConsumer:
		log.Println("Starting Scheduler Consumer")
		startConsumer(app)

	case modeWorker:
		log.Println("Starting Scheduler Worker")
		startWorker()

	default:
		log.Fatalf("Unknown APP_MODE: %s", mode)
	}
	// Block until shutdown
	app.WaitForClose()
}
