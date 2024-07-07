// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"mrtdeh/plugin/shared"

	"github.com/hashicorp/go-plugin"
)

type pluginManager struct{}

func (m *pluginManager) GetInfo() (string, error) {
	fmt.Println("get info from plugin")
	return "xyz123", nil
}

func main() {
	// We don't want to see the plugin logs.
	log.SetOutput(ioutil.Discard)

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv("SAMPLE_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer client.Kill()
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	// Request the plugin
	raw, err := rpcClient.Dispense("sample-plugin")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	// We should have a Counter store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	samplePlugin := raw.(shared.Plugin)
	err = samplePlugin.Connect(&pluginManager{})
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	for {
		res := samplePlugin.Ping()
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}
		fmt.Println(res)
		time.Sleep(time.Second * 3)
	}
}
