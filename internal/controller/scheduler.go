package controller

import (
	"context"

	schedulev1 "github.com/HariPrasath-3/scheduler-service/client/golang/proto/github.com/HariPrasath-3/scheduler-service/proto/scheduler/v1"
	"github.com/HariPrasath-3/scheduler-service/internal/service"
	"github.com/HariPrasath-3/scheduler-service/internal/service/scheduler"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
)

type SchedulerController struct {
	schedulev1.UnimplementedSchedulerServiceServer
	env              *env.Env
	schedulerService *scheduler.SchedulerService
}

func NewSchedulerController(env *env.Env) *SchedulerController {
	return &SchedulerController{
		env:              env,
		schedulerService: service.GetServiceFactory(env).SchedulerService,
	}
}

func (c *SchedulerController) Schedule(
	ctx context.Context,
	req *schedulev1.ScheduleRequest,
) (*schedulev1.ScheduleResponse, error) {
	err := c.schedulerService.Schedule(ctx, req)
	if err != nil {
		return nil, err
	}
	return &schedulev1.ScheduleResponse{Accepted: true}, nil
}
