// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"

	"google.golang.org/grpc"

	"mrtdeh/plugin/proto"

	"github.com/hashicorp/go-plugin"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"sample-plugin": &DefaultPlugin{},
}

type Maintainer interface {
	GetInfo() (string, error)
}

// KV is the interface that we're exposing as a plugin.
type Plugin interface {
	Connect(Maintainer) error
	Ping() string
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type DefaultPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Plugin
}

func (p *DefaultPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServer(s, &GRPCPluginServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *DefaultPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewPluginClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &DefaultPlugin{}
