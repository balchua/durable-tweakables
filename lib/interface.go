package lib

import (
	"context"

	"github.com/balchua/durable-tweakables/lib/proto/msg_source"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "TWEABLES_PLUGIN",
	MagicCookieValue: "tweak-me",
}

type ConfigureRequest struct{}

type ConfigureResponse struct {
	Status string
}

type StartRequest struct{}

type StartResponse struct {
	Status string
}

type StopRequest struct{}

type StopResponse struct {
	Status string
}

type TeardownRequest struct{}

type TeardownResponse struct {
	Status string
}

// SourcePlugin is the interface that we're exposing as a plugin.
type SourcePlugin interface {
	Configure(context.Context, map[string]string) (*ConfigureResponse, error)
	Start(context.Context, *StartRequest) (*StartResponse, error)
	Run(context.Context) (<-chan interface{}, error)
	Stop(context.Context, *StopRequest) (*StopResponse, error)
	Teardown(context.Context, *TeardownRequest) (*TeardownResponse, error)
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type MessageSourcePlugin struct {
	// callback to receive the messages from the source
	Callback func([]byte) error
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl SourcePlugin
}

func (p *MessageSourcePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	msg_source.RegisterSourcePluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *MessageSourcePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		callback: p.Callback,
		client:   msg_source.NewSourcePluginClient(c)}, nil
}
