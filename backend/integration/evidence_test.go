// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvidence(t *testing.T) {
	state := preamble(t)

	// Ensure there are no evidence on operation 2
	state.Alice.Get("/web/operations/" + state.Operations[1].Slug + "/evidence").Do().ExpectJSON(`[]`)

	// Add an evidence to operation 2 and ensure it's queryable
	testWebCreateEvidenceAndEnsureWebQueryable(state.Alice, state.Operations[1])

	// Add codeblock
	testWebUploadCodeblock(state.Alice, state.Operations[0], state)

	// Create Evidence with tags (via web interface)
	evidenceUUID4 := testCreateEvidenceWithTags(state.Alice, state.Operations[2], state)

	// Modify Evidence with tags
	testModifyEvidence(state, evidenceUUID4)

	t.Run("Tagging evidence", func(t *testing.T) {
		a := integration.NewTester(t)
		a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

		a.Post("/web/operations").WithJSONBody(`{"name": "Op 1", "slug": "op1"}`).Do().ExpectSuccess()
		a.Post("/web/operations").WithJSONBody(`{"name": "Op 2", "slug": "op2"}`).Do().ExpectSuccess()
		a.Post("/web/operations/op1/tags").WithJSONBody(`{"name": "Alpha", "colorName": "red"}`).Do().ExpectSubsetJSON(`{"id": 1}`)
		a.Post("/web/operations/op1/tags").WithJSONBody(`{"name": "Beta", "colorName": "blue"}`).Do().ExpectSubsetJSON(`{"id": 2}`)
		a.Post("/web/operations/op2/tags").WithJSONBody(`{"name": "Gamma", "colorName": "green"}`).Do().ExpectSubsetJSON(`{"id": 3}`)
		a.Post("/web/operations/op2/tags").WithJSONBody(`{"name": "Delta", "colorName": "yellow"}`).Do().ExpectSubsetJSON(`{"id": 4}`)

		// Ensure we can create/edit evidence with tags that belong to its operation
		evidenceUUID1 := a.Post("/web/operations/op1/evidence").WithMultipartBody(map[string]string{"description": "a", "tagIds": "[1]"}, nil).Do().ExpectStatus(http.StatusCreated).ResponseUUID()
		evidenceUUID2 := a.Post("/web/operations/op2/evidence").WithMultipartBody(map[string]string{"description": "a", "tagIds": "[3]"}, nil).Do().ExpectStatus(http.StatusCreated).ResponseUUID()
		a.Put("/web/operations/op1/evidence/"+evidenceUUID1).WithMultipartBody(map[string]string{"description": "a", "tagsToAdd": "[2]", "tagsToRemove": "[1]"}, nil).Do().ExpectSuccess()
		a.Put("/web/operations/op2/evidence/"+evidenceUUID2).WithMultipartBody(map[string]string{"description": "a", "tagsToAdd": "[4]", "tagsToRemove": "[3]"}, nil).Do().ExpectSuccess()

		// Ensure we can't create/edit evidence with tags from other operations
		a.Post("/web/operations/op1/evidence").WithMultipartBody(map[string]string{"description": "foo", "tagIds": "[3]"}, nil).Do().ExpectStatus(http.StatusBadRequest)
		a.Put("/web/operations/op1/evidence/"+evidenceUUID1).WithMultipartBody(map[string]string{"description": "a", "tagsToAdd": "[3]", "tagsToRemove": "[2]"}, nil).Do().ExpectStatus(http.StatusBadRequest)
		a.Put("/web/operations/op2/evidence/"+evidenceUUID2).WithMultipartBody(map[string]string{"description": "a", "tagsToAdd": "[2]", "tagsToRemove": "[4]"}, nil).Do().ExpectStatus(http.StatusBadRequest)
	})

	t.Run("Editing evidence", func(t *testing.T) {
		a := integration.NewTester(t)
		a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

		a.Post("/web/operations").WithJSONBody(`{"name": "Op 1", "slug": "op"}`).Do().ExpectSuccess()
		uuid := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "Evidence 1"}, nil).Do().ExpectSuccess().ResponseUUID()
		a.Put("/web/operations/op/evidence/"+uuid).WithMultipartBody(map[string]string{"description": "Updated description"}, nil).Do().ExpectSuccess()
		a.Get("/web/operations/op/evidence/" + uuid).Do().ExpectSubsetJSON(`{"description": "Updated description"}`)
	})

	t.Run("Editing evidence with an image", func(t *testing.T) {
		a := integration.NewTester(t)
		a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

		// Create evidence with an image
		a.Post("/web/operations").WithJSONBody(`{"name": "Op 1", "slug": "op"}`).Do().ExpectSuccess()
		f, _ := os.Open("fixtures/screenshot.png")
		defer f.Close()
		uuid := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "evi"}, map[string]*os.File{"content": f}).Do().ExpectStatus(http.StatusCreated).ResponseUUID()

		// evidence images cannot be edited
		f, _ = os.Open("fixtures/dummyCode.json")
		defer f.Close()
		a.Put("/web/operations/op/evidence/"+uuid).WithMultipartBody(map[string]string{"description": "new image"}, map[string]*os.File{"content": f}).Do().ExpectStatus(http.StatusBadRequest)

		// Editing image description
		a.Put("/web/operations/op/evidence/"+uuid).WithMultipartBody(map[string]string{"description": "new description"}, nil).Do().ExpectSuccess()

		// Ensure image is unmodified
		imageData, _ := ioutil.ReadFile("fixtures/screenshot.png")
		a.Get("/web/operations/op/evidence/"+uuid+"/media").Do().ExpectResponse(200, imageData)
	})

	t.Run("Deleting evidence", func(t *testing.T) {
		a := integration.NewTester(t)
		a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

		a.Post("/web/operations").WithJSONBody(`{"name": "Op 1", "slug": "op"}`).Do().ExpectSuccess()
		uuid := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "To be deleted..."}, nil).Do().ExpectSuccess().ResponseUUID()
		a.Delete("/web/operations/op/evidence/" + uuid).WithJSONBody(`{"deleteAssociatedFindings": false}`).Do().ExpectSuccess()
		a.Get("/web/operations/op/evidence").Do().ExpectJSON(`[]`)
	})

	t.Run("Deleting evidence with attached findings", func(t *testing.T) {
		a := integration.NewTester(t)
		a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

		// Create 2 evidence and 2 findings and associate them
		a.Post("/web/operations").WithJSONBody(`{"name": "Op 1", "slug": "op"}`).Do().ExpectSuccess()
		evidenceUUID1 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "Evidence 1"}, nil).Do().ExpectSuccess().ResponseUUID()
		evidenceUUID2 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "Evidence 2"}, nil).Do().ExpectSuccess().ResponseUUID()

		findingUUID1 := a.Post("/web/operations/op/findings").WithJSONBody(`{"title": "Finding 1", "category": "CD", "description": ""}`).Do().ExpectSuccess().ResponseUUID()
		findingUUID2 := a.Post("/web/operations/op/findings").WithJSONBody(`{"title": "Finding 2", "category": "CD", "description": ""}`).Do().ExpectSuccess().ResponseUUID()

		a.Put("/web/operations/op/findings/" + findingUUID1 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID1 + `"], "evidenceToRemove": []}`).Do().ExpectSuccess()
		a.Put("/web/operations/op/findings/" + findingUUID2 + "/evidence").WithJSONBody(`{"evidenceToAdd": ["` + evidenceUUID2 + `"], "evidenceToRemove": []}`).Do().ExpectSuccess()

		a.Delete("/web/operations/op/evidence/" + evidenceUUID1).WithJSONBody(`{"deleteAssociatedFindings": false}`).Do().ExpectSuccess()
		a.Delete("/web/operations/op/evidence/" + evidenceUUID2).WithJSONBody(`{"deleteAssociatedFindings": true}`).Do().ExpectSuccess()

		// Second finding should be deleted from deleteAssociatedFindings: true
		a.Get("/web/operations/op/findings").Do().ExpectSubsetJSONArray([]string{`{"uuid": "` + findingUUID1 + `", "title": "Finding 1"}`})
	})

	t.Run("Ensure users cannot edit evidence of an operation they do not have write access to", func(t *testing.T) {
		a := integration.NewTester(t)
		alice := a.NewUser("aowner", "Alice", "Owner")
		bob := a.NewUser("battacker", "Bob", "Attacker")

		a.Post("/web/operations").WithJSONBody(`{"name": "Alice's Operation", "slug": "alice"}`).AsUser(alice).Do().ExpectSuccess()
		uuid := a.Post("/web/operations/alice/evidence").WithMultipartBody(map[string]string{"description": "alice's evidence"}, nil).AsUser(alice).Do().ExpectSuccess().ResponseUUID()

		a.Put("/web/operations/alice/evidence/"+uuid).WithMultipartBody(map[string]string{"description": "bob was here"}, nil).AsUser(bob).Do().ExpectUnauthorized()

		// Ensure using an operation that bob controlls does not bypass security check
		a.Post("/web/operations").WithJSONBody(`{"name": "Bob's Operation", "slug": "bob"}`).AsUser(bob).Do().ExpectSuccess()
		a.Put("/web/operations/bob/evidence/"+uuid).WithMultipartBody(map[string]string{"description": "bob was here"}, nil).AsUser(bob).Do().ExpectUnauthorized()

		// Ensure evidence is unmodified
		a.Get("/web/operations/alice/evidence/" + uuid).AsUser(alice).Do().ExpectSubsetJSON(`{"description": "alice's evidence"}`)
	})

	t.Run("Ensure users cannot delete evidence of an operation they do not have write access to", func(t *testing.T) {
		a := integration.NewTester(t)
		alice := a.NewUser("aowner", "Alice", "Owner")
		bob := a.NewUser("battacker", "Bob", "Attacker")

		a.Post("/web/operations").WithJSONBody(`{"name": "Alice's Operation", "slug": "alice"}`).AsUser(alice).Do().ExpectSuccess()
		uuid := a.Post("/web/operations/alice/evidence").WithMultipartBody(map[string]string{"description": "alice's evidence"}, nil).AsUser(alice).Do().ExpectSuccess().ResponseUUID()

		a.Delete("/web/operations/alice/evidence/" + uuid).AsUser(bob).Do().ExpectUnauthorized()

		// Ensure using an operation that bob controlls does not bypass security check
		a.Post("/web/operations").WithJSONBody(`{"name": "Bob's Operation", "slug": "bob"}`).AsUser(bob).Do().ExpectSuccess()
		a.Delete("/web/operations/bob/evidence/" + uuid).AsUser(bob).Do().ExpectUnauthorized()

		// Ensure evidence still exists
		a.Get("/web/operations/alice/evidence/" + uuid).AsUser(alice).Do().ExpectSubsetJSON(`{"description": "alice's evidence"}`)
	})
}

