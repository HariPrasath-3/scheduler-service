package controller

import (
	"context"
	"log"

	schedulev1 "github.com/HariPrasath-3/scheduler-service/client/golang/proto/github.com/HariPrasath-3/scheduler-service/proto/scheduler/v1"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
)

type SchedulerController struct {
	schedulev1.UnimplementedSchedulerServiceServer
	env *env.Env
}

func NewSchedulerController(env *env.Env) *SchedulerController {
	return &SchedulerController{env: env}
}

func (c *SchedulerController) Schedule(
	ctx context.Context,
	req *schedulev1.ScheduleRequest,
) (*schedulev1.ScheduleResponse, error) {
	log.Printf("Received Schedule request: %+v", req)
	return &schedulev1.ScheduleResponse{Accepted: true}, nil
}
