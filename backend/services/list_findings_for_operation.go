// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/helpers/filter"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type ListFindingsForOperationInput struct {
	OperationSlug string
	Filters       helpers.TimelineFilters
}

func ListFindingsForOperation(ctx context.Context, db *database.Connection, i ListFindingsForOperationInput) ([]*dtos.Finding, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list findings for operation", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list findings for operation", backend.UnauthorizedReadErr(err))
	}

	whereClause, whereValues := buildListFindingsWhereClause(operation.ID, i.Filters)
	var findings []struct {
		models.Finding
		NumEvidence     int        `db:"num_evidence"`
		OccurredFrom    *time.Time `db:"occurred_from"`
		OccurredTo      *time.Time `db:"occurred_to"`
		TagIDs          *string    `db:"tag_ids"`
		FindingCategory *string    `db:"finding_category"`
	}

	sb := sq.Select(
		"findings.*",
		"COUNT(DISTINCT evidence_finding_map.evidence_id) AS num_evidence",
		"MIN(evidence.occurred_at) AS occurred_from",
		"MAX(evidence.occurred_at) AS occurred_to",
		"GROUP_CONCAT(DISTINCT tag_id) AS tag_ids",
		"finding_categories.category AS finding_category").
		From("findings").
		LeftJoin("evidence_finding_map ON findings.id = finding_id").
		LeftJoin("evidence ON evidence_id = evidence.id").
		LeftJoin("tag_evidence_map ON tag_evidence_map.evidence_id = evidence_finding_map.evidence_id").
		LeftJoin("finding_categories ON finding_categories.id = findings.category_id").
		Where(whereClause, whereValues...).
		GroupBy("findings.id")

	if i.Filters.SortAsc {
		sb = sb.OrderBy("occurred_to ASC").
			OrderBy("occurred_from ASC")
	} else {
		sb = sb.OrderBy("occurred_to DESC").
			OrderBy("occurred_from DESC")
	}

	err = db.Select(&findings, sb)
	if err != nil {
		return nil, backend.WrapError("Cannot list findings for operation", backend.DatabaseErr(err))
	}

	if len(findings) == 0 {
		return []*dtos.Finding{}, nil
	}

	tagsByID, err := allTagsByID(db)
	if err != nil {
		return nil, backend.WrapError("Cannot find all tags", backend.DatabaseErr(err))
	}

	findingsDTO := make([]*dtos.Finding, len(findings))
	for idx, finding := range findings {
		realCategory := ""
		if finding.FindingCategory != nil {
			realCategory = *finding.FindingCategory
		}
		findingsDTO[idx] = &dtos.Finding{
			UUID:          finding.UUID,
			Category:      realCategory,
			Title:         finding.Title,
			Description:   finding.Description,
			OccurredFrom:  finding.OccurredFrom,
			OccurredTo:    finding.OccurredTo,
			NumEvidence:   finding.NumEvidence,
			ReadyToReport: finding.ReadyToReport,
			TicketLink:    finding.TicketLink,
			Tags:          buildTags(tagsByID, finding.TagIDs),
		}
	}

	return findingsDTO, nil
}

const findingsDateRangeWhereComponent = "findings.id IN (" +
	"  SELECT findings.id FROM findings" +
	"  INNER JOIN evidence_finding_map ON evidence_finding_map.finding_id = findings.id" +
	"  INNER JOIN evidence ON evidence.id = evidence_finding_map.evidence_id" +
	"  GROUP BY findings.id HAVING MAX(evidence.occurred_at) >= ? AND MIN(evidence.occurred_at) <= ?" +
	")"

const findingsTextWhereComponent = "(findings.title LIKE ? OR findings.description LIKE ?)"
const findingsOperationIDWhereComponent = "findings.operation_id = ?"

