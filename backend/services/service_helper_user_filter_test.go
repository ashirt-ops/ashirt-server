// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/server/remux"
	"github.com/theparanoids/ashirt-server/backend/services"

	sq "github.com/Masterminds/squirrel"
)

func TestParseRequestQueryUserFilter(t *testing.T) {
	testParseRequestQueryUserFilter(t, `/users?name`, []string{})
	testParseRequestQueryUserFilter(t, `/users?name=ron`, []string{"ron"})
	testParseRequestQueryUserFilter(t, `/users?name=r%20w`, []string{"r", "w"})
}

func testParseRequestQueryUserFilter(t *testing.T, endpoint string, expectedContent []string) {
	r := httptest.NewRequest("POST", endpoint, nil)
	dr := remux.DissectJSONRequest(r)

	filter := services.ParseRequestQueryUserFilter(dr)

	require.Nil(t, dr.Error)
	require.Equal(t, expectedContent, filter.NameParts)
}

func TestAddWhere(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		filter := services.UserFilter{
			NameParts:  []string{"Ron"},
			UsersTable: "u",
		}

		query := sq.Select("count(*)").From("users AS u")

		// verify that the base db query works
		var userCount int64
		err := db.Get(&userCount, query)
		require.NoError(t, err)
		require.Equal(t, int64(len(seed.Users)), userCount)

		// verify that we have a different count after applying the filter
		expectedUserSet := filterUsersManually(filter.NameParts[0], seed.Users)
		userCount = -1
		filter.AddWhere(&query)
		err = db.Get(&userCount, query)
		require.NoError(t, err)
		require.Equal(t, int64(len(expectedUserSet)), userCount)
	})
}

func filterUsersManually(query string, userList []models.User) []models.User {
	rtn := make([]models.User, 0, len(userList))
	query = strings.ToLower(query)
	for _, user := range userList {
		userName := strings.ToLower(user.FirstName + " " + user.LastName)
		if strings.Contains(userName, query) {
			rtn = append(rtn, user)
		}
	}

	return rtn
}

func TestFilterUserManuallyTestHelper(t *testing.T) {
	users := []models.User{UserHarry, UserRon, UserHermione}

	require.Equal(t, []models.User{UserRon}, filterUsersManually("ron", users))
}
