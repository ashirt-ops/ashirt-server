// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"testing"
	"time"

	"github.com/theparanoids/ashirt-server/backend/helpers"

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

	uuid := "a"
	uuidBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{UUID: uuid})
	require.Equal(t, " WHERE evidence.operation_id = ? AND evidence.uuid = ?", toWhere(uuidBuilder))
	require.Equal(t, []interface{}{opID, uuid}, toWhereValues(uuidBuilder))

	text := []string{"one", "two"}
	descBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{Text: text})
	require.Equal(t, " WHERE evidence.operation_id = ? AND description LIKE ? AND description LIKE ?", toWhere(descBuilder))
	require.Equal(t, []interface{}{opID, "%" + text[0] + "%", "%" + text[1] + "%"}, toWhereValues(descBuilder))

	start, end := time.Now(), time.Now().Add(5*time.Second)
	dateBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{DateRange: &helpers.DateRange{From: start, To: end}})
	require.Equal(t, " WHERE evidence.operation_id = ? AND evidence.occurred_at >= ? AND evidence.occurred_at <= ?", toWhere(dateBuilder))
	require.Equal(t, []interface{}{opID, start, end}, toWhereValues(dateBuilder))

	operator := "Johnny 5"
	operatorBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{Operator: operator})
	require.Equal(t, " WHERE evidence.operation_id = ? AND "+eviForOpOperatorWhereComponent, toWhere(operatorBuilder))
	require.Equal(t, []interface{}{opID, operator}, toWhereValues(operatorBuilder))

	tags := []string{"alpha", "beta", "gamma"}
	tagBuilder := buildListEvidenceWhereClause(base, opID, helpers.TimelineFilters{Tags: tags})
	require.Equal(t, " WHERE evidence.operation_id = ? AND "+eviForOpTagWhereComponent, toWhere(tagBuilder))
	require.Equal(t, []interface{}{opID, tags, len(tags)}, toWhereValues(tagBuilder))
}
