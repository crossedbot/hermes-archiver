package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/crossedbot/common/golang/server"
	"github.com/crossedbot/hermes-archiver/pkg/indexer/models"
	"github.com/crossedbot/simplecdxj"
)

const (
	MaxRecordLimit = 1000
)

// FindRecords handles a request to find CDXJ records matching a given set of
// values
func FindRecords(w http.ResponseWriter, r *http.Request, p server.Parameters) {
	var err error
	query := r.URL.Query()

	// get surt parameter
	surt := query.Get("surt")
	matchStr := query.Get("match")
	match := models.TextMatchExact
	if matchStr != "" {
		match, err = models.ToTextMatch(matchStr)
		if err != nil {
			server.JsonResponse(w, models.Error{
				Code: models.ErrUnknownTextMatchStringCode,
				Message: fmt.Sprintf(
					"%s \"%s\"",
					err, matchStr,
				),
			}, http.StatusBadRequest)
			return
		}
	}

	// get record type parameters
	recordTypes := []simplecdxj.RecordType{}
	recordTypeStrings, present := query["type"]
	if present {
		recordTypeMap := make(map[simplecdxj.RecordType]struct{})
		for _, s := range recordTypeStrings {
			recordType, err := simplecdxj.ParseRecordType(s)
			if err != nil {
				server.JsonResponse(
					w,
					models.Error{
						Code: models.ErrUnknownRecordTypeStringCode,
						Message: fmt.Sprintf(
							"%s \"%s\"",
							err, s,
						),
					},
					http.StatusBadRequest,
				)
				return
			}
			recordTypeMap[recordType] = struct{}{}
		}
		for r, _ := range recordTypeMap {
			recordTypes = append(recordTypes, r)
		}
	}

	// get time range parameters
	before := query.Get("before")
	after := query.Get("after")

	// get limit parameter
	limit := 10
	if v := query.Get("limit"); v != "" {
		var err error
		limit, err = strconv.Atoi(v)
		if err != nil {
			server.JsonResponse(
				w,
				models.Error{
					Code:    models.ErrFailedConversionCode,
					Message: "limit is not an integer",
				},
				http.StatusBadRequest,
			)
			return
		}
		if limit > MaxRecordLimit {
			server.JsonResponse(
				w,
				models.Error{
					Code: models.ErrMaxRecordLimitCode,
					Message: fmt.Sprintf(
						"max limit exceeded [1 .. %d]",
						MaxRecordLimit,
					),
				},
				http.StatusBadRequest,
			)
			return
		}
	}

	records, err := V1().FindRecords(
		surt, recordTypes,
		before, after,
		match, limit,
	)
	if err != nil {
		server.JsonResponse(
			w,
			models.Error{
				Code: models.ErrProcessingRequestCode,
				Message: fmt.Sprintf(
					"failed to find records: %s",
					err,
				),
			},
			http.StatusInternalServerError,
		)
		return
	}
	server.JsonResponse(w, &records, http.StatusOK)
}

// GetRecord handles requests to retrieve a CDXJ record matching a given ID
func GetRecord(w http.ResponseWriter, r *http.Request, p server.Parameters) {
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
	record, err := V1().GetRecord(id)
	if err == ErrorRecordNotFound {
		server.JsonResponse(
			w,
			models.Error{
				Code: models.ErrNotFoundCode,
				Message: fmt.Sprintf(
					"failed to get record; %s",
					err,
				),
			},
			http.StatusNotFound,
		)
		return
	} else if err != nil {
		server.JsonResponse(
			w,
			models.Error{
				Code: models.ErrProcessingRequestCode,
				Message: fmt.Sprintf(
					"failed to get record; %s",
					err,
				),
			},
			http.StatusInternalServerError,
		)
		return
	}
	server.JsonResponse(w, &record, http.StatusOK)
}
