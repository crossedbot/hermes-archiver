package controller

import (
	"net/http"

	"github.com/crossedbot/common/golang/server"
)

// Route represents an Replayer server route
type Route struct {
	Handler          server.Handler
	Method           string
	Path             string
	ResponseSettings []server.ResponseSetting
}

// Routes is a list of Replayer server routes
var Routes = []Route{
	Route{
		Replay,
		http.MethodGet,
		"/replays/:id",
		[]server.ResponseSetting{},
	},
}
