package integration_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/internal/integration"
)

func TestUserAuthorization(t *testing.T) {
	a := integration.NewTester(t)
	a.NewUser("admin.user", "Al", "Admin") // Admin's have special viewing access, which break these tests. These tests focus on two plain users.
	alice := a.NewUser("asmith", "Alice", "Smith")
	bob := a.NewUser("bsmith", "Bob", "Smith")

	// Create an operation as alice and ensure bob cannot see it
	a.Post("/web/operations").WithJSONBody(`{"name": "Alice's Operation", "slug": "alice"}`).AsUser(alice).Do().ExpectSuccess()
	a.Get("/web/operations").AsUser(bob).Do().ExpectJSON(`[]`)

	// Ensure bob cannot create findings/evidence under alice's operation
	a.Post("/web/operations/alice/findings").WithJSONBody(`{"title": "finding", "category": "Detection Gap", "description": ""}`).AsUser(bob).Do().ExpectUnauthorized()
	a.Post("/web/operations/alice/evidence").WithMultipartBody(map[string]string{"description": "evi"}, nil).AsUser(bob).Do().ExpectUnauthorized()

	// Create findings as bob and ensure alice can't see them
	a.Post("/web/operations").WithJSONBody(`{"name": "Bob's Operation", "slug": "bob"}`).AsUser(bob).Do().ExpectSubsetJSON(`{"name": "Bob's Operation"}`)
	findingUUID1 := a.Post("/web/operations/bob/findings").WithJSONBody(`{"title": "f1", "category": "Detection Gap", "description": ""}`).AsUser(bob).Do().ExpectSuccess().ResponseUUID()
	a.Post("/web/operations/bob/findings").WithJSONBody(`{"title": "e2", "category": "Detection Gap", "description": ""}`).AsUser(bob).Do().ExpectSuccess()
	a.Get("/web/operations/bob/findings").AsUser(alice).Do().ExpectNotFound()
	a.Get("/web/operations/bob/findings/" + findingUUID1).AsUser(alice).Do().ExpectNotFound()
	a.Get("/web/operations/alice/findings/" + findingUUID1).AsUser(alice).Do().ExpectNotFound()

	// Create evidence as bob and ensure alice can't see them
	evidenceUUID1 := a.Post("/web/operations/bob/evidence").WithMultipartBody(map[string]string{"description": "e1"}, nil).AsUser(bob).Do().ExpectSuccess().ResponseUUID()
	a.Post("/web/operations/bob/evidence").WithMultipartBody(map[string]string{"description": "e2"}, nil).AsUser(bob).Do().ExpectSuccess()
	a.Get("/web/operations/bob/evidence").AsUser(alice).Do().ExpectNotFound()
	a.Get("/web/operations/bob/evidence/" + evidenceUUID1).AsUser(alice).Do().ExpectNotFound()
	a.Get("/web/operations/alice/evidence/" + evidenceUUID1).AsUser(alice).Do().ExpectNotFound()

	// Ensure alice cannot add evidence to bob's findings
	evidenceUUID3 := a.Post("/web/operations/alice/evidence").
		WithMultipartBody(map[string]string{"description": "Alice's evidence"}, nil).
		AsUser(alice).Do().
		ExpectSubsetJSON(`{"description": "Alice's evidence"}`).
		ResponseUUID()

	a.Put("/web/operations/bob/findings/" + findingUUID1 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID1 + `"], "evidenceToRemove": []}`).AsUser(alice).Do().ExpectUnauthorized()
	a.Put("/web/operations/bob/findings/" + findingUUID1 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID3 + `"], "evidenceToRemove": []}`).AsUser(alice).Do().ExpectUnauthorized()
	a.Put("/web/operations/alice/findings/" + findingUUID1 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID1 + `"], "evidenceToRemove": []}`).AsUser(alice).Do().ExpectUnauthorized()
	a.Put("/web/operations/alice/findings/" + findingUUID1 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID3 + `"], "evidenceToRemove": []}`).AsUser(alice).Do().ExpectUnauthorized()

	// Ensure alice cannot add bob's evidence to alice's findings
	findingUUID3 := a.Post("/web/operations/alice/findings").
		WithJSONBody(`{"title": "Alice's finding", "category": "Detection Gap", "description": ""}`).
		AsUser(alice).Do().
		ExpectSubsetJSON(`{"title": "Alice's finding"}`).
		ResponseUUID()

	evidenceUUID4 := a.Post("/web/operations/bob/evidence").
		WithMultipartBody(map[string]string{"description": "Bob's evidence"}, nil).
		AsUser(bob).Do().
		ExpectSubsetJSON(`{"description": "Bob's evidence"}`).
		ResponseUUID()

	a.Put("/web/operations/bob/findings/" + findingUUID3 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID4 + `"], "evidenceToRemove": []}`).AsUser(alice).Do().ExpectUnauthorized()
	a.Put("/web/operations/alice/findings/" + findingUUID3 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID4 + `"], "evidenceToRemove": []}`).AsUser(alice).Do().ExpectUnauthorized()

	// Ensure alice cannot see bob's evidence images
	a.Get("/web/operations/alice/evidence/1/preview").AsUser(alice).Do().ExpectNotFound()
	a.Get("/web/operations/alice/evidence/1/media").AsUser(alice).Do().ExpectNotFound()

	a.Get("/web/operations").AsUser(alice).Do().ExpectSubsetJSONArray([]string{`{"slug": "alice", "name": "Alice's Operation", "numUsers": 1}`})
}
