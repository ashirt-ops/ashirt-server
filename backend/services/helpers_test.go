package services_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"
)

func TestSanitizeSlug(t *testing.T) {
	require.Equal(t, services.SanitizeSlug("?One?Two?Three?"), "one-two-three")
	require.Equal(t, services.SanitizeSlug("Harry"), "harry")
	require.Equal(t, services.SanitizeSlug("Harry Potter"), "harry-potter")
	require.Equal(t, services.SanitizeSlug("fancy_name"), "fancy-name")
	require.Equal(t, services.SanitizeSlug("Lots_Of-Fancy! Characters"), "lots-of-fancy-characters")
	require.Equal(t, services.SanitizeSlug("$$prefixed_and_postfixed$$"), "prefixed-and-postfixed")
}
