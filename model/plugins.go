// Copyright (c) 2017 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import "net/http"

type RouteDefiner interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request))
}