type testState struct {
	Alice       *integration.Tester
	CreatedTags []dtos.Tag
	Operations  []dtos.Operation
}

func (ts testState) tagIDsAsJson(indexes ...int) string {
	tagIDs := make([]string, 0, len(indexes))
	for i := range indexes {
		idIndex := indexes[i]
		if idIndex < 0 || idIndex >= len(ts.CreatedTags) {
			continue
		}
		tagIDs = append(tagIDs, fmt.Sprintf("%v", ts.CreatedTags[idIndex].ID))
	}
	return fmt.Sprintf("[%v]", strings.Join(tagIDs, ", "))
}

func tagIndexOf(needle int64, haystack []dtos.Tag) int64 {
	for i, v := range haystack {
		if v.ID == needle {
			return int64(i)
		}
	}
	return -1
}

func preamble(t *testing.T) testState {
	a := integration.NewTester(t)
	a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

	tagsJSON := []string{
		`{"name": "one", "colorName": "red"}`,
		`{"name": "two", "colorName": "green"}`,
		`{"name": "three", "colorName": "blue"}`,
		`{"name": "foo", "colorName": "red"}`,
		`{"name": "bar", "colorName": "green"}`,
		`{"name": "baz", "colorName": "blue"}`,
	}

	fullTagJson := fmt.Sprintf("[%v]", strings.Join(tagsJSON, `, `))

	var tags []dtos.Tag
	assert.NoError(t, json.Unmarshal([]byte(fullTagJson), &tags))

	// Create three operations
	opsJSON := []string{
		`{"name": "op1", "slug": "op1"}`,
		`{"name": "op2", "slug": "op2"}`,
		`{"name": "tagOpTest", "slug": "tag-op-test"}`,
	}
	a.Post("/web/operations").WithJSONBody(opsJSON[0]).Do().ExpectSuccess()
	a.Post("/web/operations").WithJSONBody(opsJSON[1]).Do().ExpectSuccess()
	a.Post("/web/operations").WithJSONBody(opsJSON[2]).Do().ExpectSuccess()

	// Create some tags for later on
	tags[0].ID = a.Post("/web/operations/op1/tags").WithJSONBody(tagsJSON[0]).Do().ExpectSuccess().ResponseID()
	tags[1].ID = a.Post("/web/operations/op1/tags").WithJSONBody(tagsJSON[1]).Do().ExpectSuccess().ResponseID()
	tags[2].ID = a.Post("/web/operations/op1/tags").WithJSONBody(tagsJSON[2]).Do().ExpectSuccess().ResponseID()
	tags[3].ID = a.Post("/web/operations/tag-op-test/tags").WithJSONBody(tagsJSON[3]).Do().ExpectSuccess().ResponseID()
	tags[4].ID = a.Post("/web/operations/tag-op-test/tags").WithJSONBody(tagsJSON[4]).Do().ExpectSuccess().ResponseID()
	tags[5].ID = a.Post("/web/operations/tag-op-test/tags").WithJSONBody(tagsJSON[5]).Do().ExpectSuccess().ResponseID()

	var ops []dtos.Operation
	opsJsonArray := fmt.Sprintf("[%v]", strings.Join(opsJSON, `, `))
	assert.NoError(t, json.Unmarshal([]byte(opsJsonArray), &ops))
	getOpsIDs(a, &ops)

	return testState{
		Alice:       a,
		CreatedTags: tags,
		Operations:  ops,
	}
}

