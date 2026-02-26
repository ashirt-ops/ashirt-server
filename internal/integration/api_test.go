package integration_test

import (
	"net/http"
	"testing"

	"github.com/ashirt-ops/ashirt-server/internal/integration"
)

func TestAPI(t *testing.T) {
	sampleEvidenceBody := map[string]string{"notes": "evi1", "occurred_at": "1564441965400000000", "operator_id": "1"}

	t.Run("Making API requests without API keys", func(t *testing.T) {
		a := integration.NewTester(t)

		a.Get("/api/operations").Do().ExpectUnauthorized()
		a.Post("/api/operations/1/evidence").WithMultipartBody(sampleEvidenceBody, nil).Do().ExpectUnauthorized()
	})

	t.Run("Listing Operations from API", func(t *testing.T) {
		a := integration.NewTester(t)
		alice := a.NewUser("adefaultuser", "Alice", "DefaultUser")
		bob := a.NewUser("bsnooper", "Bob", "Snooper")
		aliceKey := a.APIKeyForUser(alice)
		bobKey := a.APIKeyForUser(bob)

		a.Post("/web/operations").WithJSONBody(`{"name": "Alice Op", "slug": "alice-op"}`).AsUser(alice).Do().ExpectSuccess()
		a.Post("/web/operations").WithJSONBody(`{"name": "Bob Op", "slug": "bob-op"}`).AsUser(bob).Do().ExpectSuccess()

		a.Get("/api/operations").WithAPIKey(aliceKey).Do().ExpectSubsetJSONArray([]string{`{"slug": "alice-op"}`})
		a.Get("/api/operations").WithAPIKey(bobKey).Do().ExpectSubsetJSONArray([]string{`{"slug": "bob-op"}`})
	})

	t.Run("Creating Evidence from API", func(t *testing.T) {
		a := integration.NewTester(t)
		alice := a.NewUser("adefaultuser", "Alice", "DefaultUser")
		aliceKey := a.APIKeyForUser(alice)

		a.Post("/web/operations").WithJSONBody(`{"name": "Op 1", "slug": "op1"}`).AsUser(alice).Do().ExpectSuccess()
		evidenceUUID1 := a.Post("/api/operations/op1/evidence").WithAPIKey(aliceKey).WithMultipartBody(sampleEvidenceBody, nil).Do().ExpectStatus(http.StatusCreated).ResponseUUID()
		a.Get("/web/operations/op1/evidence").AsUser(alice).Do().ExpectSubsetJSONArray([]string{
			`{"uuid": "` + evidenceUUID1 + `", "description": "evi1", "occurredAt": "2019-07-29T23:12:45Z"}`,
		})
	})
}
