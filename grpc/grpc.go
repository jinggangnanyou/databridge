package grpc

import (
	"context"
	"fmt"

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
	return g.Client, nil
}
