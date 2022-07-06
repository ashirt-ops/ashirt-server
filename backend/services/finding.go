// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/helpers/filter"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"golang.org/x/sync/errgroup"

	sq "github.com/Masterminds/squirrel"
)

type AddEvidenceToFindingInput struct {
	OperationSlug    string
	FindingUUID      string
	EvidenceToAdd    []string
	EvidenceToRemove []string
}

type CreateFindingInput struct {
	OperationSlug string
	Category      string
	Title         string
	Description   string
}

type DeleteFindingInput struct {
	OperationSlug string
	FindingUUID   string
}

type ListFindingsForOperationInput struct {
	OperationSlug string
	Filters       helpers.TimelineFilters
}

type ReadFindingInput struct {
	OperationSlug string
	FindingUUID   string
}

type UpdateFindingInput struct {
	OperationSlug string
	FindingUUID   string
	Category      string
	Title         string
	Description   string
	TicketLink    *string
	ReadyToReport bool
}

func CreateFinding(ctx context.Context, db *database.Connection, i CreateFindingInput) (*dtos.Finding, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to create finding", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unable to create finding", backend.UnauthorizedWriteErr(err))
	}

	if i.Title == "" {
		return nil, backend.MissingValueErr("Title")
	}

	if i.Category == "" {
		return nil, backend.MissingValueErr("Category")
	}

	useCategoryID, err := getFindingCategoryID(i.Category, db.Select)

	if err != nil {
		return nil, backend.WrapError("Unable create finding", err)
	}
	if useCategoryID == nil {
		return nil, backend.BadInputErr(errors.New("no such category"), "Unknown Category")
	}

	findingUUID := uuid.New().String()
	_, err = db.Insert("findings", map[string]interface{}{
		"uuid":         findingUUID,
		"operation_id": operation.ID,
		"category_id":  useCategoryID,
		"title":        i.Title,
		"description":  i.Description,
	})
	if err != nil {
		return nil, backend.WrapError("Unable to insert finding", backend.DatabaseErr(err))
	}

	return &dtos.Finding{
		UUID:        findingUUID,
		Title:       i.Title,
		Description: i.Description,
	}, nil
}

