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
	Status        models.OperationStatus
}

type operationListItem struct {
	Op *dtos.Operation
	ID int64
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
			"name":   i.Name,
			"status": models.OperationStatusPlanning,
			"slug":   cleanSlug,
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
		Status:   models.OperationStatusPlanning,
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
	operations, err := listAllOperations(db)

	if err != nil {
		return nil, err
	}

	var operationPreference []models.UserOperationPermission

	err = db.Select(&operationPreference, sq.Select("operation_id", "is_favorite").
		From("user_operation_permissions").
		Where(sq.Eq{"user_id": middleware.UserID(ctx)}))

	if err != nil {
		return nil, backend.WrapError("Cannot get user operation permissions", backend.DatabaseErr(err))
	}

	operationPreferenceMap := make(map[int64]bool)
	for _, op := range operationPreference {
		operationPreferenceMap[op.OperationID] = op.IsFavorite
	}

	operationsDTO := make([]*dtos.Operation, 0, len(operations))
	for _, operation := range operations {
		if middleware.Policy(ctx).Check(policy.CanReadOperation{OperationID: operation.ID}) {
			fave, _ := operationPreferenceMap[operation.ID]
			operation.Op.Favorite = fave
			operationsDTO = append(operationsDTO, operation.Op)
		}
	}
	return operationsDTO, nil
}

func ReadOperation(ctx context.Context, db *database.Connection, operationSlug string) (*dtos.Operation, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to read operation", backend.UnauthorizedReadErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to read operation", backend.UnauthorizedReadErr(err))
	}

	var numUsers int
	var favorite bool

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Get(&numUsers, sq.Select("count(*)").From("user_operation_permissions").
			Where(sq.Eq{"operation_id": operation.ID}))

		tx.Get(&favorite, sq.Select("is_favorite").
			From("user_operation_permissions").
			Where(sq.Eq{"user_id": middleware.UserID(ctx), "operation_id": operation.ID}))
	})

	if err != nil {
		return nil, backend.WrapError("Cannot read favorite operation", backend.DatabaseErr(err))
	}

	return &dtos.Operation{
		Slug:     operationSlug,
		Name:     operation.Name,
		Status:   operation.Status,
		NumUsers: numUsers,
		Favorite: favorite,
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
			"name":   i.Name,
			"status": i.Status,
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
	ops, err := listAllOperations(db)

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
func listAllOperations(db *database.Connection) ([]operationListItem, error) {
	var operations []struct {
		models.Operation
		NumUsers int `db:"num_users"`
	}

	err := db.Select(&operations, sq.Select("id", "slug", "name", "status", "count(user_id) AS num_users").
		From("operations").
		LeftJoin("user_operation_permissions ON user_operation_permissions.operation_id = operations.id").
		GroupBy("operations.id").
		OrderBy("operations.created_at DESC"))
	if err != nil {
		return nil, backend.WrapError("Cannot list all operations", backend.DatabaseErr(err))
	}

	operationsDTO := []operationListItem{}
	for _, operation := range operations {
		operationsDTO = append(operationsDTO, operationListItem{
			ID: operation.ID,
			Op: &dtos.Operation{
				Slug:     operation.Slug,
				Name:     operation.Name,
				Status:   operation.Status,
				NumUsers: operation.NumUsers,
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
		return backend.WrapError("Unwilling to read operatoin", backend.UnauthorizedReadErr(err))
	}

	err = db.Update(sq.Update("user_operation_permissions").
		SetMap(map[string]interface{}{
			"is_favorite": i.IsFavorite,
		}).
		Where(sq.Eq{"operation_id": operation.ID, "user_id": middleware.UserID(ctx)}))
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
