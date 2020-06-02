package network

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetBaseURL(t *testing.T) {
	require.False(t, BaseURLSet())
	SetBaseURL("Something")
	require.Equal(t, "Something/api", apiURL)
	require.True(t, BaseURLSet())
}
