package lib

import (
	"context"

	"github.com/balchua/durable-tweakables/proto/msg_source"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	client msg_source.MessageSourceClient
}

func (m *GRPCClient) Receive(batch int) error {
	_, err := m.client.Receive(context.Background(), &msg_source.GetRequest{
		Batch: int32(batch),
	})
	return err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl MessageSource
}

func (m *GRPCServer) Receive(
	ctx context.Context,
	req *msg_source.GetRequest) (*msg_source.GetResponse, error) {
	v, err := m.Impl.Receive(req.Batch)
	return &msg_source.GetResponse{Value: v}, err
}