func DeleteFinding(ctx context.Context, db *database.Connection, i DeleteFindingInput) error {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return backend.WrapError("Unable to delete finding", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to delete finding", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"finding_id": finding.ID}))
		tx.Delete(sq.Delete("findings").Where(sq.Eq{"id": finding.ID}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete finding", backend.DatabaseErr(err))
	}

	return nil
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

func ReadFinding(ctx context.Context, db *database.Connection, i ReadFindingInput) (*dtos.Finding, error) {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return nil, backend.WrapError("Unable to read finding", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to read finding", backend.UnauthorizedReadErr(err))
	}

	var evidenceIDs []int64

	err = db.Select(&evidenceIDs, sq.Select("evidence_id").
		From("evidence_finding_map").
		Where(sq.Eq{"finding_id": finding.ID}))
	if err != nil {
		return nil, backend.WrapError("Cannot load evidence for finding", backend.DatabaseErr(err))
	}

	_, allTags, err := tagsForEvidenceByID(db, evidenceIDs)
	if err != nil {
		return nil, backend.WrapError("Cannot load tags for evidence", backend.DatabaseErr(err))
	}

	var realCategory = ""
	if finding.CategoryID != nil {
		realCategory, err = getFindingCategory(db, *finding.CategoryID)
		if err != nil {
			return nil, backend.WrapError("Cannot load finding category for finding", backend.DatabaseErr(err))
		}
	}

	return &dtos.Finding{
		UUID:          i.FindingUUID,
		Title:         finding.Title,
		Category:      realCategory,
		Description:   finding.Description,
		NumEvidence:   len(evidenceIDs),
		Tags:          allTags,
		ReadyToReport: finding.ReadyToReport,
		TicketLink:    finding.TicketLink,
	}, nil
}

func UpdateFinding(ctx context.Context, db *database.Connection, i UpdateFindingInput) error {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return backend.WrapError("Unable to lookup operation", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Failed permission check", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		useCategoryID, _ := getFindingCategoryID(i.Category, tx.Select)

		tx.Update(sq.Update("findings").
			SetMap(map[string]interface{}{
				"category_id":     useCategoryID,
				"title":           i.Title,
				"description":     i.Description,
				"ticket_link":     i.TicketLink,
				"ready_to_report": i.ReadyToReport,
			}).
			Where(sq.Eq{"id": finding.ID}))
	})

	if err != nil {
		return backend.WrapError("Unable to update database", backend.UnauthorizedWriteErr(err))
	}
	return nil
}

func AddEvidenceToFinding(ctx context.Context, db *database.Connection, i AddEvidenceToFindingInput) error {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return backend.WrapError("Unable to add evidence to finding", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyFindingsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unable to add evidence to finding", backend.UnauthorizedWriteErr(err))
	}

	var g errgroup.Group
	g.Go(func() (err error) { return batchAddEvidenceToFinding(db, i.EvidenceToAdd, operation.ID, finding.ID) })
	g.Go(func() (err error) { return batchRemoveEvidenceFromFinding(db, i.EvidenceToRemove, finding.ID) })
	if err = g.Wait(); err != nil {
		return backend.WrapError("Unable to add evidence to finding", backend.UnauthorizedWriteErr(err))
	}

	return nil
}

func buildQueryForEvidenceFromUUIDs(evidenceUUIDs []string) sq.SelectBuilder {
	return sq.Select("*").
		From("evidence").
		Where(sq.Eq{"uuid": evidenceUUIDs})
}

func batchAddEvidenceToFinding(db *database.Connection, evidenceUUIDs []string, operationID int64, findingID int64) error {
	if len(evidenceUUIDs) == 0 {
		return nil
	}
	var evidence []models.Evidence
	if err := db.Select(&evidence, buildQueryForEvidenceFromUUIDs(evidenceUUIDs)); err != nil {
		return backend.WrapError("Unable to get evidence from uuids", err)
	}
	evidenceIDs := []int64{}
	for _, evi := range evidence {
		if evi.OperationID != operationID {
			return fmt.Errorf(
				"Cannot add evidence %d to operation %d. Evidence belongs to operation %d",
				evi.ID, operationID, evi.OperationID,
			)
		}
		evidenceIDs = append(evidenceIDs, evi.ID)
	}
	return db.BatchInsert("evidence_finding_map", len(evidenceIDs), func(idx int) map[string]interface{} {
		return map[string]interface{}{
			"finding_id":  findingID,
			"evidence_id": evidenceIDs[idx],
		}
	})
}

func batchRemoveEvidenceFromFinding(db *database.Connection, evidenceUUIDs []string, findingID int64) error {
	if len(evidenceUUIDs) == 0 {
		return nil
	}
	var evidence []models.Evidence
	if err := db.Select(&evidence, buildQueryForEvidenceFromUUIDs(evidenceUUIDs)); err != nil {
		return backend.WrapError("Unable to get evidence from uuids", err)
	}
	evidenceIDs := []int64{}
	for _, evi := range evidence {
		evidenceIDs = append(evidenceIDs, evi.ID)
	}

	return db.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"finding_id": findingID, "evidence_id": evidenceIDs}))
}

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

	if values := filters.DateRanges; len(values) > 0 {
		// we're only going to support a single date range for now TODO
		dateFilter := values[0]
		include := !(dateFilter.Modifier == filter.Not)

		queryFilters = append(queryFilters, findingDateRangeWhere(include))
		queryValues = append(queryValues, dateFilter.Value.From, dateFilter.Value.To)
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

func findingDateRangeWhere(in bool) string {
	return "findings.id " + inOrNotIn(in) + " (" +
		"  SELECT findings.id FROM findings" +
		"  INNER JOIN evidence_finding_map ON evidence_finding_map.finding_id = findings.id" +
		"  INNER JOIN evidence ON evidence.id = evidence_finding_map.evidence_id" +
		"  GROUP BY findings.id HAVING MAX(evidence.occurred_at) >= ? AND MIN(evidence.occurred_at) <= ?" +
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
