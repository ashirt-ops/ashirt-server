// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package integration_test

import (
	"testing"

	"github.com/theparanoids/ashirt/backend/integration"
)

func TestOperations(t *testing.T) {
	a := integration.NewTester(t)
	a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

	// Ensure there are no operations
	a.Get("/web/operations").Do().ExpectJSON(`[]`)

	// Create an operation from web and ensure it's queryable
	a.Post("/web/operations").
		WithJSONBody(`{"name": "My operation", "slug": "my-op"}`).Do().
		ExpectSubsetJSON(`{"slug": "my-op", "name": "My operation", "numUsers": 1, "status": 0}`)
	a.Get("/web/operations").Do().ExpectSubsetJSONArray([]string{`{"name": "My operation", "status": 0, "numUsers": 1}`})

	// Updating operations
	bob := a.NewUser("bsnooper", "Bob", "Snooper")
	a.Post("/web/operations").
		WithJSONBody(`{"name": "Original Name", "slug": "other-op"}`).Do().
		ExpectSubsetJSON(`{"slug": "other-op", "name": "Original Name", "numUsers": 1, "status": 0}`)
	a.Put("/web/operations/other-op").WithJSONBody(`{"name": "New Name", "status": 2}`).Do().ExpectSuccess()
	a.Get("/web/operations/other-op").Do().ExpectSubsetJSON(`{"name": "New Name", "status": 2}`)

	// Ensure other users cannot update your operations
	a.Put("/web/operations/other-op").WithJSONBody(`{"name": "Bob's operation"}`).AsUser(bob).Do().ExpectUnauthorized()
	a.Get("/web/operations/other-op").Do().ExpectSubsetJSON(`{"name": "New Name", "status": 2}`)
}

func TestRequestingNonexistentOperations(t *testing.T) {
	a := integration.NewTester(t)
	a.DefaultUser = a.NewUser("adefaultuser", "Alice", "DefaultUser")

	// Querying operations that don't exist returns not found
	a.Get("/web/operations/op/findings").Do().ExpectNotFound()
	a.Get("/web/operations/op/findings/1").Do().ExpectNotFound()
	a.Get("/web/operations/op/evidence").Do().ExpectNotFound()
	a.Get("/web/operations/op/evidence/1").Do().ExpectNotFound()

	// Writing operations that don't exist returns unauthorized
	a.Post("/web/operations/op/findings").WithJSONBody(`{"title": "e1", "category": "CD", "description": ""}`).Do().ExpectUnauthorized()
	a.Put("/web/operations/op").WithJSONBody(`{"name": "new name"}`).Do().ExpectUnauthorized()
}
