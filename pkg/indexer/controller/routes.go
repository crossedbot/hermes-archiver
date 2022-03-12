package controller

import (
	"net/http"

	"github.com/crossedbot/common/golang/server"
)

// Route represents an Indexer server route
type Route struct {
	Handler          server.Handler
	Method           string
	Path             string
	ResponseSettings []server.ResponseSetting
}

// Routes is a list of Indexer server routes
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
