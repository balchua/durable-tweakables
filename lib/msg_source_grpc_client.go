package lib

import (
	"context"
	"io"
	"log"

	"github.com/balchua/durable-tweakables/lib/proto/msg_source"
)

// GRPCClient is an implementation of KV that talks over RPC.
// This is used by the host application to talk to the plugin.
type GRPCClient struct {
	callback func([]byte) error
	client   msg_source.SourcePluginClient
}

func (m *GRPCClient) Configure(ctx context.Context, config map[string]string) (*ConfigureResponse, error) {
	cfRequest := &msg_source.Source_Configure_Request{}
	cfRequest.Config = config
	resp, err := m.client.Configure(ctx, cfRequest)

	if err != nil {
		return nil, err
	}
	resp.Status = "OK"

	configResponse := &ConfigureResponse{
		Status: resp.Status,
	}

	return configResponse, err
}

// Run runs the GRPC client and streams records from the plugin to the host application (client).
func (m *GRPCClient) Run(ctx context.Context) (<-chan interface{}, error) {

	// currently not used, it is just a placeholder
	src := &msg_source.Source_Run_Request{}
	stream, err := m.client.Run(context.Background(), src)
	stopChan := make(chan interface{})

	go func() {
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				close(stopChan)
				break
			}
			if err != nil {
				log.Printf("%v.ListFeatures(_) = _, %v", m.client, err)
			}
			m.callback(data.Value)
		}
	}()
	<-stopChan
	return nil, err
}
func (m *GRPCClient) Start(ctx context.Context, startRequest *StartRequest) (*StartResponse, error) {

	req := &msg_source.Source_Start_Request{}
	resp, err := m.client.Start(ctx, req)

	if err != nil {
		return nil, err
	}
	resp.Status = "OK"

	startResponse := &StartResponse{
		Status: resp.Status,
	}

	return startResponse, nil

}

func (m *GRPCClient) Teardown(ctx context.Context, teardownRequest *TeardownRequest) (*TeardownResponse, error) {

	req := &msg_source.Source_Teardown_Request{}
	resp, err := m.client.Teardown(ctx, req)

	if err != nil {
		return nil, err
	}
	resp.Status = "OK"

	teardownResponse := &TeardownResponse{
		Status: resp.Status,
	}

	return teardownResponse, nil
}

func (m *GRPCClient) Stop(ctx context.Context, stopRequest *StopRequest) (*StopResponse, error) {
	req := &msg_source.Source_Stop_Request{}
	resp, err := m.client.Stop(ctx, req)

	if err != nil {
		return nil, err
	}
	resp.Status = "OK"

	stopResponse := &StopResponse{
		Status: resp.Status,
	}

	return stopResponse, nil
}