func buildListFindingsWhereClause(operationID int64, filters helpers.TimelineFilters) (string, []interface{}) {
	queryFilters := []string{findingsOperationIDWhereComponent}
	queryValues := []interface{}{operationID}

	addWhere := func(vals filter.Values, whereFunc func(bool) string) {
		findingAddWhereAndNot(&queryFilters, &queryValues, vals, whereFunc)
	}

	if len(filters.UUID) > 0 {
		addWhere(filters.UUID, findingUUIDWhere)
	}

	if len(filters.Tags) > 0 {
		addWhere(filters.Tags, findingTagOrWhere)
	}

	for _, text := range filters.Text {
		fuzzyText := "%" + text + "%"
		queryFilters = append(queryFilters, findingsTextWhereComponent)
		queryValues = append(queryValues, fuzzyText, fuzzyText)
	}

	if len(filters.DateRanges) > 0 {
		// we're only going to support a single date range for now TODO
		queryFilters = append(queryFilters, findingsDateRangeWhereComponent)
		queryValues = append(queryValues, filters.DateRanges[0].From, filters.DateRanges[0].To)
	}

	if len(filters.Operator) > 0 {
		addWhere(filters.Operator, findingOperatorWhere)
	}

	if len(filters.WithEvidenceUUID) > 0 {
		addWhere(filters.WithEvidenceUUID, findingEvidenceUUIDWhere)
	}

	return strings.Join(queryFilters, " AND "), queryValues
}

func buildTags(tagsByID map[int64]dtos.Tag, tagIDs *string) []dtos.Tag {
	tags := []dtos.Tag{}
	if tagIDs == nil {
		return tags
	}
	for _, tagIDStr := range strings.Split(*tagIDs, ",") {
		tagID, _ := strconv.ParseInt(tagIDStr, 10, 64)
		tags = append(tags, tagsByID[tagID])
	}
	return tags
}

func allTagsByID(db *database.Connection) (map[int64]dtos.Tag, error) {
	tagsByID := map[int64]dtos.Tag{}

	var tags []models.Tag
	err := db.Select(&tags, sq.Select("id", "name", "color_name").From("tags"))
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		tagsByID[tag.ID] = dtos.Tag{
			ID:        tag.ID,
			Name:      tag.Name,
			ColorName: tag.ColorName,
		}
	}
	return tagsByID, nil
}

func findingAddWhereAndNot(queryFilters *[]string, queryValues *[]interface{}, vals filter.Values, whereFunc func(bool) string) {
	splitValues := vals.SplitByModifier()

	if values := splitValues[filter.Normal]; len(values) > 0 {
		*queryFilters = append(*queryFilters, whereFunc(true))
		*queryValues = append(*queryValues, values)
	}
	if values := splitValues[filter.Not]; len(values) > 0 {
		*queryFilters = append(*queryFilters, whereFunc(false))
		*queryValues = append(*queryValues, values)
	}
}

func inOrNotIn(in bool) string {
	if in {
		return "IN"
	}
	return "NOT IN"
}

func findingUUIDWhere(in bool) string {
	return "findings.uuid " + inOrNotIn(in) + " (?)"
}

func findingOperatorWhere(in bool) string {
	return "findings.id " + inOrNotIn(in) + " (" +
		"  SELECT findings.id FROM findings" +
		"  INNER JOIN evidence_finding_map ON evidence_finding_map.finding_id = findings.id" +
		"  INNER JOIN evidence ON evidence.id = evidence_finding_map.evidence_id" +
		"  LEFT JOIN users ON users.id = evidence.operator_id" +
		"  WHERE users.slug IN(?)" +
		")"
}

func findingEvidenceUUIDWhere(in bool) string {
	return "findings.id " + inOrNotIn(in) + " (" +
		"  SELECT finding_id FROM evidence_finding_map" +
		"  LEFT JOIN evidence ON evidence.id = evidence_finding_map.evidence_id" +
		"  WHERE evidence.uuid IN (?)" +
		")"
}

// func findingTagAndWhere(is bool) string {
// 	return findingTagWhere(is, true)
// }

func findingTagOrWhere(in bool) string {
	return findingTagWhere(in, false)
}

func findingTagWhere(in, all bool) string {
	groupBy := ""
	if all {
		groupBy = "  GROUP BY findings.id HAVING COUNT(DISTINCT tags.id) = ?"
	}
	return "findings.id " + inOrNotIn(in) + " (" +
		"  SELECT findings.id FROM findings" +
		"  INNER JOIN evidence_finding_map ON evidence_finding_map.finding_id = findings.id" +
		"  INNER JOIN tag_evidence_map ON tag_evidence_map.evidence_id = evidence_finding_map.evidence_id" +
		"  LEFT JOIN tags ON tags.id = tag_evidence_map.tag_id" +
		"  WHERE tags.name IN (?)" +
		groupBy +
		")"
}
