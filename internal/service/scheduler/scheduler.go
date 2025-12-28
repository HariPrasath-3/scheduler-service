package scheduler

import (
	"context"
	"log"

	schedulev1 "github.com/HariPrasath-3/scheduler-service/client/golang/proto/github.com/HariPrasath-3/scheduler-service/proto/scheduler/v1"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
)

type SchedulerService struct {
	env *env.Env
}

func NewSchedulerService(env *env.Env) *SchedulerService {
	return &SchedulerService{env: env}
}

func (s *SchedulerService) Schedule(
	ctx context.Context,
	req *schedulev1.ScheduleRequest,
) error {
	log.Printf("Received Schedule request: %+v", req)
	// TODO: implement repository.Save()

	return nil
}
