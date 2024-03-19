package service_test

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/ankorstore/yokai-grpc-template/internal"
	"github.com/ankorstore/yokai-grpc-template/proto"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestExampleUnary(t *testing.T) {
	var grpcServer *grpc.Server
	var lis *bufconn.Listener
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&grpcServer, &lis, &logBuffer, &traceExporter))

	defer func() {
		err := lis.Close()
		assert.NoError(t, err)

		grpcServer.GracefulStop()
	}()

	// client preparation
	conn, err := prepareGrpcClientTestConnection(lis)
	assert.NoError(t, err)

	client := proto.NewExampleServiceClient(conn)

	// call
	response, err := client.ExampleUnary(context.Background(), &proto.ExampleRequest{
		Text: "test",
	})
	assert.NoError(t, err)

	// response assertions
	assert.Equal(t, "response from grpc-app: you sent test", response.Text)

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "received: test",
	})

	// traces assertions
	tracetest.AssertHasTraceSpan(t, traceExporter, "ExampleUnary")
}

func TestTransformAndSplitText(t *testing.T) {
	var grpcServer *grpc.Server
	var lis *bufconn.Listener
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&grpcServer, &lis, &logBuffer, &traceExporter))

	defer func() {
		err := lis.Close()
		assert.NoError(t, err)

		grpcServer.GracefulStop()
	}()

	// client preparation
	conn, err := prepareGrpcClientTestConnection(lis)
	assert.NoError(t, err)

	client := proto.NewExampleServiceClient(conn)

	// send
	stream, err := client.ExampleStreaming(context.Background())
	assert.NoError(t, err)

	for _, text := range []string{"this", "is", "a", "test"} {
		err = stream.Send(&proto.ExampleRequest{
			Text: text,
		})
		assert.NoError(t, err)
	}

	err = stream.CloseSend()
	assert.NoError(t, err)

	// receive
	var responses []*proto.ExampleResponse

	wait := make(chan struct{})

	go func() {
		for {
			resp, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			assert.NoError(t, err)

			responses = append(responses, resp)
		}

		close(wait)
	}()

	<-wait

	// responses assertions
	assert.Len(t, responses, 4)
	assert.Equal(t, "response from grpc-app: you sent this", responses[0].Text)
	assert.Equal(t, "response from grpc-app: you sent is", responses[1].Text)
	assert.Equal(t, "response from grpc-app: you sent a", responses[2].Text)
	assert.Equal(t, "response from grpc-app: you sent test", responses[3].Text)

	// logs assertions
	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "received: this",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "received: is",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "received: a",
	})

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":   "info",
		"message": "received: test",
	})

	// traces assertions
	tracetest.AssertHasTraceSpan(t, traceExporter, "ExampleStreaming")

	span, err := traceExporter.Span("ExampleStreaming")
	assert.NoError(t, err)

	assert.Equal(t, "received: this", span.Events[0].Name)
	assert.Equal(t, "received: is", span.Events[1].Name)
	assert.Equal(t, "received: a", span.Events[2].Name)
	assert.Equal(t, "received: test", span.Events[3].Name)
}

func prepareGrpcClientTestConnection(lis *bufconn.Listener) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		"",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
