package main

import (
	"fmt"
	"os"
	"os/exec"

	shared "github.com/balchua/durable-tweakables/lib"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func run() error {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:       "host-app",
		JSONFormat: true,
	})

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("/home/thor/workspace/durable-tweakables/plugin/sample-plugin"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
		Logger: appLogger,
	})
	defer client.Kill()

	var (
		rpcClient plugin.ClientProtocol
		err       error
	)
	// Connect via RPC
	if rpcClient, err = client.Client(); err != nil {
		panic(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("msg_src_grpc")
	if err != nil {
		return err
	}

	msg := raw.(shared.MessageSource)
	if x, err := msg.Receive(1); err == nil {
		appLogger.Info("message received", "message", string(x))
	}

	return nil
}

func main() {

	if err := run(); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
