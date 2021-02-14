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
		GetRecord,
		http.MethodGet,
		"/records/:id",
		[]server.ResponseSetting{},
	},
	Route{
		FindRecords,
		http.MethodGet,
		"/records",
		[]server.ResponseSetting{},
	},
}
