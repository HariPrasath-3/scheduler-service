package controller

import (
	"context"

	common "github.com/HariPrasath-3/scheduler-service/client/golang/proto/github.com/HariPrasath-3/scheduler-service/proto/common"
	schedulev1 "github.com/HariPrasath-3/scheduler-service/client/golang/proto/github.com/HariPrasath-3/scheduler-service/proto/scheduler/v1"
	"github.com/HariPrasath-3/scheduler-service/internal/service"
	"github.com/HariPrasath-3/scheduler-service/internal/service/scheduler"
	"github.com/HariPrasath-3/scheduler-service/internal/validator"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
)

type SchedulerController struct {
	schedulev1.UnimplementedSchedulerServiceServer
	env              *env.Env
	schedulerService *scheduler.SchedulerService
	validator        *validator.SchedulerValidator
}

func NewSchedulerController(env *env.Env) *SchedulerController {
	return &SchedulerController{
		env:              env,
		schedulerService: service.GetServiceFactory(env).SchedulerService,
		validator:        validator.NewSchedulerValidator(),
	}
}

func (c *SchedulerController) Schedule(
	ctx context.Context,
	req *schedulev1.ScheduleRequest,
) (*schedulev1.ScheduleResponse, error) {
	if err := c.validator.ValidateScheduleRequest(req); err != nil {
		return &schedulev1.ScheduleResponse{
			Status: &common.Status{
				Code:    common.StatusCode_BAD_REQUEST,
				Message: err.Error(),
			},
		}, nil
	}

	err := c.schedulerService.Schedule(ctx, req)
	if err != nil {
		return &schedulev1.ScheduleResponse{
			Status: &common.Status{
				Code:    common.StatusCode_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	return &schedulev1.ScheduleResponse{
		Status: &common.Status{
			Code: common.StatusCode_SUCCESS,
		},
	}, nil
}

func (c *SchedulerController) Cancel(
	ctx context.Context,
	req *schedulev1.CancelRequest,
) (*schedulev1.CancelResponse, error) {
	if err := c.validator.ValidateCancelRequest(req); err != nil {
		return &schedulev1.CancelResponse{
			Status: &common.Status{
				Code:    common.StatusCode_BAD_REQUEST,
				Message: err.Error(),
			},
		}, nil
	}

	err := c.schedulerService.Cancel(ctx, req)
	if err != nil {
		return &schedulev1.CancelResponse{
			Status: &common.Status{
				Code:    common.StatusCode_INTERNAL_SERVER_ERROR,
				Message: err.Error(),
			},
		}, nil
	}

	return &schedulev1.CancelResponse{
		Status: &common.Status{
			Code: common.StatusCode_SUCCESS,
		},
	}, nil
}
