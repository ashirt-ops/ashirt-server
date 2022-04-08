// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"testing"
	"time"

	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/helpers/filter"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

// TestBuildListEvidenceWhereClause is a unit-test for the buildTags function.
func TestBuildListEvidenceWhereClause(t *testing.T) {
	base := sq.Select("a").From("b")
	baseSQL, _, _ := base.ToSql()
	opID := int64(1234)

	toWhere := func(s sq.SelectBuilder) string {
		query, _, _ := s.ToSql()
		return query[len(baseSQL):]
	}
	toWhereValues := func(s sq.SelectBuilder) []interface{} {
		_, v, _ := s.ToSql()
		return v
	}

	noFilterBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{})
	require.Equal(t, " WHERE evidence.operation_id = ?", toWhere(noFilterBuilder))
	require.Equal(t, []interface{}{opID}, toWhereValues(noFilterBuilder))

	uuids := filter.Values{filter.Val("a")}
	uuidBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{UUID: uuids})
	require.Equal(t, " WHERE evidence.operation_id = ? AND evidence.uuid IN (?)", toWhere(uuidBuilder))
	require.Equal(t, []interface{}{opID, uuids.Values()}, toWhereValues(uuidBuilder))

	text := []string{"one", "two"}
	descBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{Text: text})
	require.Equal(t, " WHERE evidence.operation_id = ? AND description LIKE ? AND description LIKE ?", toWhere(descBuilder))
	require.Equal(t, []interface{}{opID, "%" + text[0] + "%", "%" + text[1] + "%"}, toWhereValues(descBuilder))

	start, end := time.Now(), time.Now().Add(5*time.Second)
	singleDate := filter.DateValues{
		filter.DateVal(filter.DateRange{From: start, To: end}),
	}
	datePart := "(evidence.occurred_at >= ? AND evidence.occurred_at <= ?)"
	singleDateBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{DateRanges: singleDate})
	require.Equal(t, " WHERE evidence.operation_id = ? AND ("+datePart+")", toWhere(singleDateBuilder))
	require.Equal(t, []interface{}{opID, start, end}, toWhereValues(singleDateBuilder))

	start2, end2 := time.Now(), time.Now().Add(5*time.Second)
	dates := filter.DateValues{
		filter.DateVal(filter.DateRange{From: start, To: end}),
		filter.DateVal(filter.DateRange{From: start2, To: end2}),
	}
	multiDateBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{DateRanges: dates})
	require.Equal(t, " WHERE evidence.operation_id = ? AND ("+datePart+" OR "+datePart+")", toWhere(multiDateBuilder))
	require.Equal(t, []interface{}{opID, start, end, start2, end2}, toWhereValues(multiDateBuilder))

	operators := filter.Values{filter.Val("Johnny 5")}
	operatorBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{Operator: operators})
	require.Equal(t, " WHERE evidence.operation_id = ? AND "+evidenceOperatorWhere(true), toWhere(operatorBuilder))
	require.Equal(t, []interface{}{opID, operators.Values()}, toWhereValues(operatorBuilder))

	tags := filter.Values{filter.Val("alpha"), filter.Val("beta"), filter.Val("gamma")}
	tagBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{Tags: tags})
	require.Equal(t, " WHERE evidence.operation_id = ? AND "+evidenceTagOrWhere(true), toWhere(tagBuilder))
	require.Equal(t, []interface{}{opID, tags.Values()}, toWhereValues(tagBuilder))
}
