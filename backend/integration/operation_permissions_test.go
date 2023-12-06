package integration_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/integration"
)

func TestUserPermissions(t *testing.T) {
	a := integration.NewTester(t)
	a.NewUser("True.Admin", "Terrence", "Administrator") // admins necessarily break permission steps, so ignoring them here
	creator := a.NewUser("creator", "Charlie", "Creator")
	reader := a.NewUser("reader", "Rupert", "Reader")
	writer := a.NewUser("writer", "Wendy", "Writer")
	admin := a.NewUser("admin", "Alice", "Admin")
	otherUser := a.NewUser("otheruser", "Oscar", "Otheruser")

	a.Post("/web/operations").WithJSONBody(`{"name": "op", "slug": "op"}`).AsUser(creator).Do().ExpectSuccess()
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "rupert.reader", "role": "read"}`).AsUser(creator).Do().ExpectSuccess()
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "wendy.writer", "role": "write"}`).AsUser(creator).Do().ExpectSuccess()
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "alice.admin", "role": "admin"}`).AsUser(creator).Do().ExpectSuccess()

	a.Get("/web/operations/op/users").AsUser(creator).Do().ExpectJSON(`
	  [
		{"role": "admin", "user": {"slug": "charlie.creator", "firstName": "Charlie", "lastName": "Creator"}},
		{"role": "read",  "user": {"slug": "rupert.reader",   "firstName": "Rupert",  "lastName": "Reader"}},
		{"role": "write", "user": {"slug": "wendy.writer",    "firstName": "Wendy",   "lastName": "Writer"}},
		{"role": "admin", "user": {"slug": "alice.admin",     "firstName": "Alice",   "lastName": "Admin"}}
	  ]`)

	// Setting nonexistent user permissions results in a 401
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "nonexistentuser", "role": "admin"}`).AsUser(creator).Do().ExpectStatus(400)

	// Users not belonging to an operation cannot read or set permissions
	a.Get("/web/operations/op/users").AsUser(otherUser).Do().ExpectNotFound()
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "rupert.reader", "role": "write"}`).AsUser(otherUser).Do().ExpectUnauthorized()

	// Readers and writers cannot set permissions
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "oscar.otheruser", "role": "read"}`).AsUser(reader).Do().ExpectUnauthorized()
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "oscar.otheruser", "role": "read"}`).AsUser(writer).Do().ExpectUnauthorized()
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "wendy.writer",    "role": "read"}`).AsUser(reader).Do().ExpectUnauthorized()
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "rupert.reader",   "role": "write"}`).AsUser(writer).Do().ExpectUnauthorized()

	// Admins cannot demote themselves (anti-lockout)
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "charlie.creator", "role": "write"}`).AsUser(creator).Do().ExpectUnauthorized()

	// Admins can change other's user permissions
	a.Patch("/web/operations/op/users").WithJSONBody(`{"userSlug": "charlie.creator", "role": "write"}`).AsUser(admin).Do().ExpectSuccess()
	a.Get("/web/operations/op/users").AsUser(admin).Do().ExpectJSON(`
	  [
		{"role": "write", "user": {"slug": "charlie.creator", "firstName": "Charlie", "lastName": "Creator"}},
		{"role": "read",  "user": {"slug": "rupert.reader",   "firstName": "Rupert",  "lastName": "Reader"}},
		{"role": "write", "user": {"slug": "wendy.writer",    "firstName": "Wendy",   "lastName": "Writer"}},
		{"role": "admin", "user": {"slug": "alice.admin",     "firstName": "Alice",   "lastName": "Admin"}}
	  ]`)
}
