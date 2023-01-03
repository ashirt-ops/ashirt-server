// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"strings"

	"github.com/theparanoids/ashirt-server/backend/server/dissectors"

	sq "github.com/Masterminds/squirrel"
)

// UserFilter provides a mechanism to alter queries such that users are filtered
type UserGroupFilter struct {
	NameParts       []string
	UserGroupsTable string
}

// ParseRequestQueryUserFilter generates a UserFilter object from a given request.
// This expects that filtering is specified by the query parameter "name"
func ParseRequestQueryUserGroupFilter(dr dissectors.DissectedRequest) UserGroupFilter {
	return UserGroupFilter{
		NameParts:       strings.Fields(dr.FromQuery("name").OrDefault("").AsString()),
		UserGroupsTable: "user_groups",
	}
}

// AddWhere adds to the given SelectBuilder a Where clause that will apply the filtering
func (uf *UserGroupFilter) AddWhere(sb *sq.SelectBuilder) {
	if len(uf.NameParts) > 0 {
		*sb = sb.Where(sq.Like{"name": "%" + strings.Join(uf.NameParts, "%") + "%"})
	}
}
