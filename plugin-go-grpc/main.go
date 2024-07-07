// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"mrtdeh/plugin/shared"
	"time"

	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type PluginWrapper struct {
	pm shared.Maintainer
}

func (k *PluginWrapper) Connect(a shared.Maintainer) error {

	r, err := a.GetInfo()
	if err != nil {
		return err
	}

	fmt.Println("get info from main : ", r)
	k.pm = a
	return nil
}

func (k *PluginWrapper) Ping() string {
	return "pong"
}

func main() {
	var a = &PluginWrapper{}

	go func() {
		for {
			time.Sleep(time.Second * 1)
			if a.pm != nil {
				_, err := a.pm.GetInfo()
				if err != nil {
					log.Fatal("error in getInfo : ", err)
				}
			}
		}
	}()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"sample-plugin": &shared.DefaultPlugin{Impl: a},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
