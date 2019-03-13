// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/theparanoids/ashirt/backend/helpers"

	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/stretchr/testify/require"
)

// TestBuildTags is a unit-test suite for the buildTags function.
func TestBuildTags(t *testing.T) {
	tagsByID := map[int64]dtos.Tag{
		1: dtos.Tag{ID: 1, ColorName: "blue", Name: "blueTag"},
		2: dtos.Tag{ID: 2, ColorName: "red", Name: "redTag"},
		3: dtos.Tag{ID: 3, ColorName: "aqua", Name: "aquaTag"},
		4: dtos.Tag{ID: 4, ColorName: "maroon", Name: "maroonTag"},
	}
	tagStr := "1,2,4"
	tags := buildTags(tagsByID, &tagStr)
	sort.Slice(tags, func(a, b int) bool { return tags[a].ID < tags[b].ID })

	require.Equal(t, []dtos.Tag{tagsByID[1], tagsByID[2], tagsByID[4]}, tags)

	require.Equal(t, []dtos.Tag{}, buildTags(tagsByID, nil))
}

// TestBuildListFindingsWhereClause is a unit-test for the buildTags function.
func TestBuildListFindingsWhereClause(t *testing.T) {
	test := func(filters helpers.TimelineFilters, queryParts []string, queryValues []interface{}) {
		targetQuery := strings.Join(append([]string{findingsOperationIDWhereComponent}, queryParts...), " AND ")
		query, values := buildListFindingsWhereClause(1, filters)
		require.Equal(t, targetQuery, query)
		require.Equal(t, append([]interface{}{int64(1)}, queryValues...), values)
	}

	test(helpers.TimelineFilters{}, []string{}, []interface{}{}) // no filters test
	test(helpers.TimelineFilters{UUID: "abc"}, []string{findingsUUIDWhereComponent}, []interface{}{"abc"})
	test(helpers.TimelineFilters{Tags: []string{"fraggle", "rock"}}, []string{findingsTagWhereComponent}, []interface{}{[]string{"fraggle", "rock"}, 2})
	test(helpers.TimelineFilters{Text: []string{"some", "text"}}, []string{findingsTextWhereComponent, findingsTextWhereComponent}, []interface{}{"%some%", "%some%", "%text%", "%text%"})

	start, end := time.Now(), time.Now().Add(5*time.Second)
	test(helpers.TimelineFilters{DateRange: &helpers.DateRange{From: start, To: end}}, []string{findingsDateRangeWhereComponent}, []interface{}{start, end})
	test(helpers.TimelineFilters{Operator: "MyOp"}, []string{findingsOperatorWhereComponent}, []interface{}{"MyOp"})
	test(helpers.TimelineFilters{WithEvidenceUUID: "abc"}, []string{findingsEvidenceUUIDWhereComponent}, []interface{}{"abc"})
}

// TestAllTagsByID is a unit-test suite for the allTagsByID function.
func TestAllTagsByID(t *testing.T) {
	db := internalTestDBSetup(t)

	tags := []models.Tag{
		models.Tag{ID: 1, OperationID: 1, Name: "firstTag", ColorName: "black", CreatedAt: time.Now()},
		models.Tag{ID: 2, OperationID: 1, Name: "secondTag", ColorName: "white", CreatedAt: time.Now()},
		models.Tag{ID: 3, OperationID: 1, Name: "thirdTag", ColorName: "gray", CreatedAt: time.Now()},
		models.Tag{ID: 7, OperationID: 1, Name: "fourthTag", ColorName: "red", CreatedAt: time.Now()},
	}
	db.BatchInsert("tags", len(tags), func(i int) map[string]interface{} {
		return map[string]interface{}{
			"id":           tags[i].ID,
			"operation_id": tags[i].OperationID,
			"name":         tags[i].Name,
			"color_name":   tags[i].ColorName,
			"created_at":   tags[i].CreatedAt,
		}
	})

	foundTags, err := allTagsByID(db)
	require.NoError(t, err)
	for _, tag := range tags {
		require.Equal(t, tag.ID, foundTags[tag.ID].ID)
		require.Equal(t, tag.Name, foundTags[tag.ID].Name)
		require.Equal(t, tag.ColorName, foundTags[tag.ID].ColorName)
	}
}
