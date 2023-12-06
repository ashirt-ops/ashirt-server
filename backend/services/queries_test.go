package services_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"
)

func TestCreateQuery(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		op := OpChamberOfSecrets
		i := services.CreateQueryInput{
			OperationSlug: op.Slug,
			Name:          "Evidence By author",
			Query:         "<query goes here>",
			Type:          "findings",
		}
		createdQuery, err := services.CreateQuery(ctx, db, i)
		require.NoError(t, err)
		fullQuery := getQueryByID(t, db, createdQuery.ID)
		require.Equal(t, i.Name, fullQuery.Name)
		require.Equal(t, i.Type, fullQuery.Type)
		require.Equal(t, i.Query, fullQuery.Query)
	})
}

func TestDeleteQuery(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		i := services.DeleteQueryInput{
			OperationSlug: OpChamberOfSecrets.Slug,
			ID:            QuerySalazarsHier.ID,
		}

		getQueryCount := makeDBRowCounter(t, db, "queries", "id=?", i.ID)
		require.Equal(t, int64(1), getQueryCount(), "Database should have item to delete")

		err := services.DeleteQuery(ctx, db, i)
		require.NoError(t, err)
		require.Equal(t, int64(0), getQueryCount(), "Database should have deleted the item")
	})
}

func TestListQueriesForOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		masterOp := OpChamberOfSecrets
		allQueries := getQueriesForOperationID(t, db, masterOp.ID)
		require.NotEqual(t, len(allQueries), 0, "Some number of queries should exist")

		foundQueries, err := services.ListQueriesForOperation(ctx, db, masterOp.Slug)
		require.NoError(t, err)
		require.Equal(t, len(foundQueries), len(allQueries))
		validateQuerySets(t, foundQueries, allQueries, validateQuery)
	})
}

func TestUpdateQuery(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		masterOp := OpChamberOfSecrets
		masterQuery := QuerySalazarsHier
		input := services.UpdateQueryInput{
			OperationSlug: masterOp.Slug,
			ID:            masterQuery.ID,
			Name:          "New Name",
			Query:         "New Query",
		}

		err := services.UpdateQuery(ctx, db, input)
		require.NoError(t, err)

		updatedQuery := getQueryByID(t, db, masterQuery.ID)

		require.NoError(t, err)
		require.Equal(t, input.Name, updatedQuery.Name)
		require.Equal(t, input.Query, updatedQuery.Query)
	})
}

func TestUpsertQuery(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		user := UserRon
		ctx := contextForUser(user, db)
		op := OpChamberOfSecrets

		checkQueryData := func(i services.UpsertQueryInput, dbModel models.Query) {
			require.Equal(t, i.Name, dbModel.Name)
			require.Equal(t, i.Type, dbModel.Type)
			require.Equal(t, i.Query, dbModel.Query)
		}

		i := services.UpsertQueryInput{
			CreateQueryInput: services.CreateQueryInput{
				OperationSlug: op.Slug,
				Name:          "Ron's amazing query",
				Query:         "I can sepll!",
				Type:          "evidence",
			},
			ReplaceName: false, // initial value doesn't matter
		}
		createdQuery, err := services.UpsertQuery(ctx, db, i)
		require.NoError(t, err)
		fullQuery := getQueryByID(t, db, createdQuery.ID)
		checkQueryData(i, fullQuery)

		// Update the query
		i.CreateQueryInput.Query = "I can spell!"
		updatedQuery, err := services.UpsertQuery(ctx, db, i)
		require.Equal(t, createdQuery.ID, updatedQuery.ID)
		require.NoError(t, err)
		fullUpdatedQuery := getQueryByID(t, db, updatedQuery.ID)
		checkQueryData(i, fullUpdatedQuery)

		// update the name
		i.CreateQueryInput.Name = "Ron's pretty good query"
		i.ReplaceName = true
		updatedQuery, err = services.UpsertQuery(ctx, db, i)
		require.Equal(t, createdQuery.ID, updatedQuery.ID)
		require.NoError(t, err)
		fullUpdatedQuery = getQueryByID(t, db, updatedQuery.ID)
		checkQueryData(i, fullUpdatedQuery)
	})
}

type queryValidator func(*testing.T, models.Query, *dtos.Query)

func validateQuery(t *testing.T, expected models.Query, actual *dtos.Query) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.Query, actual.Query)
	require.Equal(t, expected.Type, actual.Type)
}

func validateQuerySets(t *testing.T, dtoSet []*dtos.Query, dbSet []models.Query, validator queryValidator) {
	var expected *models.Query = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.ID == dtoItem.ID {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validator(t, *expected, dtoItem)
	}
}
