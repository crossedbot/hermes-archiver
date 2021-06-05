package controller

import (
	"fmt"
	"net/http"

	"github.com/crossedbot/common/golang/server"
	"github.com/crossedbot/hermes-archiver/pkg/replayer/models"
)

// Replay handles requests to retrieve a WARC replay for the given record Id
func Replay(w http.ResponseWriter, r *http.Request, p server.Parameters) {
	id := p.Get("id")
	if id == "" {
		server.JsonResponse(
			w,
			models.Error{
				Code:    models.ErrRequiredParamCode,
				Message: "path parameter 'id' is required",
			},
			http.StatusBadRequest,
		)
		return
	}
	replay, err := V1().Replay(id)
	if err == ErrorReplayNotFound {
		server.JsonResponse(
			w,
			models.Error{
				Code:    models.ErrNotFoundCode,
				Message: fmt.Sprintf("failed to get replay; %s", err),
			},
			http.StatusNotFound,
		)
		return
	} else if err != nil {
		server.JsonResponse(
			w,
			models.Error{
				Code:    models.ErrProcessingRequestCode,
				Message: fmt.Sprintf("failed to get replay; %s", err),
			},
			http.StatusInternalServerError,
		)
		return
	}
	server.JsonResponse(w, &replay, http.StatusOK)
}
