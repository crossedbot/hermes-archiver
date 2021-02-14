package models

import (
	"fmt"
)

const (
	// Error Codes
	ErrMaxRecordLimitCode = iota + 1000
	ErrFailedConversionCode
	ErrUnknownRecordTypeStringCode
	ErrRequiredParamCode
	ErrUnauthorizedCode
	ErrProcessingRequestCode
	ErrNotFoundCode
)

type Error struct {
	Code    int
	Message string
}

func (err Error) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}
