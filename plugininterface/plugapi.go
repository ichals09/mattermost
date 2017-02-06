// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package plugininterface

import (
	"net/http"

	l4g "github.com/alecthomas/log4go"
	"github.com/gorilla/mux"
	"github.com/mattermost/platform/model"
)

type GorillaMuxRouteDefiner struct {
	router *mux.Router
}

func (me *GorillaMuxRouteDefiner) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	me.router.HandleFunc(path, f)
}

func NewGorillaMuxRouteDefiner(router *mux.Router) model.RouteDefiner {
	return &GorillaMuxRouteDefiner{
		router: router,
	}
}

func RegisterPluginRoutes(rootRouter *mux.Router) {
	for _, pluginInfo := range pluginInfos {
		if sym, err := pluginInfo.GoPlugin.Lookup("RegisterRoutes"); err == nil {
			if len(pluginInfo.RouteName) <= 0 {
				l4g.Error("Plugin " + pluginInfo.DisplayName + " has register routes function but no RouteName!")
				continue
			}
			pluginRouter := rootRouter.PathPrefix("/" + pluginInfo.RouteName).Subrouter()
			sym.(func(model.RouteDefiner))(NewGorillaMuxRouteDefiner(pluginRouter))
		}
	}
}
