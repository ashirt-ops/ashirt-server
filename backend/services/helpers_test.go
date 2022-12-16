// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestSanitizeSlug(t *testing.T) {
	require.Equal(t, services.SanitizeSlug("?One?Two?Three?"), "one-two-three")
	require.Equal(t, services.SanitizeSlug("Harry"), "harry")
	require.Equal(t, services.SanitizeSlug("Harry Potter"), "harry-potter")
	require.Equal(t, services.SanitizeSlug("fancy_name"), "fancy-name")
	require.Equal(t, services.SanitizeSlug("Lots_Of-Fancy! Characters"), "lots-of-fancy-characters")
}
