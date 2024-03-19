package internal

import (
	"github.com/ankorstore/yokai-grpc-template/internal/service"
	"github.com/ankorstore/yokai-grpc-template/proto"
	"github.com/ankorstore/yokai/fxgrpcserver"
	"go.uber.org/fx"
)

// ProvideServices is used to register the application services.
func ProvideServices() fx.Option {
	return fx.Options(
		// gRPC server service
		fxgrpcserver.AsGrpcServerService(service.NewExampleService, &proto.ExampleService_ServiceDesc),
	)
}
