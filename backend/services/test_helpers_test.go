package services_test

import (
	"sync"
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/database/seeding"
)

var retrieveMutex sync.Mutex
var servicesDB *database.Connection

func GetReusableTestDatabase(t *testing.T) *database.Connection {
	retrieveMutex.Lock()
	if servicesDB == nil {
		servicesDB = initTest(t)
	}
	retrieveMutex.Unlock()
	return servicesDB
}

// RunResettableDBTest creates a database connection, seeds the database with the standard seed, then
// runs the test. This can be used for all tests to provide a consistant test environment. Note that
// this uses RunReusableDBTestWithSeed under-the-hood, and so is not appropriate for tests that require
// a fresh database.
func RunResettableDBTest(t *testing.T, fn func(*database.Connection, TestSeedData)) {
	RunReusableDBTestWithSeed(t, HarryPotterSeedData, fn)
}

// RunReusableDBTestWithSeed re-uses the already-established database connection. The seed data is reset,
// but any auto-incremented values within the database will remain as-is. This may be important for
// specific tests. See RunDisposableDBTestWithSeed if you think you might need a completely fresh
// database instance.
func RunReusableDBTestWithSeed(t *testing.T, seed TestSeedData, fn func(*database.Connection, TestSeedData)) {
	db := GetReusableTestDatabase(t)
	seed.ApplyTo(t, db)
	defer seeding.ClearDB(db)
	fn(db, seed)
}

// RunDisposableDBTestWithSeed creates a connection to a database server, creates a fresh database instance,
// and returns that instance with some seed data applied. This differs from RunReusableDBTestWithSeed
// by the destroy-create process. At the end of this process, the network connection is closed, and
// future usages will have a fresh instance.
func RunDisposableDBTestWithSeed(t *testing.T, seed TestSeedData, fn func(*database.Connection, TestSeedData)) {
	db := initTest(t)
	defer db.DB.Close()
	seed.ApplyTo(t, db)
	defer seeding.ClearDB(db)
	fn(db, seed)
}
