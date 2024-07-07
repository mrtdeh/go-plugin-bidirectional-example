// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package shared

import (
	"context"

	"mrtdeh/plugin/proto"

	plugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.PluginClient
}

// Main -> plugin.Connect request
func (m *GRPCClient) Connect(a Maintainer) error {
	addHelperServer := &GRPCMaintainerServer{Impl: a}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterMaintainerServer(s, addHelperServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	// Call plugin gRPC.Connect and pass broker id for connect back
	_, err := m.client.Connect(context.Background(), &proto.ConnectRequest{
		ServerId: brokerID,
	})

	// s.Stop()
	return err
}

// Main -> plugin.Ping request
func (m *GRPCClient) Ping() string {
	_, err := m.client.Ping(context.Background(), &proto.Empty{})
	if err != nil {
		return ""
	}

	return "poong"
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCPluginServer struct {
	// This is the real implementation
	Impl Plugin

	broker *plugin.GRPCBroker
}

// Plugin -> gRPC.Connect body
func (m *GRPCPluginServer) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.Empty, error) {
	conn, err := m.broker.Dial(req.ServerId)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()

	a := &GRPCMaintainerClient{proto.NewMaintainerClient(conn)}
	return &proto.Empty{}, m.Impl.Connect(a)
}

// Plugin -> gRPC.Ping body
func (m *GRPCPluginServer) Ping(ctx context.Context, req *proto.Empty) (*proto.PongResponse, error) {
	return &proto.PongResponse{}, nil
}

// Cover struct for Main gRPC-methods
type GRPCMaintainerClient struct{ client proto.MaintainerClient }

// Implement Cover method GetInfo for main.GetInfo
func (m *GRPCMaintainerClient) GetInfo() (string, error) {
	resp, err := m.client.GetInfo(context.Background(), &proto.Empty{})
	if err != nil {
		return "", err
	}
	return resp.Info, err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCMaintainerServer struct {
	// This is the real implementation
	Impl Maintainer
}

func (m *GRPCMaintainerServer) GetInfo(ctx context.Context, req *proto.Empty) (resp *proto.GetInfoResponse, err error) {
	r, err := m.Impl.GetInfo()
	if err != nil {
		return nil, err
	}
	return &proto.GetInfoResponse{Info: r}, err
}
