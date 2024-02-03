package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	shared "github.com/balchua/durable-tweakables/lib"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type receiver struct {
	appLogger hclog.Logger
}

func (r *receiver) processMessage(b []byte) error {
	r.appLogger.Info("message received", "data", string(b))
	return nil
}

func run() error {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:       "host-app",
		JSONFormat: true,
		Level:      hclog.LevelFromString("DEBUG"),
	})
	r := &receiver{
		appLogger: appLogger,
	}

	// PluginMap is the map of plugins we can dispense.
	pluginMap := map[string]plugin.Plugin{
		"msg_src_grpc": &shared.MessageSourcePlugin{
			Callback: r.processMessage,
		},
	}

	//cmd := exec.Command("/bin/sh", "/home/thor/workspace/durable-tweakables/plugin/runDebug.sh")
	cmd := exec.Command("/home/thor/workspace/durable-tweakables/plugin/sample-plugin")
	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         pluginMap,
		Cmd:             cmd,
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

	ctx := context.Background()

	msg := raw.(shared.SourcePlugin)
	// Call the plugin's configure
	configurations := map[string]string{
		"config1": "value1",
		"config2": "value2",
	}
	if cfgResp, err := msg.Configure(ctx, configurations); err != nil {
		appLogger.Error("unable to configure the plugin", "error", err)
	} else {
		appLogger.Info("plugin configured", "status", cfgResp.Status)
	}

	go doSomething(appLogger, msg)

	ch := time.After(10 * time.Second)
	<-ch

	msg.Stop(ctx, &shared.StopRequest{})
	return nil
}

func doSomething(appLogger hclog.Logger, msg shared.SourcePlugin) {
	// we dont need to wait for the channel in this case.
	if _, err := msg.Run(context.Background()); err != nil {
		appLogger.Error("unable to start the plugin", "error", err)
	}
}

func main() {

	if err := run(); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
