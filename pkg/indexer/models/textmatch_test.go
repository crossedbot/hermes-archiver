package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchString(t *testing.T) {
	expected := "exact"
	actual := TextMatchExact.String()
	require.Equal(t, expected, actual)
}

func TestToTextMatch(t *testing.T) {
	expected := TextMatchExact
	actual, err := ToTextMatch("Exact")
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}
