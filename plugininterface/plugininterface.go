// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package plugininterface

import (
	"fmt"
	"plugin"

	l4g "github.com/alecthomas/log4go"
)

const (
	PLUGIN_DIR = "plugins"
)

type PluginInfo struct {
	DisplayName string
	File        string
	RouteName   string
	GoPlugin    *plugin.Plugin
}

var pluginInfos []*PluginInfo

func getPluginInfo(infoin map[string]string) *PluginInfo {
	var infoout PluginInfo

	if res, ok := infoin["DisplayName"]; ok {
		infoout.DisplayName = res
	}

	if res, ok := infoin["RouteName"]; ok {
		infoout.RouteName = res
	}

	return &infoout
}

func getPluginFiles(plugdir string) []string {
	return []string{"plugins/mattermost-plugin-jira.so"}
}

func registerPlugin(plugin *plugin.Plugin) {
	registerPostPlugin(plugin)
}

func InitPlugins() {
	for _, pluginfile := range getPluginFiles(PLUGIN_DIR) {
		plug, err := plugin.Open(pluginfile)
		if err != nil {
			l4g.Error("Plugin Failed to load", err)
			continue
		}

		var pluginInfo *PluginInfo
		if sym, err := plug.Lookup("GetInfo"); err == nil {
			pluginInfo = getPluginInfo(sym.(func() map[string]string)())
		} else {
			pluginInfo = &PluginInfo{
				DisplayName: pluginfile,
			}
		}

		pluginInfo.File = pluginfile
		pluginInfo.GoPlugin = plug

		pluginInfos = append(pluginInfos, pluginInfo)

		l4g.Info("Loaded plugin: " + fmt.Sprint(pluginInfo))

		registerPlugin(plug)
	}

}