func testWebCreateEvidenceAndEnsureWebQueryable(a *integration.Tester, op dtos.Operation) string {
	resultUUID := a.Post("/web/operations/"+op.Slug+"/evidence").WithMultipartBody(map[string]string{"description": "evi2"}, nil).Do().ResponseUUID()
	a.Get("/web/operations/" + op.Slug + "/evidence").Do().ExpectSubsetJSONArray([]string{`{"uuid": "` + resultUUID + `", "description": "evi2"}`})
	return resultUUID
}

func testWebUploadCodeblock(a *integration.Tester, op dtos.Operation, state testState) string {
	file, err := os.Open("fixtures/dummyCode.json")
	require.NoError(a.TestingT(), err)
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	require.NoError(a.TestingT(), err)
	file.Seek(0, 0)

	desc := "Codeblock With Source"
	resultUUID := a.
		Post("/web/operations/"+op.Slug+"/evidence").
		WithMultipartBody(
			map[string]string{"description": desc, "operator_id": "1",
				"tagIds": state.tagIDsAsJson(0, 1), "contentType": "codeblock"},
			map[string]*os.File{"content": file},
		).
		Do().
		ExpectStatus(http.StatusCreated).
		ResponseUUID()

	a.Get("/web/operations/" + op.Slug + "/evidence/" + resultUUID).Do().ExpectSubsetJSON(fmt.Sprintf(`{"description": "%v"}`, desc))
	encodedCodeBlock := a.Get("/web/operations/" + op.Slug + "/evidence/" + resultUUID + "/media").Do().ResponseBody()

	require.Equal(a.TestingT(), fileContent, encodedCodeBlock)

	return resultUUID
}

