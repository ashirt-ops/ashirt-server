// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
	"golang.org/x/sync/errgroup"

	sq "github.com/Masterminds/squirrel"
)

type CreateOperationInput struct {
	Slug    string
	OwnerID int64
	Name    string
}

type UpdateOperationInput struct {
	OperationSlug string
	Name          string
}

type OperationWithID struct {
	Op *dtos.Operation
	ID int64
}

type TopContribWithID struct {
	dtos.TopContrib
	OperationID int64 `db:"operation_id" json:"operationId"`
}

type EvidenceCountWithID struct {
	dtos.EvidenceCount
	OperationID int64 `db:"operation_id" json:"operationId"`
}

func CreateOperation(ctx context.Context, db *database.Connection, i CreateOperationInput) (*dtos.Operation, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.CanCreateOperations{}); err != nil {
		return nil, backend.WrapError("Unable to create operation", backend.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, backend.MissingValueErr("Name")
	}

	if i.Slug == "" {
		return nil, backend.MissingValueErr("Slug")
	}

	cleanSlug := SanitizeOperationSlug(i.Slug)
	if cleanSlug == "" {
		return nil, backend.BadInputErr(errors.New("Unable to create operation. Invalid operation slug"), "Slug must contain english letters or numbers")
	}

	err := db.WithTx(ctx, func(tx *database.Transactable) {
		operationID, _ := tx.Insert("operations", map[string]interface{}{
			"name": i.Name,
			"slug": cleanSlug,
		})
		tx.Insert("user_operation_permissions", map[string]interface{}{
			"user_id":      i.OwnerID,
			"operation_id": operationID,
			"role":         policy.OperationRoleAdmin,
		})

		// Copy default tags into new operation
		tx.Exec(sq.Insert("tags").
			Columns(
				"name", "color_name",
				"operation_id",
			).
			Select(sq.Select(
				"name", "color_name",
				fmt.Sprintf("%v AS operation_id", operationID),
			).From("default_tags")),
		)
	})
	if err != nil {
		if database.IsAlreadyExistsError(err) {
			return nil, backend.WrapError("Unable to create operation. Operation slug already exists.", backend.BadInputErr(err, "An operation with this slug already exists"))
		}
		return nil, backend.WrapError("Unable to add new operation", backend.DatabaseErr(err))
	}

	return &dtos.Operation{
		Slug:     cleanSlug,
		Name:     i.Name,
		NumUsers: 1,
	}, nil
}

func DeleteOperation(ctx context.Context, db *database.Connection, contentStore contentstore.Store, slug string) error {
	operation, err := lookupOperation(db, slug)
	if err != nil {
		return backend.WrapError("Unable to delete operation", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanDeleteOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to delete operation", backend.UnauthorizedWriteErr(err))
	}
	log := logging.ReqLogger(ctx)

	var g errgroup.Group
	g.Go(func() error {
		err := db.WithTx(ctx, func(tx *database.Transactable) {
			var evidence []models.Evidence
			err = tx.Select(&evidence, sq.Select("*").From("evidence").Where(sq.Eq{"operation_id": operation.ID}))

			// remove evidence content
			if err == nil {
				for _, evi := range evidence {
					copy := evi
					g.Go(func() error {
						err := deleteEvidenceContent(contentStore, copy)
						if err != nil {
							log.Log("task", "delete operation", "msg", "error deleting evidence content", "uniqueKey", "orphanedDelete",
								"keys", fmt.Sprintf(`["%v", "%v"]`, copy.FullImageKey, copy.ThumbImageKey), "error", err.Error())
							return backend.DeleteErr(err)
						}
						return nil
					})
				}
			}

			// remove all tags for an operation
			var tagIDs []int64
			tx.Select(&tagIDs, sq.Select("id").From("tags").Where(sq.Eq{"operation_id": operation.ID}))
			tx.Delete(sq.Delete("tag_evidence_map").Where(sq.Eq{"tag_id": tagIDs}))
			tx.Delete(sq.Delete("tags").Where(sq.Eq{"id": tagIDs}))

			// remove all findings for an operation
			var findingIDs []int64
			tx.Select(&findingIDs, sq.Select("id").From("findings").Where(sq.Eq{"operation_id": operation.ID}))
			tx.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"finding_id": findingIDs}))
			tx.Delete(sq.Delete("findings").Where(sq.Eq{"id": findingIDs}))

			var evidenceIDs = make([]int64, len(evidence))
			for i, evi := range evidence {
				evidenceIDs[i] = evi.ID
			}

			// remove evidence metadata
			tx.Delete(sq.Delete("evidence_metadata").Where(sq.Eq{"evidence_id": evidenceIDs}))

			// remove all evidence
			tx.Delete(sq.Delete("evidence").Where(sq.Eq{"id": evidenceIDs}))

			// remove user/operations map
			tx.Delete(sq.Delete("user_operation_permissions").Where(sq.Eq{"operation_id": operation.ID}))
			// remove user preferences for operation
			tx.Delete(sq.Delete("user_operation_preferences").Where(sq.Eq{"operation_id": operation.ID}))

			tx.Delete(sq.Delete("operations").Where(sq.Eq{"id": operation.ID}))
		})
		if err != nil {
			log.Log("task", "delete operation", "msg", "Failed to fully delete operation data",
				"error", err.Error())
			return backend.WrapError("Cannot delete operation", backend.DatabaseErr(err))
		}
		return nil
	})

	return g.Wait()
}

