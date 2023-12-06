package services

import (
	"context"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/server/dissectors"

	sq "github.com/Masterminds/squirrel"
)

type Pagination struct {
	PageSize    int64
	Page        int64
	maxPageSize int64
	TotalCount  int64
}

const defaultMaxPageSize = 250

// ParseRequestQueryPagination retreives the part of the request set aside for pagination
// Note that this retrieves the values and hopes for the best. Since this uses a DissectedRequest,
// it is the caller of the function to ensure no error occurred _after_ this has been called.
func ParseRequestQueryPagination(dr dissectors.DissectedRequest, defaultMaxItems int64) Pagination {
	return Pagination{
		Page:     dr.FromQuery("page").OrDefault(1).AsInt64(),
		PageSize: dr.FromQuery("pageSize").OrDefault(defaultMaxItems).AsInt64(),
	}
}

// SetMaxItems sets the maximum number of items that can be returned in a request/page. This
// must be called before Select to have any effect
func (p *Pagination) SetMaxItems(maxItems int64) *Pagination {
	p.maxPageSize = maxItems
	return p
}

// constrain restricts the page and pageSize fields to reasonable values (i.e. page >=1, 0 < maxitems <= defaultMaxItems (currently defaultMaxPageSize))
func (p *Pagination) constrain() {
	if p.Page < 1 {
		p.Page = 1
	}
	maxItems := p.maxPageSize
	if maxItems == 0 {
		maxItems = defaultMaxPageSize
	}
	p.PageSize = helpers.Clamp(p.PageSize, 1, maxItems)
}

// WrapData is a small helper to turn the desired content of a request into a pagination result set
func (p *Pagination) WrapData(data interface{}) *dtos.PaginationWrapper {
	pageQuotient := p.TotalCount / p.PageSize
	totalPages := pageQuotient + helpers.Clamp(p.TotalCount%p.PageSize, 0, 1)

	return &dtos.PaginationWrapper{
		PageNumber: p.Page,
		PageSize:   p.PageSize,
		Content:    data,
		TotalCount: p.TotalCount,
		TotalPages: totalPages,
	}
}

// Select is a wrapper around database.Connection.Select. This performs a query that returns multiple rows.
// In addition, this counts the total number of rows matching this query, and saves the result inside the
// pagination structure.
//
// This actually performs two queries: the intended query, plus a second query to discover the total
// number of matching rows. I think this works differently in other databases, but this
// seems to be the preferred route for mysql. See: https://dev.mysql.com/doc/refman/8.0/en/information-functions.html#function_found-rows
// for more details
//
// Note: It is possible to have the initial query succeed and the count query to fail. In order to prevent
// odd issues, you should always do an error check before using the resulting value.
//
// Note 2: This is really only useful for communicating size back to the enduser. For other pagination techniques,
// you may want to use LIMIT and OFFSET directly
func (p *Pagination) Select(ctx context.Context, db *database.Connection, resultSlice interface{}, sb sq.SelectBuilder) error {
	limitedSelect := sb.Limit(uint64(p.PageSize)).Offset(uint64(p.PageSize * (p.Page - 1)))

	var count int64
	err := db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Select(resultSlice, limitedSelect)
		tx.Get(&count, sq.Select("count(*)").FromSelect(sb, "T"))
	})
	if err != nil {
		return backend.WrapError("Unable to get pagination count", err)
	}

	p.TotalCount = count
	return nil
}
