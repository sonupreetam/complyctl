// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hashicorp/go-hclog"

	"github.com/complytime/complytime/cmd/openscap-plugin/server"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

var logger hclog.Logger

func init() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:  "openscap-plugin",
		Level: hclog.Debug,
	})
	hclog.SetDefault(logger)
}

func main() {
	hclog.Default().Info("Starting OpenSCAP plugin")
	openSCAPPlugin := server.New()
	pluginByType := map[string]hplugin.Plugin{
		plugin.PVPPluginName: &plugin.PVPPlugin{Impl: openSCAPPlugin},
	}
	plugin.Register(pluginByType)
}
