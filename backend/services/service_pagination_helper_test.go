package services_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/server/remux"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"
)

func TestParseRequestQueryPagination(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever?page=2&pageSize=3", nil)
	dr := remux.DissectJSONRequest(r)

	pagination := services.ParseRequestQueryPagination(dr, 10)

	require.Nil(t, dr.Error)
	require.Equal(t, int64(2), pagination.Page)
	require.Equal(t, int64(3), pagination.PageSize)
}

func TestPaginationWrapData(t *testing.T) {
	p := services.Pagination{
		PageSize:   100,
		Page:       2,
		TotalCount: 10000,
	}
	expectedContent := []string{"one", "two", "three"}
	resp := p.WrapData(expectedContent)

	require.Equal(t, p.Page, resp.PageNumber)
	require.Equal(t, p.PageSize, resp.PageSize)
	require.Equal(t, expectedContent, resp.Content)
	require.Equal(t, p.TotalCount, resp.TotalCount)
	require.Equal(t, int64(100), resp.TotalPages)
}

func TestPaginationSelect(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		// full page
		checkTagSubset(t, services.Pagination{PageSize: 2, Page: 1}, db)

		// Second page
		checkTagSubset(t, services.Pagination{PageSize: 2, Page: 2}, db)

		// partial page
		checkTagSubset(t, services.Pagination{PageSize: 100, Page: 1}, db)

		// emptySet
		checkTagSubset(t, services.Pagination{PageSize: 100, Page: 2}, db)
	})
}

func checkTagSubset(t *testing.T, p services.Pagination, db *database.Connection) {
	opID := OpChamberOfSecrets.ID

	fullSet := getTagFromOperationID(t, db, opID)

	var tags []models.Tag
	err := p.Select(context.Background(), db, &tags, sq.Select("*").
		From("tags").
		Where(sq.Eq{"operation_id": opID}))
	require.NoError(t, err)

	offset := (p.Page - 1) * p.PageSize
	limit := helpers.Clamp(int64(len(fullSet))-offset, 0, p.PageSize)

	require.Equal(t, limit, int64(len(tags)))

	for i := int64(0); i < limit; i++ {
		require.Equal(t, fullSet[i+offset], tags[i])
	}
	require.Equal(t, int64(len(fullSet)), p.TotalCount)
}
