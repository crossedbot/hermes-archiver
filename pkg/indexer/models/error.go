package models

import (
	"fmt"
)

const (
	// Error Codes
	ErrMaxRecordLimitCode = iota + 1000
	ErrFailedConversionCode
	ErrUnknownRecordTypeStringCode
	ErrUnknownTextMatchStringCode
	ErrRequiredParamCode
	ErrUnauthorizedCode
	ErrProcessingRequestCode
	ErrNotFoundCode
)

// Error represents an error code and message
type Error struct {
	Code    int
	Message string
}

// Error returns the error as a formatted string
func (err Error) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}
