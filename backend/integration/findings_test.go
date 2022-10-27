// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package integration_test

import (
	"net/http"
	"testing"

	"github.com/theparanoids/ashirt-server/backend/integration"
)

func TestFindings(t *testing.T) {
	t.Run("Editing Findings", func(t *testing.T) {
		a := integration.NewTester(t)
		a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

		a.Post("/web/operations").WithJSONBody(`{"name": "Op 1", "slug": "op"}`).Do().ExpectSuccess()

		a.Post("/web/operations/op/tags").WithJSONBody(`{"name": "one", "colorName": "red"}`).Do().ExpectSuccess()
		a.Post("/web/operations/op/tags").WithJSONBody(`{"name": "two", "colorName": "green"}`).Do().ExpectSuccess()
		a.Post("/web/operations/op/tags").WithJSONBody(`{"name": "three", "colorName": "blue"}`).Do().ExpectSuccess()

		uuid := a.Post("/web/operations/op/findings").WithJSONBody(`{"title": "Finding 1", "category": "Product", "description": "Here is my finding"}`).Do().ExpectStatus(http.StatusCreated).ResponseUUID()
		a.Get("/web/operations/op/findings").Do().ExpectSubsetJSONArray([]string{`{"uuid": "` + uuid + `", "title": "Finding 1", "category": "Product", "description": "Here is my finding"}`})
		a.Put("/web/operations/op/findings/" + uuid).WithJSONBody(`{
			"title": "Updated title", 
			"category": "Network", 
			"description": "Updated description",
			"readyToReport": false,
			"ticketLink": null
		}`).Do().ExpectSuccess()
		a.Get("/web/operations/op/findings").Do().ExpectSubsetJSONArray([]string{`{"uuid": "` + uuid + `", "title": "Updated title", "category": "Network", "description": "Updated description"}`})
	})

	t.Run("Deleting findings", func(t *testing.T) {
		a := integration.NewTester(t)
		a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

		a.Post("/web/operations").WithJSONBody(`{"name": "Op 1", "slug": "op"}`).Do().ExpectSuccess()
		uuid := a.Post("/web/operations/op/findings").WithJSONBody(`{"title": "To be deleted...", "category": "Enterprise", "description": ""}`).Do().ExpectStatus(http.StatusCreated).ResponseUUID()
		a.Get("/web/operations/op/findings").Do().ExpectSubsetJSONArray([]string{`{"uuid": "` + uuid + `", "title": "To be deleted..."}`})
		a.Delete("/web/operations/op/findings/" + uuid).Do().ExpectSuccess()
		a.Get("/web/operations/op/findings").Do().ExpectJSON("[]")
	})

	t.Run("Ensure users cannot edit findings of an operation they do not have write access to", func(t *testing.T) {
		a := integration.NewTester(t)
		alice := a.NewUser("aowner", "Alice", "Owner")
		bob := a.NewUser("battacker", "Bob", "Attacker")

		a.Post("/web/operations").WithJSONBody(`{"name": "Alice's Operation", "slug": "alice"}`).AsUser(alice).Do().ExpectSuccess()
		uuid := a.Post("/web/operations/alice/findings").WithJSONBody(`{"title": "Alice's finding", "category": "Enterprise", "description": ""}`).AsUser(alice).Do().ExpectSuccess().ResponseUUID()
		a.Put("/web/operations/alice/findings/" + uuid).WithJSONBody(`{
			"title": "bob was here", 
			"category": "Enterprise", 
			"description": "",
			"readyToReport": false,
			"ticketLink": null
		}`).AsUser(bob).Do().ExpectUnauthorized()

		// Ensure using an operation that bob controlls does not bypass security check
		a.Post("/web/operations").WithJSONBody(`{"name": "Bob's Operation", "slug": "bob"}`).AsUser(bob).Do().ExpectSuccess()
		a.Put("/web/operations/bob/findings/" + uuid).WithJSONBody(`{"title": "bob was here", "category": "Enterprise", "description": "", "readyToReport": false}`).AsUser(bob).Do().ExpectUnauthorized()

		// Ensure finding is unmodified
		a.Get("/web/operations/alice/findings").AsUser(alice).Do().ExpectSubsetJSONArray([]string{`{"title": "Alice's finding"}`})
	})

	t.Run("Ensure users cannot delete findings of an operation they do not have write access to", func(t *testing.T) {
		a := integration.NewTester(t)
		alice := a.NewUser("aowner", "Alice", "Owner")
		bob := a.NewUser("battacker", "Bob", "Attacker")

		a.Post("/web/operations").WithJSONBody(`{"name": "Alice's Operation", "slug": "alice"}`).AsUser(alice).Do().ExpectSuccess()
		uuid := a.Post("/web/operations/alice/findings").WithJSONBody(`{"title": "Alice's finding", "category": "Product", "description": ""}`).AsUser(alice).Do().ExpectSuccess().ResponseUUID()
		a.Delete("/web/operations/alice/findings/" + uuid).AsUser(bob).Do().ExpectUnauthorized()

		// Ensure using an operation that bob controlls does not bypass security check
		a.Post("/web/operations").WithJSONBody(`{"name": "Bob's Operation", "slug": "bob"}`).AsUser(bob).Do().ExpectSuccess()
		a.Delete("/web/operations/bob/findings/" + uuid).AsUser(bob).Do().ExpectUnauthorized()

		// Ensure finding is unmodified
		a.Get("/web/operations/alice/findings").AsUser(alice).Do().ExpectSubsetJSONArray([]string{`{"title": "Alice's finding"}`})
	})
}

