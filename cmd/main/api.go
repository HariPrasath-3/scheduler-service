package main

import (
	"log"

	schedulev1 "github.com/HariPrasath-3/scheduler-service/client/golang/proto/github.com/HariPrasath-3/scheduler-service/proto/scheduler/v1"
	"github.com/HariPrasath-3/scheduler-service/internal/controller"
	"github.com/HariPrasath-3/scheduler-service/pkg/grpc"
)

func startAPI(app *application) {
	ctx := app.confCtx

	grpcServer, err := grpc.NewGrpcServer(ctx, app.env)
	if err != nil {
		log.Fatalf("failed to create grpc server: %v", err)
	}

	schedulev1.RegisterSchedulerServiceServer(grpcServer, controller.NewSchedulerController(app.env))

	app.AddShutdownCallback(func() {
		log.Println("gracefully stopping grpc server")
		grpcServer.GracefulStop()
	})

	err = grpc.Serve(ctx, grpcServer, &app.appConfig.Grpc)
	if err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}

	log.Println("Scheduler API started")
}
