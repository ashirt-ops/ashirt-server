// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"

	sq "github.com/Masterminds/squirrel"
)

type detailedSchemeTable struct {
	AuthScheme      string     `db:"auth_scheme"`
	UserCount       int64      `db:"num_users"`
	UniqueUserCount int64      `db:"unique_users"`
	LastUsed        *time.Time `db:"last_used"`
}

func ListAuthDetails(ctx context.Context, db *database.Connection, supportedAuthSchemes *[]dtos.SupportedAuthScheme) ([]*dtos.DetailedAuthenticationInfo, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var detailedAuthData []detailedSchemeTable
	err := db.Select(&detailedAuthData,
		sq.Select("auth_scheme",
			"COUNT(*) AS num_users",
			"COALESCE(SUM(is_unique), 0) AS unique_users",
			"MAX(last_login) AS last_used").
			From("auth_scheme_data").
			LeftJoin(
				"(SELECT user_id, 1 AS is_unique FROM auth_scheme_data WHERE auth_scheme != 'recovery' GROUP BY user_id HAVING COUNT(*) = 1) AS t "+
					"ON t.user_id = auth_scheme_data.user_id",
			).
			Where(sq.NotEq{"auth_scheme": "recovery"}).
			GroupBy("auth_scheme"))

	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	return mergeSchemes(detailedAuthData, supportedAuthSchemes), nil
}

// mergeSchemes cobbles together the list of known supported schemes (whether used or not), and the actual
// schemes used (whether supported or not). The result here would be a list with schemes that _can be_
// used, schemes that _are_ used, and schemes that _were previously_ used.
func mergeSchemes(foundSchemes []detailedSchemeTable, supportedAuthSchemes *[]dtos.SupportedAuthScheme) []*dtos.DetailedAuthenticationInfo {
	clonedSchemes := make([]dtos.SupportedAuthScheme, len(*supportedAuthSchemes))
	copy(clonedSchemes, *supportedAuthSchemes)

	// create space for all possible elements
	schemes := make([]*dtos.DetailedAuthenticationInfo, 0, len(foundSchemes)+len(clonedSchemes))

	// Add schemes that are used (whether supported or not)
	for i, scheme := range foundSchemes {
		schemes = append(schemes, &dtos.DetailedAuthenticationInfo{ // pre-populate known values
			AuthSchemeCode:  scheme.AuthScheme,
			UserCount:       scheme.UserCount,
			UniqueUserCount: scheme.UniqueUserCount,
			LastUsed:        scheme.LastUsed,
			Labels:          []string{},
		})
		matchingSchemeIndex := getMatchingSchemeIndex(&clonedSchemes, scheme.AuthScheme)
		if matchingSchemeIndex == -1 {
			schemes[i].AuthSchemeName = scheme.AuthScheme
			schemes[i].Labels = append(schemes[i].Labels, "Unsupported")
		} else {
			schemes[i].AuthSchemeName = clonedSchemes[matchingSchemeIndex].SchemeName

			// Remove the used element (swap + remove last)
			clonedSchemes[matchingSchemeIndex], clonedSchemes[len(clonedSchemes)-1] = clonedSchemes[len(clonedSchemes)-1], clonedSchemes[matchingSchemeIndex]
			clonedSchemes = clonedSchemes[:len(clonedSchemes)-1]
		}
	}

	// Add schemes that are supported (whether used or not)
	for _, scheme := range clonedSchemes {
		schemes = append(schemes, &dtos.DetailedAuthenticationInfo{
			AuthSchemeName: scheme.SchemeName,
			AuthSchemeCode: scheme.SchemeCode,
			Labels:         []string{},
		})
	}

	return schemes
}

func getMatchingSchemeIndex(allSchemes *[]dtos.SupportedAuthScheme, targetSchemeCode string) int {
	if allSchemes == nil {
		return -1
	}
	for i, scheme := range *allSchemes {
		if scheme.SchemeCode == targetSchemeCode {
			return i
		}
	}
	return -1
}
