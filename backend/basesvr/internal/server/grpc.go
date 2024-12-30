package server

import (
	v1 "github.com/IsaacIsRebooting/Metok/backend/basesvr/api/basesvr/v1"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(options ...Option) *grpc.Server {
	params := &Params{}
	for _, option := range options {
		option(params)
	}
	warmUp(params)
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			tracing.Server(),
		),
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterAccountServer(srv, initAccountApplication())
	return srv
}
