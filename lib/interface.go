package lib

import (
	"context"

	"github.com/balchua/durable-tweakables/proto/msg_source"
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

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"msg_src_grpc": &MessageSourcePlugin{},
}

type MessageSource interface {
	Receive(batch int32) ([]byte, error)
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type MessageSourcePlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl MessageSource
}

func (p *MessageSourcePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	msg_source.RegisterMessageSourceServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *MessageSourcePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: msg_source.NewMessageSourceClient(c)}, nil
}
