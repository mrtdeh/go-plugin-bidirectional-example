// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log"
	"mrtdeh/plugin/shared"
	"time"

	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type AliReza struct {
	shared.PluginImplementation
}

func main() {
	var ali = &AliReza{}

	go func() {
		for {
			time.Sleep(time.Second * 1)
			if pm := ali.GetMaintainer(); pm != nil {
				_, err := pm.GetInfo()
				if err != nil {
					log.Fatal("error in getInfo : ", err)
				}
			}
		}
	}()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"sample-plugin": &shared.DefaultPlugin{Impl: ali},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
