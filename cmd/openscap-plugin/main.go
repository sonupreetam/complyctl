// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/complytime/complytime/cmd/openscap-plugin/server"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

func main() {
	openSCAPPlugin := server.New()
	pluginByType := map[string]hplugin.Plugin{
		plugin.PVPPluginName: &plugin.PVPPlugin{Impl: openSCAPPlugin},
	}
	plugin.Register(pluginByType)
}