func testCreateEvidenceWithTags(a *integration.Tester, op dtos.Operation, state testState) string {
	tags := state.CreatedTags

	desc := "blah blah blah"

	evidenceUUID4 := a.Post("/web/operations/"+op.Slug+"/evidence").WithMultipartBody(
		map[string]string{
			"description": desc,
			"tagIds":      state.tagIDsAsJson(3, 4),
		},
		nil,
	).Do().ExpectSuccess().ResponseUUID()

	reconstituted, _ := json.Marshal([]dtos.Tag{tags[3], tags[4]})
	expectedJSON := `{"uuid": "` + evidenceUUID4 + `", "description": "` + desc + `", "tags": ` + string(reconstituted) + `}`
	a.Get("/web/operations/" + op.Slug + "/evidence").Do().ExpectSubsetJSONArray([]string{expectedJSON})
	return evidenceUUID4
}

func testModifyEvidence(state testState, evidenceUUID string) {
	a := state.Alice
	tags := state.CreatedTags

	newDesc := "new"
	a.Put("/web/operations/tag-op-test/evidence/"+evidenceUUID).WithMultipartBody(
		map[string]string{"description": newDesc, "tagsToAdd": "[6]", "tagsToRemove": "[4]"},
		nil,
	).Do().ExpectSuccess()
	reconstituted, _ := json.Marshal([]dtos.Tag{tags[4], tags[5]})
	expectedJSON := `{"uuid": "` + evidenceUUID + `", "description": "` + newDesc + `", "tags": ` + string(reconstituted) + `}`
	a.Get("/web/operations/tag-op-test/evidence").
		Do().ExpectSubsetJSONArray([]string{expectedJSON})
}

func getOpsIDs(a *integration.Tester, ops *[]dtos.Operation) {
	rb := a.Get("/web/operations").Do().ResponseBody()
	var temp []dtos.Operation
	assert.NoError(a.TestingT(), json.Unmarshal([]byte(rb), &temp))
	for idx, op := range *ops {
		for _, apiOp := range temp {
			if apiOp.Slug == op.Slug {
				(*ops)[idx].ID = apiOp.ID
				break
			}
		}
	}
}
