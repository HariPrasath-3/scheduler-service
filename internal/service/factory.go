package service

import (
	"sync"

	"github.com/HariPrasath-3/scheduler-service/internal/service/scheduler"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
)

var serviceFactory Factory

var mutexForFactory sync.Once

type Factory struct {
	SchedulerService *scheduler.SchedulerService
}

func GetServiceFactory(env *env.Env) *Factory {
	mutexForFactory.Do(func() {
		serviceFactory = Factory{
			SchedulerService: scheduler.NewSchedulerService(env),
		}
	})
	return &serviceFactory
}
