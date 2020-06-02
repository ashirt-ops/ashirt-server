// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package integration_test

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/integration"
)

func expectQueryToReturnIDs(t *testing.T, a *integration.Tester, query string, expectedUUIDs []string) {
	t.Helper()
	params := url.Values{}
	params.Add("query", query)

	body := a.Get("/web/operations/op/evidence?" + params.Encode()).Do().ExpectSuccess().ResponseBody()
	var evidence []map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &evidence))

	actualUUIDs := []string{}
	for _, evi := range evidence {
		actualUUIDs = append(actualUUIDs, evi["uuid"].(string))
	}

	require.Equal(t, expectedUUIDs, actualUUIDs)
}

func TestSavedQueries(t *testing.T) {
	a := integration.NewTester(t)
	a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

	// Setup
	a.Post("/web/operations").WithJSONBody(`{"name": "Operation", "slug": "op"}`).Do().ExpectSubsetJSON(`{"name": "Operation", "status": 0}`)

	// Test "starts empty"
	a.Get("/web/operations/op/queries").Do().ExpectSuccess().ExpectJSON(`[]`)

	// Test "add one entry"
	const firstEntry string = `{"name": "First", "query": "one", "type": "evidence"}`
	a.Post("/web/operations/op/queries").
		WithJSONBody(firstEntry).
		Do().
		ExpectJSON(`{"id": 1, "name": "First", "query": "one", "type": "evidence"}`)

	// Test "contains one entry"
	a.Get("/web/operations/op/queries").Do().ExpectSuccess().
		ExpectSubsetJSONArray([]string{firstEntry})

		// Test "supports tags" (setup)
	const complexEntry string = `{"name": "Second", "query": "one tag:Plain tag:\"Is Complex\"", "type": "evidence"}`
	a.Post("/web/operations/op/queries").
		WithJSONBody(complexEntry).
		Do().
		ExpectJSON(`{"id": 2, "name": "Second", "query": "one tag:Plain tag:\"Is Complex\"", "type": "evidence"}`)

	// Test "supports tags"
	a.Get("/web/operations/op/queries").Do().ExpectSuccess().
		ExpectSubsetJSONArray([]string{firstEntry, complexEntry})

	// Test "won't allow duplicates"
	a.Post("/web/operations/op/queries").
		WithJSONBody(complexEntry).
		Do().
		ExpectStatus(http.StatusBadRequest)
}

func TestQueryingEvidence(t *testing.T) {
	a := integration.NewTester(t)
	a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

	// Setup
	a.Post("/web/operations").WithJSONBody(`{"name": "Operation", "slug": "op"}`).Do().ExpectSuccess()

	a.Post("/web/operations/op/tags").WithJSONBody(`{"name": "Lateral Movement"}`).Do().ExpectSuccess()
	a.Post("/web/operations/op/tags").WithJSONBody(`{"name": "Exploitation"}`).Do().ExpectSuccess()
	a.Post("/web/operations/op/tags").WithJSONBody(`{"name": "Password Cracking"}`).Do().ExpectSuccess()

	e1 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "Moved laterally", "tagIds": "[1]", "occurredAt": "2019-05-01T15:00:00Z"}, nil).Do().ExpectSuccess().ResponseUUID()
	e2 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "Cracked weak password found on dev", "tagIds": "[3]", "occurredAt": "2019-05-29T15:00:00Z"}, nil).Do().ExpectSuccess().ResponseUUID()
	e3 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "Used found password on prod", "tagIds": "[1, 3]", "occurredAt": "2019-06-21T15:00:00Z"}, nil).Do().ExpectSuccess().ResponseUUID()
	e4 := a.Post("/web/operations/op/evidence").WithMultipartBody(map[string]string{"description": "Last Evidence", "tagIds": "[2]", "occurredAt": "2019-08-09T15:00:00Z"}, nil).Do().ExpectSuccess().ResponseUUID()

	// Test querys
	expectQueryToReturnIDs(t, a, ``, []string{e4, e3, e2, e1})
	expectQueryToReturnIDs(t, a, `found password`, []string{e3, e2})
	expectQueryToReturnIDs(t, a, `"found password"`, []string{e3})
	expectQueryToReturnIDs(t, a, `tag:"Password Cracking"`, []string{e3, e2})
	expectQueryToReturnIDs(t, a, `tag:"Password Cracking" tag:"Lateral Movement"`, []string{e3})
	expectQueryToReturnIDs(t, a, `dev tag:"Lateral Movement"`, []string{})
	expectQueryToReturnIDs(t, a, `tag:UnknownTag`, []string{})
	expectQueryToReturnIDs(t, a, `range:2018-01-01,2018-12-31`, []string{})
	expectQueryToReturnIDs(t, a, `range:2019-01-01,2019-12-31`, []string{e4, e3, e2, e1})
	expectQueryToReturnIDs(t, a, `range:2019-05-25,2019-06-25`, []string{e3, e2})
}
