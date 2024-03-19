package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ankorstore/yokai-grpc-template/proto"
	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
)

// ExampleService is the gRPC server service for ExampleService.
type ExampleService struct {
	proto.UnimplementedExampleServiceServer
	config *config.Config
}

// NewExampleService returns a new [TransformTextServiceService] instance.
func NewExampleService(cfg *config.Config) *ExampleService {
	return &ExampleService{
		config: cfg,
	}
}

// ExampleUnary returns the text provided in the [proto.ExampleRequest] by adding the application name.
func (s *ExampleService) ExampleUnary(ctx context.Context, in *proto.ExampleRequest) (*proto.ExampleResponse, error) {
	ctx, span := trace.CtxTracerProvider(ctx).Tracer("ExampleService").Start(ctx, "ExampleUnary")
	defer span.End()

	log.CtxLogger(ctx).Info().Msgf("received: %s", in.Text)

	return &proto.ExampleResponse{
		Text: fmt.Sprintf("response from %s: you sent %s", s.config.AppName(), in.Text),
	}, nil
}

// ExampleStreaming streams the text provided in the [proto.ExampleRequest] by adding the application name.
func (s *ExampleService) ExampleStreaming(stream proto.ExampleService_ExampleStreamingServer) error {
	ctx := stream.Context()

	ctx, span := trace.CtxTracerProvider(ctx).Tracer("ExampleService").Start(ctx, "ExampleStreaming")
	defer span.End()

	logger := log.CtxLogger(ctx)

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			logger.Info().Msg("end of rpc")

			return nil
		}

		if err != nil {
			logger.Error().Err(err).Msgf("error while receiving: %v", err)
		}

		logger.Info().Msgf("received: %s", req.Text)

		span.AddEvent(fmt.Sprintf("received: %s", req.Text))

		err = stream.Send(&proto.ExampleResponse{
			Text: fmt.Sprintf("response from %s: you sent %s", s.config.AppName(), req.Text),
		})

		if err != nil {
			logger.Error().Err(err).Msgf("error while sending: %v", err)

			return err
		}
	}
}