// ListOperations retrieves a list of all operations that the contextual user can see
func ListOperations(ctx context.Context, db *database.Connection) ([]*dtos.Operation, error) {
	operations, err := listAllOperations(ctx, db)

	if err != nil {
		return nil, err
	}

	var operationPreference []models.UserOperationPreferences

	err = db.Select(&operationPreference, sq.Select("operation_id", "is_favorite").
		From("user_operation_preferences").
		Where(sq.Eq{"user_id": middleware.UserID(ctx)}))

	if err != nil {
		return nil, backend.WrapError("Cannot get user operation preferences", backend.DatabaseErr(err))
	}

	operationPreferenceMap := make(map[int64]bool)
	for _, op := range operationPreference {
		operationPreferenceMap[op.OperationID] = op.IsFavorite
	}

	operationsDTO := make([]*dtos.Operation, 0, len(operations))
	for _, operation := range operations {
		if middleware.Policy(ctx).Check(policy.CanReadOperation{OperationID: operation.ID}) {
			operation.Op.Favorite = operationPreferenceMap[operation.ID]
			operationsDTO = append(operationsDTO, operation.Op)
		}
	}
	return operationsDTO, nil
}

func ReadOperation(ctx context.Context, db *database.Connection, operationSlug string) (*dtos.Operation, error) {
	operation, err := lookupOperationWithCounts(db, operationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to read operation", backend.UnauthorizedReadErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to read operation", backend.UnauthorizedReadErr(err))
	}

	var numUsers int
	favorite := false
	var topContribs []TopContribWithID
	var evidenceCount []EvidenceCountWithID

	evidenceCountForOneOperation := fmt.Sprintf(`
	%s
	WHERE operation_id = ?
		GROUP BY operation_id`, getCountsFromEvidence)

	getTopContributorsForOperation := GetTopContributorsForEachOperation + ` AND t1.operation_id = ?`

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Get(&numUsers, sq.Select("count(*)").From("user_operation_permissions").
			Where(sq.Eq{"operation_id": operation.ID}))

		var favorites []bool
		tx.Select(&favorites, sq.Select("is_favorite").
			From("user_operation_preferences").
			Where(sq.Eq{"user_id": middleware.UserID(ctx), "operation_id": operation.ID}))
		if len(favorites) > 0 {
			favorite = favorites[0]
		}

		tx.SelectRawWithIntArg(&topContribs, getTopContributorsForOperation, operation.ID)

		tx.SelectRawWithIntArg(&evidenceCount, evidenceCountForOneOperation, operation.ID)
	})

	if err != nil {
		return nil, backend.WrapError("Cannot read operation", backend.DatabaseErr(err))
	}

	var evidenceCountForOp dtos.EvidenceCount
	if len(evidenceCount) > 0 {
		evidenceCountForOp.CodeblockCount = evidenceCount[0].CodeblockCount
		evidenceCountForOp.ImageCount = evidenceCount[0].ImageCount
		evidenceCountForOp.HarCount = evidenceCount[0].HarCount
		evidenceCountForOp.EventCount = evidenceCount[0].EventCount
		evidenceCountForOp.RecordingCount = evidenceCount[0].RecordingCount
	} else {
		evidenceCountForOp = dtos.EvidenceCount{}
	}

	var topContribsForOp []dtos.TopContrib
	if len(topContribs) > 0 {
		for i := range topContribs {
			var topContrib dtos.TopContrib
			topContrib.Slug = topContribs[i].Slug
			topContrib.Count = topContribs[i].Count
			topContribsForOp = append(topContribsForOp, topContrib)
		}
	} else {
		topContribsForOp = []dtos.TopContrib{}
	}

	return &dtos.Operation{
		Slug:          operationSlug,
		Name:          operation.Name,
		NumUsers:      numUsers,
		Favorite:      favorite,
		NumEvidence:   operation.NumEvidence,
		NumTags:       operation.NumTags,
		TopContribs:   topContribsForOp,
		EvidenceCount: evidenceCountForOp,
	}, nil
}

