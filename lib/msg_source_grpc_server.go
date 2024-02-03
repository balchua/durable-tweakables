package lib

import (
	"context"

	"github.com/balchua/durable-tweakables/lib/proto/msg_source"
)

// Here is the gRPC server that GRPCClient talks to.
// This is mostly implemented by the plugin.
type GRPCServer struct {
	// This is the real implementation
	Impl SourcePlugin
	msg_source.UnimplementedSourcePluginServer
}

// Run runs the GRPC server and streams records from the implementation to the client.
// It takes a request and a stream as input parameters.
// It returns an error if there is any issue with running the server or streaming the records.
func (m *GRPCServer) Run(
	req *msg_source.Source_Run_Request, stream msg_source.SourcePlugin_RunServer) error {
	recordsCh, err := m.Impl.Run(stream.Context())

	if err != nil {
		return err
	}

	for record := range recordsCh {
		err := stream.Send(&msg_source.Source_Run_Response{
			Value: record.([]byte)})
		if err != nil {
			return err
		}
	}

	return nil

}

func (m *GRPCServer) Configure(ctx context.Context, req *msg_source.Source_Configure_Request) (*msg_source.Source_Configure_Response, error) {
	resp, err := m.Impl.Configure(ctx, req.Config)
	if err != nil {
		return nil, err
	}
	return &msg_source.Source_Configure_Response{
		Status: resp.Status,
	}, nil
}

func (m *GRPCServer) Start(ctx context.Context, req *msg_source.Source_Start_Request) (*msg_source.Source_Start_Response, error) {
	startReq := &StartRequest{}
	resp, err := m.Impl.Start(ctx, startReq)
	if err != nil {
		return nil, err
	}
	return &msg_source.Source_Start_Response{
		Status: resp.Status,
	}, nil
}

func (m *GRPCServer) Stop(ctx context.Context, req *msg_source.Source_Stop_Request) (*msg_source.Source_Stop_Response, error) {
	stopRequest := &StopRequest{}

	resp, err := m.Impl.Stop(ctx, stopRequest)
	if err != nil {
		return nil, err
	}
	return &msg_source.Source_Stop_Response{
		Status: resp.Status,
	}, nil
}

func (m *GRPCServer) Teardown(ctx context.Context, req *msg_source.Source_Teardown_Request) (*msg_source.Source_Teardown_Response, error) {
	teardownReq := &TeardownRequest{}
	resp, err := m.Impl.Teardown(ctx, teardownReq)
	if err != nil {
		return nil, err
	}
	return &msg_source.Source_Teardown_Response{
		Status: resp.Status,
	}, nil
}
