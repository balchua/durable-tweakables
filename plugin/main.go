package main

import (
	"context"
	"math/rand"
	"time"

	shared "github.com/balchua/durable-tweakables/lib"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type MessageSource struct {
	appLogger hclog.Logger
	stopCh    chan bool
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (m *MessageSource) randStringBytes(ctx context.Context, n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

// Run runs the GRPC server and streams records from the plugin to the host application (client).
func (m *MessageSource) Run(ctx context.Context) (<-chan interface{}, error) {
	m.stopCh = make(chan bool)
	recordsCh := make(chan interface{})
	go func() {
		defer func() {
			m.appLogger.Info("closing records channel")
			close(recordsCh)
		}()
		// Fetch records and send them through the channel
		for {
			select {
			case <-m.stopCh:
				// External event triggered, stop the loop

				return
			default:
				b := m.randStringBytes(ctx, 10)
				recordsCh <- b
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return recordsCh, nil

}

func (m *MessageSource) Configure(ctx context.Context, config map[string]string) (*shared.ConfigureResponse, error) {
	m.appLogger.Info("MessageSource.Configure called with config", "config", config)
	return &shared.ConfigureResponse{Status: "OK"}, nil
}

func (m *MessageSource) Start(ctx context.Context, req *shared.StartRequest) (*shared.StartResponse, error) {
	m.appLogger.Info("MessageSource.Start called")
	return &shared.StartResponse{Status: "OK"}, nil
}

func (m *MessageSource) Stop(ctx context.Context, req *shared.StopRequest) (*shared.StopResponse, error) {
	m.appLogger.Info("MessageSource.Stop called")
	m.stopCh <- true
	return &shared.StopResponse{Status: "OK"}, nil
}

func (m *MessageSource) Teardown(ctx context.Context, req *shared.TeardownRequest) (*shared.TeardownResponse, error) {
	m.appLogger.Info("MessageSource.Teardown called")
	return &shared.TeardownResponse{Status: "OK"}, nil
}

func main() {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:       "my-plugin",
		JSONFormat: true,
		Level:      hclog.LevelFromString("DEBUG"),
	})

	messageSource := &MessageSource{
		appLogger: appLogger,
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"msg_src_grpc": &shared.MessageSourcePlugin{Impl: messageSource},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})

}