func UpdateOperation(ctx context.Context, db *database.Connection, i UpdateOperationInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to update operation", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to update operation", backend.UnauthorizedWriteErr(err))
	}

	err = db.Update(sq.Update("operations").
		SetMap(map[string]interface{}{
			"name": i.Name,
		}).
		Where(sq.Eq{"id": operation.ID}))
	if err != nil {
		return backend.WrapError("Cannot update operation", backend.DatabaseErr(err))
	}
	return nil
}

// ListOperationsForAdmin is a specialized version of ListOperations where no operations are filtered
// For use in admin screens only
func ListOperationsForAdmin(ctx context.Context, db *database.Connection) ([]*dtos.Operation, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}
	ops, err := listAllOperations(ctx, db)

	if err != nil {
		return nil, err
	}

	fixedOps := make([]*dtos.Operation, len(ops))

	for i, v := range ops {
		fixedOps[i] = v.Op
	}

	return fixedOps, nil
}

// listAllOperations is a helper function for both ListOperations and ListOpperationsForAdmin.
// This retrieves all operations, then relies on the caller to sort which operations are visible
// to the enduser
func listAllOperations(ctx context.Context, db *database.Connection) ([]OperationWithID, error) {
	var operations []struct {
		models.Operation
		NumUsers    int `db:"num_users"`
		NumEvidence int `db:"num_evidence"`
		NumTags     int `db:"num_tags"`
	}

	var topContribs []TopContribWithID

	var evidenceCount []EvidenceCountWithID

	err := db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Select(&operations, sq.Select("operations.id", "slug", "operations.name", "count(distinct(user_operation_permissions.user_id)) AS num_users", "count(distinct(evidence.id)) AS num_evidence", "count(distinct(tags.id)) AS num_tags").
			From("operations").
			LeftJoin("user_operation_permissions ON user_operation_permissions.operation_id = operations.id").
			LeftJoin("evidence ON evidence.operation_id = operations.id").
			LeftJoin("tags ON tags.operation_id = operations.id").
			GroupBy("operations.id").
			OrderBy("operations.created_at DESC"))

		tx.SelectRaw(&topContribs, GetTopContributorsForEachOperation)

		tx.SelectRaw(&evidenceCount, EvidenceCountForAllOperations)
	})

	if err != nil {
		return nil, backend.WrapError("Cannot list all operations", backend.DatabaseErr(err))
	}

	operationsDTO := []OperationWithID{}
	for _, operation := range operations {

		filteredTopContribs := helpers.Filter(topContribs, func(contributor TopContribWithID) bool {
			return contributor.OperationID == operation.ID
		})

		topContribsForOp := make([]dtos.TopContrib, 0, len(filteredTopContribs))
		for i := range filteredTopContribs {
			var topContrib dtos.TopContrib
			topContrib.Slug = filteredTopContribs[i].Slug
			topContrib.Count = filteredTopContribs[i].Count
			topContribsForOp = append(topContribsForOp, topContrib)
		}

		var evidenceCountForOp dtos.EvidenceCount
		idx, _ := helpers.Find(evidenceCount, func(item EvidenceCountWithID) bool {
			return item.OperationID == operation.ID
		})
		if idx > -1 {
			evidenceCountForOp.CodeblockCount = evidenceCount[idx].CodeblockCount
			evidenceCountForOp.ImageCount = evidenceCount[idx].ImageCount
			evidenceCountForOp.HarCount = evidenceCount[idx].HarCount
			evidenceCountForOp.EventCount = evidenceCount[idx].EventCount
			evidenceCountForOp.RecordingCount = evidenceCount[idx].RecordingCount
		}
		operationsDTO = append(operationsDTO, OperationWithID{
			ID: operation.ID,
			Op: &dtos.Operation{
				Slug:          operation.Slug,
				Name:          operation.Name,
				NumUsers:      operation.NumUsers,
				NumEvidence:   operation.NumEvidence,
				NumTags:       operation.NumTags,
				TopContribs:   topContribsForOp,
				EvidenceCount: evidenceCountForOp,
			},
		})
	}
	return operationsDTO, nil
}

