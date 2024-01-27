package main

import (
	shared "github.com/balchua/durable-tweakables/lib"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type MessageSource struct {
	appLogger hclog.Logger
}

func (m *MessageSource) Receive(batch int32) ([]byte, error) {
	m.appLogger.Info("MessageSource.Receive called with batch", "batch", batch)
	b := []byte("hello")
	return b, nil
}

func main() {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:       "my-plugin",
		JSONFormat: true,
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
