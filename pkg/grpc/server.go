package grpc

import (
	"context"

	appconfig "github.com/HariPrasath-3/scheduler-service/pkg/config"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	"google.golang.org/grpc"
)

func NewGrpcServer(
	ctx context.Context,
	cfg *appconfig.GrpcConfig,
	environment *env.Env,
) (*grpc.Server, error) {

	serverOpts := []grpc.ServerOption{
		// grpc.UnaryInterceptor(env.UnaryServerInterceptor(environment)),
	}

	grpcServer := grpc.NewServer(serverOpts...)

	// ctx is not directly used by grpc.Server,
	// but passed for future lifecycle management.
	_ = ctx
	_ = cfg

	return grpcServer, nil
}
