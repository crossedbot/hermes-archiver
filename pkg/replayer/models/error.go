package models

import (
	"fmt"
)

const (
	// Error Codes
	ErrRequiredParamCode = iota + 1000
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
