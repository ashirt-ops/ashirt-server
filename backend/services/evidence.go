// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/enhancementservices"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/helpers/filter"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"
	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"
)

type CreateEvidenceInput struct {
	OperatorID    int64
	OperationSlug string
	Description   string
	Content       io.Reader
	ContentType   string
	TagIDs        []int64
	OccurredAt    time.Time
}

type DeleteEvidenceInput struct {
	OperationSlug            string
	EvidenceUUID             string
	DeleteAssociatedFindings bool
}

type ListEvidenceForFindingInput struct {
	OperationSlug string
	FindingUUID   string
}

type ListEvidenceForOperationInput struct {
	OperationSlug string
	Filters       helpers.TimelineFilters
}

type ReadEvidenceInput struct {
	OperationSlug string
	EvidenceUUID  string
	LoadPreview   bool
	LoadMedia     bool
}

type ReadEvidenceOutput struct {
	UUID        string    `json:"uuid"`
	Description string    `json:"description"`
	ContentType string    `json:"contentType"`
	OccurredAt  time.Time `json:"occurredAt"`
	Preview     io.Reader `json:"-"`
	Media       io.Reader `json:"-"`
}

type UpdateEvidenceInput struct {
	OperationSlug string
	EvidenceUUID  string
	Description   *string
	TagsToAdd     []int64
	TagsToRemove  []int64
	Content       io.Reader
}

type MoveEvidenceInput struct {
	SourceOperationSlug string
	EvidenceUUID        string
	TargetOperationSlug string
}

func CreateEvidence(ctx context.Context, db *database.Connection, contentStore contentstore.Store, i CreateEvidenceInput) (*dtos.Evidence, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to create evidence", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unable to create evidence", backend.UnauthorizedWriteErr(err))
	}

	if i.OccurredAt.IsZero() {
		i.OccurredAt = time.Now()
	}

	if err := ensureTagIDsBelongToOperation(db, i.TagIDs, operation); err != nil {
		return nil, backend.BadInputErr(err, err.Error())
	}

	keys := contentstore.ContentKeys{}

	if i.Content != nil {
		var content contentstore.Storable
		switch i.ContentType {
		case "http-request-cycle":
			fallthrough
		case "terminal-recording":
			fallthrough
		case "codeblock":
			fallthrough
		case "event":
			content = contentstore.NewBlob(i.Content)

		case "image":
			fallthrough
		default:
			content = contentstore.NewImage(i.Content)
		}

		keys, err = content.ProcessPreviewAndUpload(contentStore)
		if err != nil {
			if httpErr, ok := err.(*backend.HTTPError); ok {
				return nil, httpErr
			}
			return nil, backend.WrapError("Unable to upload evidence", backend.UploadErr(err))
		}
	}

	evidenceUUID := uuid.New().String()
	var evidenceID int64
	err = db.WithTx(ctx, func(tx *database.Transactable) {
		evidenceID, _ = tx.Insert("evidence", map[string]interface{}{
			"uuid":            evidenceUUID,
			"description":     i.Description,
			"content_type":    i.ContentType,
			"occurred_at":     i.OccurredAt,
			"operation_id":    operation.ID,
			"operator_id":     middleware.UserID(ctx),
			"full_image_key":  keys.Full,
			"thumb_image_key": keys.Thumbnail,
		})
		tx.BatchInsert("tag_evidence_map", len(i.TagIDs), func(idx int) map[string]interface{} {
			return map[string]interface{}{
				"tag_id":      i.TagIDs[idx],
				"evidence_id": evidenceID,
			}
		})
	})

	if err != nil {
		return nil, backend.WrapError("Could not create evidence and tags", backend.DatabaseErr(err))
	}

	err = enhancementservices.SendEvidenceCreatedEvent(db, logging.ReqLogger(ctx), operation.ID, []string{evidenceUUID}, enhancementservices.AllWorkers())

	if err != nil {
		logging.Log(ctx, "msg", "Unable to run workers", "error", err.Error())
	}

	return &dtos.Evidence{
		UUID:        evidenceUUID,
		Description: i.Description,
		OccurredAt:  i.OccurredAt,
	}, nil
}

