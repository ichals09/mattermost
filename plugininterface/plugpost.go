// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package plugininterface

import (
	"plugin"
)

var modifyPostTextFuncs []func(string) string
var handlePostedEventFuncs []func(string, string)

func registerPostPlugin(plugin *plugin.Plugin) {
	if sym, err := plugin.Lookup("ModifyPostText"); err == nil {
		modifyPostTextFuncs = append(modifyPostTextFuncs, sym.(func(string) string))
	}
	if sym, err := plugin.Lookup("HandlePostedEvent"); err == nil {
		handlePostedEventFuncs = append(handlePostedEventFuncs, sym.(func(string, string)))
	}
}

func ModifyPostText(posttext string) string {
	var output string = posttext
	for _, f := range modifyPostTextFuncs {
		output = f(output)
	}

	return output
}

func HandlePostedEvent(id string, message string) {
	for _, f := range handlePostedEventFuncs {
		f(id, message)
	}
}
