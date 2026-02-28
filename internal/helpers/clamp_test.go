package helpers_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestClamp(t *testing.T) {
	require.Equal(t, int64(10), helpers.Clamp(12, 1, 10))
	require.Equal(t, int64(12), helpers.Clamp(12, 10, 20))
	require.Equal(t, int64(20), helpers.Clamp(12, 20, 30))
}
