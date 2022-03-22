// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"

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

type ListEvidenceForOperationInput struct {
	OperationSlug string
	Filters       helpers.TimelineFilters
}

// ListEvidenceForOperation retrieves all evidence for a particular operation id matching a particular
// set of filters (e.g. tag:some_tag)
func ListEvidenceForOperation(ctx context.Context, db *database.Connection, i ListEvidenceForOperationInput) ([]*dtos.Evidence, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list evidence for an operation", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list evidence for an operation", backend.UnauthorizedReadErr(err))
	}

	var evidence []struct {
		models.Evidence
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Slug      string `db:"slug"`
	}

	sb := sq.Select("evidence.id", "evidence.uuid", "description", "evidence.content_type", "occurred_at", "users.first_name", "users.last_name", "users.slug").
		From("evidence").
		LeftJoin("users ON evidence.operator_id = users.id")

	if i.Filters.SortAsc {
		sb = sb.OrderBy("occurred_at ASC")
	} else {
		sb = sb.OrderBy("occurred_at DESC")
	}

	sb = buildListEvidenceWhereClause(sb, operation.ID, i.Filters)

	err = db.Select(&evidence, sb)
	if err != nil {
		return nil, backend.WrapError("Cannot list evidence for an operation", backend.DatabaseErr(err))
	}

	if len(evidence) == 0 {
		return []*dtos.Evidence{}, nil
	}

	evidenceIDs := make([]int64, len(evidence))
	for idx, ev := range evidence {
		evidenceIDs[idx] = ev.ID
	}

	tagsByEvidenceID, _, err := tagsForEvidenceByID(db, evidenceIDs)
	if err != nil {
		return nil, backend.WrapError("Cannot get tags for evidence", backend.DatabaseErr(err))
	}

	evidenceDTO := make([]*dtos.Evidence, len(evidence))
	for idx, evi := range evidence {
		tags, ok := tagsByEvidenceID[evi.ID]

		if !ok {
			tags = []dtos.Tag{}
		}

		evidenceDTO[idx] = &dtos.Evidence{
			UUID:        evi.UUID,
			Description: evi.Description,
			Operator:    dtos.User{FirstName: evi.FirstName, LastName: evi.LastName, Slug: evi.Slug},
			OccurredAt:  evi.OccurredAt,
			ContentType: evi.ContentType,
			Tags:        tags,
		}
	}
	return evidenceDTO, nil
}

func buildListEvidenceWhereClause(sb sq.SelectBuilder, operationID int64, filters helpers.TimelineFilters) sq.SelectBuilder {
	sb = sb.Where(sq.Eq{"evidence.operation_id": operationID})
	if len(filters.UUID) > 0 {
		sb = addWhereAndNot(sb, filters.UUID, evidenceUUIDWhere)
	}

	for _, text := range filters.Text {
		sb = sb.Where(sq.Like{"description": "%" + text + "%"})
	}

	if len(filters.DateRanges) > 0 {
		// sq.Or is a []sq.Sqlizer, so we can treat it like a slice (because it is one).
		stmts := make(sq.Or, len(filters.DateRanges))
		for i, v := range filters.DateRanges {
			stmts[i] = sq.And{
				sq.GtOrEq{"evidence.occurred_at": v.From},
				sq.LtOrEq{"evidence.occurred_at": v.To},
			}
		}
		sb = sb.Where(stmts)
	}

	if len(filters.Operator) > 0 {
		sb = addWhereAndNot(sb, filters.Operator, evidenceOperatorWhere)
	}

	if len(filters.Tags) > 0 {
		sb = sb.Where(eviForOpTagWhereComponent, filters.Tags, len(filters.Tags))
	}

	if len(filters.Type) > 0 {
		sb = addWhereAndNot(sb, filters.Type, evidenceTypeWhere)
	}

	if filters.Linked != nil {
		query := "evidence.id"
		if *filters.Linked {
			query += " IN "
		} else {
			query += " NOT IN "
		}
		query += eviLinkedSubquery
		sb = sb.Where(query)
	}

	return sb
}

const eviForOpTagWhereComponent = "evidence.id IN (" +
	"  SELECT evidence_id FROM tags" +
	"  LEFT JOIN tag_evidence_map ON tag_evidence_map.tag_id = tags.id" +
	"  WHERE tags.name IN (?)" +
	"  GROUP BY evidence_id HAVING COUNT(*) = ?" +
	")"
const eviLinkedSubquery = "(SELECT evidence_id FROM evidence_finding_map)"

func evidenceUUIDWhere(is bool) string {
	return "evidence.uuid " + isOrIsNot(is) + " (?)"
}

func evidenceOperatorWhere(is bool) string {
	return "evidence.operator_id " + isOrIsNot(is) + " (SELECT id FROM users WHERE slug IN (?))"
}

func evidenceTypeWhere(is bool) string {
	return "evidence.content_type " + isOrIsNot(is) + " (?)"
}

func addWhereAndNot(sb sq.SelectBuilder, vals filter.Values, whereFunc func(bool) string) sq.SelectBuilder {
	splitValues := vals.SplitByModifier()

	if values := splitValues[filter.Normal]; len(values) > 0 {
		sb = sb.Where(whereFunc(true), values)
	}
	if values := splitValues[filter.Not]; len(values) > 0 {
		sb = sb.Where(whereFunc(false), values)
	}
	return sb
}
