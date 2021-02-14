package controller

import (
	"net/http"

	"github.com/crossedbot/common/golang/server"
)

type Route struct {
	Handler          server.Handler
	Method           string
	Path             string
	ResponseSettings []server.ResponseSetting
}

var Routes = []Route{
	Route{
		Replay,
		http.MethodGet,
		"/replays/:id",
		[]server.ResponseSetting{},
	},
}