func TestAssociatingEvidenceWithFindings(t *testing.T) {
	a := integration.NewTester(t)
	a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

	// Initialize test with an operation with an finding and three evidence
	a.Post("/web/operations").WithJSONBody(`{"name": "op", "slug": "op1"}`).Do().ExpectSuccess()
	evidenceUUID1 := a.Post("/web/operations/op1/evidence").WithMultipartBody(map[string]string{"description": "evi1"}, nil).Do().ExpectSuccess().ResponseUUID()
	evidenceUUID2 := a.Post("/web/operations/op1/evidence").WithMultipartBody(map[string]string{"description": "evi2"}, nil).Do().ExpectSuccess().ResponseUUID()
	evidenceUUID3 := a.Post("/web/operations/op1/evidence").WithMultipartBody(map[string]string{"description": "evi3"}, nil).Do().ExpectSuccess().ResponseUUID()
	findingUUID := a.Post("/web/operations/op1/findings").WithJSONBody(`{"title": "finding", "category": "Product", "description": ""}`).Do().ExpectSuccess().ResponseUUID()

	// Check adding and removing evidence
	a.Put("/web/operations/op1/findings/" + findingUUID + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID1 + `", "` + evidenceUUID2 + `"], "evidenceToRemove": []}`).Do().ExpectSuccess()
	a.Get("/web/operations/op1/findings/" + findingUUID + "/evidence").Do().ExpectSubsetJSONArray([]string{`{"uuid": "` + evidenceUUID1 + `"}`, `{"uuid": "` + evidenceUUID2 + `"}`})
	a.Put("/web/operations/op1/findings/" + findingUUID + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID3 + `"], "evidenceToRemove": ["` + evidenceUUID1 + `"]}`).Do().ExpectSuccess()
	a.Get("/web/operations/op1/findings/" + findingUUID + "/evidence").Do().ExpectSubsetJSONArray([]string{`{"uuid": "` + evidenceUUID2 + `"}`, `{"uuid": "` + evidenceUUID3 + `"}`})

	// Ensure evidence cannot be added to findings in a different operation
	a.Post("/web/operations").WithJSONBody(`{"name": "op", "slug": "op2"}`).Do().ExpectSubsetJSON(
		`{"slug": "op2", "name": "op"}`,
	)
	findingUUID2 := a.Post("/web/operations/op2/findings").WithJSONBody(`{"title": "other finding", "category": "Product", "description": ""}`).Do().ResponseUUID()
	a.Put("/web/operations/op2/findings/" + findingUUID2 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID1 + `"], "evidenceToRemove": []}`).Do().ExpectUnauthorized()
}

func TestTaggingFindings(t *testing.T) {
	a := integration.NewTester(t)
	a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

	// Initialize some tags, evidence, and findings
	a.Post("/web/operations").WithJSONBody(`{"name": "op", "slug": "op"}`).Do().ExpectSuccess()
	a.Post("/web/operations/op/tags").WithJSONBody(`{"name": "Exploitation", "colorName": "red"}`).Do().ExpectSuccess()
	a.Post("/web/operations/op/tags").WithJSONBody(`{"name": "Lateral Movement", "colorName": "blue"}`).Do().ExpectSuccess()
	evidenceUUID1 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "e1", "tagIds": "[1]", "occurredAt": "2019-05-01T10:00:00Z"}, nil).Do().ExpectSuccess().ResponseUUID()
	evidenceUUID2 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "e2", "tagIds": "[2]", "occurredAt": "2019-06-01T10:00:00Z"}, nil).Do().ExpectSuccess().ResponseUUID()
	evidenceUUID3 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "e3", "occurredAt": "2019-07-01T10:00:00Z"}, nil).Do().ExpectSuccess().ResponseUUID()
	findingUUID1 := a.Post("/web/operations/op/findings").WithJSONBody(`{"title": "f1", "category": "Network", "description": ""}`).Do().ExpectSuccess().ResponseUUID()
	findingUUID2 := a.Post("/web/operations/op/findings").WithJSONBody(`{"title": "f2", "category": "Network", "description": ""}`).Do().ExpectSuccess().ResponseUUID()
	findingUUID3 := a.Post("/web/operations/op/findings").WithJSONBody(`{"title": "f3", "category": "Network", "description": ""}`).Do().ExpectSuccess().ResponseUUID()

	// Add evidence to findings
	a.Put("/web/operations/op/findings/" + findingUUID1 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID1 + `", "` + evidenceUUID2 + `"], "evidenceToRemove": []}`).Do().ExpectSuccess()
	a.Put("/web/operations/op/findings/" + findingUUID2 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID2 + `", "` + evidenceUUID3 + `"], "evidenceToRemove": []}`).Do().ExpectSuccess()
	a.Put("/web/operations/op/findings/" + findingUUID3 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID3 + `", "` + evidenceUUID1 + `"], "evidenceToRemove": []}`).Do().ExpectSuccess()

	// Check that tags are populated from attached evidence
	a.Get("/web/operations/op/findings").Do().ExpectSubsetJSONArray([]string{
		`{"title": "f2", "tags": [{"id": 2, "name": "Lateral Movement", "colorName": "blue"}], "occurredFrom": "2019-06-01T10:00:00Z", "occurredTo": "2019-07-01T10:00:00Z"}`,
		`{"title": "f3", "tags": [{"id": 1, "name": "Exploitation", "colorName": "red"}], "occurredFrom": "2019-05-01T10:00:00Z", "occurredTo": "2019-07-01T10:00:00Z"}`,
		`{"title": "f1", "tags": [{"id": 1,"name": "Exploitation", "colorName": "red"}, {"id": 2, "name": "Lateral Movement", "colorName": "blue"}], "occurredFrom": "2019-05-01T10:00:00Z", "occurredTo": "2019-06-01T10:00:00Z"}`,
	})
}
