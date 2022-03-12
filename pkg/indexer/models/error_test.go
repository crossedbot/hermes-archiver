package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	expected := fmt.Sprintf("%d: hello world", ErrRequiredParamCode)
	err := Error{ErrRequiredParamCode, "hello world"}
	actual := err.Error()
	require.Equal(t, expected, actual)
}
