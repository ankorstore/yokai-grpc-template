package service_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/ankorstore/yokai-grpc-template/internal"
	"github.com/ankorstore/yokai-grpc-template/proto"
	"github.com/ankorstore/yokai/grpcserver/grpcservertest"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestExampleUnary(t *testing.T) {
	var connFactory grpcservertest.TestBufconnConnectionFactory
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&connFactory, &logBuffer, &traceExporter))

	// conn preparation
	conn, err := connFactory.Create(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	defer func() {
		err = conn.Close()
		assert.NoError(t, err)
	}()

	// client preparation
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

func TestExampleStreaming(t *testing.T) {
	var connFactory grpcservertest.TestBufconnConnectionFactory
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter

	internal.RunTest(t, fx.Populate(&connFactory, &logBuffer, &traceExporter))

	// conn preparation
	conn, err := connFactory.Create(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err)

	defer func() {
		err = conn.Close()
		assert.NoError(t, err)
	}()

	// client preparation
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
