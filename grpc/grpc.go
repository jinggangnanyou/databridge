package grpc

import (
	"context"
	"fmt"

	"databridge/common"

	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcConn struct {
	Client *grpc.ClientConn
}

// GrpcPool GrpcPool
var GrpcPool *GrpcConn

// InitGrpc InitGrpc
func InitGrpc(grpcAddr string) error {
	tracer := otel.Tracer(common.ModuleName)
	_, span := tracer.Start(context.Background(), "init grpc")
	fmt.Printf("trace_id:%s,span_id:%s\n",
		span.SpanContext().TraceID(), span.SpanContext().SpanID())
	defer span.End()
	maxSizeOption := grpc.MaxCallRecvMsgSize(1024 * 1024 * 1024 * 1024)
	pool, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(
		insecure.NewCredentials()), grpc.WithDefaultCallOptions(maxSizeOption))
	if err != nil || pool == nil {
		return fmt.Errorf("grpc pool init fail, %v", err.Error())
	}
	GrpcPool = &GrpcConn{
		Client: pool,
	}
	return nil
}

func (g *GrpcConn) Get(c context.Context) (*grpc.ClientConn, error) {
	tracer := otel.Tracer("code-go-api")
	_, span := tracer.Start(c, "get grpc client")
	fmt.Printf("trace_id:%s,span_id:%s\n",
		span.SpanContext().TraceID(), span.SpanContext().SpanID())
	defer span.End()
	return g.Client, nil
}