func DeleteEvidence(ctx context.Context, db *database.Connection, contentStore contentstore.Store, i DeleteEvidenceInput) error {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return backend.WrapError("Unable to delete evidence", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to delete evidence", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		if i.DeleteAssociatedFindings {
			tx.Exec(sq.Expr("DELETE findings FROM findings INNER JOIN evidence_finding_map ON findings.id = evidence_finding_map.finding_id WHERE evidence_id = ?", evidence.ID))
		}
		tx.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"evidence_id": evidence.ID}))
		tx.Delete(sq.Delete("evidence_metadata").Where(sq.Eq{"evidence_id": evidence.ID}))
		tx.Delete(sq.Delete("evidence").Where(sq.Eq{"id": evidence.ID}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete evidence", backend.DatabaseErr(err))
	}

	if err = deleteEvidenceContent(contentStore, *evidence); err != nil {
		return backend.WrapError("Cannot delete evidence content", backend.DeleteErr(err))
	}

	return nil
}

func ListEvidenceForFinding(ctx context.Context, db *database.Connection, i ListEvidenceForFindingInput) ([]dtos.Evidence, error) {
	operation, finding, err := lookupOperationFinding(db, i.OperationSlug, i.FindingUUID)
	if err != nil {
		return nil, backend.WrapError("Unable to list evidence for finding", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list evidence for finding", backend.UnauthorizedReadErr(err))
	}

	var evidenceForFinding []struct {
		models.Evidence
		Slug      string `db:"slug"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	err = db.Select(&evidenceForFinding, sq.Select("evidence.*", "slug", "first_name", "last_name").
		From("evidence").
		LeftJoin("evidence_finding_map ON evidence.id = evidence_id").
		LeftJoin("users ON users.id = evidence.operator_id").
		Where(sq.Eq{"finding_id": finding.ID}))

	if err != nil {
		return nil, backend.WrapError("Cannot list evidence for finding", backend.UnauthorizedReadErr(err))
	}

	evidenceIDs := make([]int64, len(evidenceForFinding))
	for idx, evi := range evidenceForFinding {
		evidenceIDs[idx] = evi.Evidence.ID
	}

	tagsByEvidenceID, _, err := tagsForEvidenceByID(db, evidenceIDs)
	if err != nil {
		return nil, backend.WrapError("Cannot get tags for evidnece", backend.UnauthorizedReadErr(err))
	}

	var evidenceDTOs = make([]dtos.Evidence, len(evidenceForFinding))
	for i, evi := range evidenceForFinding {
		tags := tagsByEvidenceID[evi.Evidence.ID]
		if tags == nil {
			tags = []dtos.Tag{}
		}
		evidenceDTOs[i] = dtos.Evidence{
			UUID:        evi.UUID,
			ContentType: evi.ContentType,
			Description: evi.Description,
			OccurredAt:  evi.OccurredAt,
			Tags:        tags,
			Operator: dtos.User{
				Slug:      evi.Slug,
				FirstName: evi.FirstName,
				LastName:  evi.LastName,
			},
		}
	}

	return evidenceDTOs, nil
}

// ListEvidenceForOperation retrieves all evidence for a particular operation id matching a particular
// set of filters (e.g. tag:some_tag)
func ListEvidenceForOperation(ctx context.Context, db *database.Connection, contentStore contentstore.Store, i ListEvidenceForOperationInput) ([]*dtos.Evidence, error) {
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

	sb := sq.Select().
		From("evidence").
		LeftJoin("users ON evidence.operator_id = users.id").
		Columns(
			"evidence.id",
			"evidence.uuid",
			"description",
			"evidence.content_type",
			"occurred_at",
			"users.first_name",
			"users.last_name",
			"users.slug",
		)

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

	sendUrl := false
	fmt.Println("about to check contentSTore")
	if _, ok := contentStore.(*contentstore.S3Store); ok {
		fmt.Println("S3Store")
		sendUrl = true
	}

	fmt.Println("sendUrl", sendUrl)

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
			SendUrl:     sendUrl,
		}
	}
	return evidenceDTO, nil
}

func SendUrl(ctx context.Context, db *database.Connection, contentStore *contentstore.S3Store, i ReadEvidenceInput) (*string, error) {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return nil, backend.WrapError("Unable to read evidence", backend.UnauthorizedReadErr(err))
	}
	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to read evidence", backend.UnauthorizedReadErr(err))
	}
	str, err := contentStore.SendURL(evidence.FullImageKey)
	if err != nil {
		return nil, backend.WrapError("Unable to get image URL", backend.ServerErr(err))
	}

	return str, nil

}
func ReadEvidence(ctx context.Context, db *database.Connection, contentStore contentstore.Store, i ReadEvidenceInput) (*ReadEvidenceOutput, error) {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return nil, backend.WrapError("Unable to read evidence", backend.UnauthorizedReadErr(err))
	}
	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to read evidence", backend.UnauthorizedReadErr(err))
	}

	var media io.Reader
	var preview io.Reader
	if i.LoadPreview {
		preview, err = contentStore.Read(evidence.ThumbImageKey)
		if err != nil {
			return nil, backend.WrapError("Cannot read evidence preview", err)
		}
	}

	if i.LoadMedia {
		media, err = contentStore.Read(evidence.FullImageKey)
		if err != nil {
			return nil, backend.WrapError("Cannot read evidence media", err)
		}
	}

	return &ReadEvidenceOutput{
		UUID:        evidence.UUID,
		Description: evidence.Description,
		ContentType: evidence.ContentType,
		OccurredAt:  evidence.OccurredAt,
		Media:       media,
		Preview:     preview,
	}, nil
}

func UpdateEvidence(ctx context.Context, db *database.Connection, contentStore contentstore.Store, i UpdateEvidenceInput) error {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return backend.WrapError("Unable to update evidence", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to update evidence", backend.UnauthorizedWriteErr(err))
	}

	if err := ensureTagIDsBelongToOperation(db, i.TagsToAdd, operation); err != nil {
		return backend.WrapError("Unable to update evidence", backend.BadInputErr(err, err.Error()))
	}

	var keys *contentstore.ContentKeys
	if i.Content != nil {
		switch evidence.ContentType {
		case "http-request-cycle":
			fallthrough
		case "codeblock":
			fallthrough
		case "terminal-recording":
			content := contentstore.NewBlob(i.Content)
			processedKeys, err := content.ProcessPreviewAndUpload(contentStore)
			if err != nil {
				return backend.WrapError("Cannot update evidence content", backend.BadInputErr(err, "Failed to process content"))
			}
			keys = &processedKeys

		case "image":
			fallthrough
		default:
			err := errors.New("Content cannot be updated")
			return backend.BadInputErr(err, err.Error())
		}
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		ub := sq.Update("evidence").Where(sq.Eq{"id": evidence.ID})
		if i.Description != nil {
			ub = ub.Set("description", i.Description)
		}
		if keys != nil {
			ub = ub.SetMap(map[string]interface{}{
				"full_image_key":  keys.Full,
				"thumb_image_key": keys.Thumbnail,
			})
		}

		if _, _, err := ub.ToSql(); err == nil {
			tx.Update(ub)
		}

		tx.Delete(sq.Delete("tag_evidence_map").Where(sq.Eq{"evidence_id": evidence.ID, "tag_id": i.TagsToRemove}))

		if len(i.TagsToAdd) > 0 {
			tx.BatchInsert("tag_evidence_map", len(i.TagsToAdd), func(idx int) map[string]interface{} {
				return map[string]interface{}{
					"tag_id":      i.TagsToAdd[idx],
					"evidence_id": evidence.ID,
				}
			})
		}
	})
	if err != nil {
		return backend.WrapError("Cannot update evidence", backend.DatabaseErr(err))
	}

	return nil
}

func buildListEvidenceWhereClause(sb sq.SelectBuilder, operationID int64, filters helpers.TimelineFilters) sq.SelectBuilder {
	sb = sb.Where(sq.Eq{"evidence.operation_id": operationID})
	if len(filters.UUID) > 0 {
		sb = addWhereAndNot(sb, filters.UUID, evidenceUUIDWhere)
	}

	for _, text := range filters.Text {
		sb = sb.Where(sq.Like{"description": "%" + text + "%"})
	}

	if len(filters.Metadata) > 0 {
		metadataSubquery := sq.Select("evidence_id").From("evidence_metadata")
		for _, text := range filters.Metadata {
			metadataSubquery = metadataSubquery.Where(sq.Like{"body": "%" + text + "%"})
		}
		if q, v, e := metadataSubquery.ToSql(); e == nil {
			sb = sb.Where("evidence.id IN ("+q+")", v...)
		}
	}

	if values := filters.DateRanges; len(values) > 0 {
		splitValues := values.SplitByModifier()

		if splitVals := splitValues[filter.Normal]; len(splitVals) > 0 {
			stmts := make(sq.Or, len(splitVals))
			for i, v := range splitVals {
				stmts[i] = sq.And{
					sq.GtOrEq{"evidence.occurred_at": v.From},
					sq.LtOrEq{"evidence.occurred_at": v.To},
				}
			}
			sb = sb.Where(stmts)
		}
		if splitVals := splitValues[filter.Not]; len(splitVals) > 0 {
			// there's not a great way to do this, so falling back to expr and string construction
			stmts := make(sq.And, len(splitVals))
			for i, v := range splitVals {
				stmts[i] = sq.Expr(
					"NOT( evidence.occurred_at >= ? AND evidence.occurred_at <= ?)", v.From, v.To,
				)
			}
			sb = sb.Where(stmts)
		}
	}

	if len(filters.Operator) > 0 {
		sb = addWhereAndNot(sb, filters.Operator, evidenceOperatorWhere)
	}

	if len(filters.Tags) > 0 {
		sb = addWhereAndNot(sb, filters.Tags, evidenceTagOrWhere)
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

const eviLinkedSubquery = "(SELECT evidence_id FROM evidence_finding_map)"

func evidenceUUIDWhere(in bool) string {
	return "evidence.uuid " + inOrNotIn(in) + " (?)"
}

func evidenceOperatorWhere(in bool) string {
	return "evidence.operator_id " + inOrNotIn(in) + " (SELECT id FROM users WHERE slug IN (?))"
}

func evidenceTypeWhere(in bool) string {
	return "evidence.content_type " + inOrNotIn(in) + " (?)"
}

func evidenceTagOrWhere(in bool) string {
	return evidenceTagWhere(in, false)
}

// func evidenceTagAndWhere(is bool) string {
// 	return evidenceTagWhere(is, true)
// }

func evidenceTagWhere(in, all bool) string {
	groupBy := ""
	if all {
		groupBy = "  GROUP BY evidence_id HAVING COUNT(*) = ?"
	}
	return "evidence.id " + inOrNotIn(in) + " (" +
		"  SELECT evidence_id FROM tag_evidence_map" +
		"  LEFT JOIN tags ON tag_evidence_map.tag_id = tags.id" +
		"  WHERE tags.name IN (?)" +
		groupBy +
		")"
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

func deleteEvidenceContent(contentStore contentstore.Store, evidence models.Evidence) error {
	keys := make([]string, 0, 2)
	if evidence.FullImageKey != "" {
		keys = append(keys, evidence.FullImageKey)
	}
	if evidence.ThumbImageKey != "" && evidence.ThumbImageKey != evidence.FullImageKey {
		keys = append(keys, evidence.ThumbImageKey)
	}
	for _, key := range keys {
		err := contentStore.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func MoveEvidence(ctx context.Context, db *database.Connection, i MoveEvidenceInput) error {
	sourceOperation, evidence, err := lookupOperationEvidence(db, i.SourceOperationSlug, i.EvidenceUUID)
	if err != nil {
		return backend.WrapError("Unable to move evidence (src)", backend.UnauthorizedReadErr(err))
	}

	destinationOperation, err := lookupOperation(db, i.TargetOperationSlug)
	if err != nil {
		return backend.WrapError("Unable to move evidence (dst op)", backend.UnauthorizedReadErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx,
		policy.CanModifyOperation{OperationID: sourceOperation.ID},
		policy.CanModifyOperation{OperationID: destinationOperation.ID},
	); err != nil {
		return backend.WrapError("Unwilling to move evidence", backend.UnauthorizedWriteErr(err))
	}

	//Check which tags can be migrated
	tagDifferences, err := ListTagDifferenceForEvidence(ctx, db, ListTagDifferenceForEvidenceInput{
		ListTagsDifferenceInput: ListTagsDifferenceInput{
			SourceOperationSlug:      i.SourceOperationSlug,
			DestinationOperationSlug: i.TargetOperationSlug,
		},
		SourceEvidenceUUID: i.EvidenceUUID,
	})

	if err != nil {
		return backend.WrapError("Unable to list tag differences for moving", err)
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		// remove findings
		tx.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"evidence_id": evidence.ID}))
		// remove tags
		tx.Delete(sq.Delete("tag_evidence_map").Where(sq.Eq{"evidence_id": evidence.ID}))
		// reassociate evidence with new operation
		tx.Update(sq.Update("evidence").Set("operation_id", destinationOperation.ID).Where(sq.Eq{"id": evidence.ID}))
		// associate with common tags
		tx.BatchInsert("tag_evidence_map", len(tagDifferences.Included), func(idx int) map[string]interface{} {
			pair := tagDifferences.Included[idx]
			return map[string]interface{}{
				"tag_id":      pair.DestinationTag.ID,
				"evidence_id": evidence.ID,
			}
		})
	})
	if err != nil {
		return backend.WrapError("Cannot move evidence", err)
	}

	return nil
}
