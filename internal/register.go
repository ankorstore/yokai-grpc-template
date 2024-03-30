package internal

import (
	"github.com/ankorstore/yokai-grpc-template/internal/service"
	"github.com/ankorstore/yokai-grpc-template/proto"
	"github.com/ankorstore/yokai/fxgrpcserver"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		fxgrpcserver.AsGrpcServerService(service.NewExampleService, &proto.ExampleService_ServiceDesc),
	)
}