type SetFavoriteInput struct {
	OperationSlug string
	IsFavorite    bool
}

func SetFavoriteOperation(ctx context.Context, db *database.Connection, i SetFavoriteInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to read operation", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to read operation", backend.UnauthorizedReadErr(err))
	}

	_, err = db.Insert("user_operation_preferences", map[string]interface{}{
		"user_id":      middleware.UserID(ctx),
		"operation_id": operation.ID,
		"is_favorite":  i.IsFavorite,
	}, "ON DUPLICATE KEY UPDATE is_favorite=VALUES(is_favorite)")

	if err != nil {
		return backend.WrapError("Cannot set operation as favorite", backend.DatabaseErr(err))
	}

	return nil
}

var disallowedCharactersRegex = regexp.MustCompile(`[^A-Za-z0-9]+`)

// SanitizeOperationSlug removes objectionable characters from a slug and returns the new slug.
// Current logic: only allow alphanumeric characters and hyphen, with hypen excluded at the start
// and end
func SanitizeOperationSlug(slug string) string {
	return strings.Trim(
		disallowedCharactersRegex.ReplaceAllString(strings.ToLower(slug), "-"),
		"-",
	)
}

var getDataFromEvidence string = `
	SELECT
		slug,
		operation_id,
		count(evidence.id) AS count
	FROM
		evidence
		LEFT JOIN users ON evidence.operator_id = users.id`

var GetTopContributorsForEachOperation string = fmt.Sprintf(`
	SELECT
		t1.*
	FROM (%s
	GROUP BY
		operation_id,
		users.id) t1
		LEFT JOIN (
			%s	
			GROUP BY
				operation_id,
				users.id) t2 ON t1.operation_id = t2.operation_id
		AND t1.count < t2.count
	WHERE
		t2.count IS NULL`, getDataFromEvidence, getDataFromEvidence)

var getCountsFromEvidence string = `
SELECT operation_id,
	COUNT(CASE WHEN content_type = "image" THEN 1 END) image_count,
	COUNT(CASE WHEN content_type = "codeblock" THEN 1 END) codeblock_count,
	COUNT(CASE WHEN content_type = "terminal-recording" THEN 1 END) recording_count,
	COUNT(CASE WHEN content_type = "event" THEN 1 END) event_count,
	COUNT(CASE WHEN content_type = "http-request-cycle" THEN 1 END) har_count
FROM
	evidence
`

var EvidenceCountForAllOperations string = fmt.Sprintf(`
	%s
	GROUP BY 
		operation_id`, getCountsFromEvidence)
