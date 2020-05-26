// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/helpers"
)

func TestClamp(t *testing.T) {
	require.Equal(t, int64(10), helpers.Clamp(12, 1, 10))
	require.Equal(t, int64(12), helpers.Clamp(12, 10, 20))
	require.Equal(t, int64(20), helpers.Clamp(12, 20, 30))
}
