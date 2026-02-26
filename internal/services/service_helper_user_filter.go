package services

import (
	"strings"

	"github.com/ashirt-ops/ashirt-server/internal/server/dissectors"

	sq "github.com/Masterminds/squirrel"
)

// UserFilter provides a mechanism to alter queries such that users are filtered
type UserFilter struct {
	NameParts  []string
	UsersTable string
}

// ParseRequestQueryUserFilter generates a UserFilter object from a given request.
// This expects that filtering is specified by the query parameter "name"
func ParseRequestQueryUserFilter(dr dissectors.DissectedRequest) UserFilter {
	return UserFilter{
		NameParts:  strings.Fields(dr.FromQuery("name").OrDefault("").AsString()),
		UsersTable: "users",
	}
}

// AddWhere adds to the given SelectBuilder a Where clause that will apply the filtering
func (uf *UserFilter) AddWhere(sb *sq.SelectBuilder) {
	if len(uf.NameParts) > 0 {
		baseQuery := "concat(" + uf.UsersTable + ".first_name, ' ', " + uf.UsersTable + ".last_name)"
		*sb = sb.Where(sq.Like{baseQuery: "%" + strings.Join(uf.NameParts, "%") + "%"})
	}

}
