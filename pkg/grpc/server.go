package grpc

import (
	"context"
	"log"
	"net"

	appconfig "github.com/HariPrasath-3/scheduler-service/pkg/config"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	"google.golang.org/grpc"
)

func NewGrpcServer(
	ctx context.Context,
	environment *env.Env,
) (*grpc.Server, error) {

	serverOpts := []grpc.ServerOption{
		// grpc.UnaryInterceptor(env.UnaryServerInterceptor(environment)),
	}

	grpcServer := grpc.NewServer(serverOpts...)
	return grpcServer, nil
}

func Serve(
	ctx context.Context,
	grpcServer *grpc.Server,
	cfg *appconfig.GrpcConfig,
) error {
	lis, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		return err
	}
	log.Printf("gRPC API listening on %s", cfg.Host)
	return grpcServer.Serve(lis)
}
