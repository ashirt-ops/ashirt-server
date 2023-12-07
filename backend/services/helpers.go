// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

// getFindingCategory returns the category associated with the provided category id
// if this record is not found, then an empty string will be returned. If an error occurs,
// then an error will be returned.
func getFindingCategory(db *database.Connection, findingCategoryID int64) (string, error) {
	var category string
	err := db.Get(&category, sq.Select("category").
		From("finding_categories").
		Where(sq.Eq{"id": findingCategoryID}),
	)

	return category, err
}

// getFindingCategoryID retrieves the ID associated with the provided category.
// this function accepts a select function, which is intended to be a (*database.Connection).Select,
// or a (*database.Transactable).Select
func getFindingCategoryID(findingCategory string, selectFunc func(modalSlice interface{}, sb sq.SelectBuilder) error) (*int64, error) {
	var foundCategoryID []int64

	// look up the category -- it might be null. We don't want to create it here.
	err := selectFunc(&foundCategoryID, sq.Select("id").
		From("finding_categories").
		Where(sq.Eq{"category": findingCategory}),
	)
	if err != nil {
		return nil, err
	}
	if len(foundCategoryID) == 0 {
		return nil, nil
	}
	return &foundCategoryID[0], nil

}

// tagsForEvidenceByID retrieves a list of Tag structures for the specified evidence ids
func tagsForEvidenceByID(db *database.Connection, evidenceIDs []int64) (tagsByEvidenceID map[int64][]dtos.Tag, allTags []dtos.Tag, err error) {
	if len(evidenceIDs) == 0 {
		allTags = []dtos.Tag{}
		return
	}
	var tags []struct {
		models.Tag
		EvidenceID int64 `db:"evidence_id"`
	}

	err = db.Select(&tags, sq.Select("evidence_id", "tags.*").
		From("tag_evidence_map").
		LeftJoin("tags ON tag_id = tags.id").
		Where(sq.Eq{"evidence_id": evidenceIDs}).
		OrderBy("tag_id ASC"))
	if err != nil {
		return
	}

	allTagsByTagID := map[int64]bool{}
	allTags = []dtos.Tag{}
	tagsByEvidenceID = map[int64][]dtos.Tag{}
	for _, tag := range tags {
		tagDTO := dtos.Tag{
			ID:        tag.ID,
			Name:      tag.Name,
			ColorName: tag.ColorName,
		}
		if !allTagsByTagID[tag.ID] {
			allTags = append(allTags, tagDTO)
			allTagsByTagID[tag.ID] = true
		}
		tagsByEvidenceID[tag.EvidenceID] = append(tagsByEvidenceID[tag.EvidenceID], tagDTO)
	}

	return
}

// lookupOperation returns an operation model for the given slug
func lookupOperation(db *database.Connection, operationSlug string) (*models.Operation, error) {
	var operation models.Operation

	err := db.Get(&operation, sq.Select("id", "name").
		From("operations").
		Where(sq.Eq{"slug": operationSlug}))
	if err != nil {
		return &operation, backend.WrapError("Unable to lookup operation by slug", err)
	}
	return &operation, nil
}

type operationWithCounts struct {
	models.Operation
	NumEvidence int `db:"num_evidence"`
	NumTags     int `db:"num_tags"`
}

// lookupOperation returns an operation model for the given slug
func lookupOperationWithCounts(db *database.Connection, operationSlug string) (*operationWithCounts, error) {
	var opAndData operationWithCounts
	err := db.Get(&opAndData, sq.Select("operations.id", "operations.name", "count(distinct(tags.id)) AS num_tags", "count(distinct(evidence.id)) AS num_evidence").
		LeftJoin("evidence ON evidence.operation_id = operations.id").
		LeftJoin("tags ON tags.operation_id = operations.id").
		From("operations").
		GroupBy("operations.id").
		Where(sq.Eq{"slug": operationSlug}))
	if err != nil {
		return &opAndData, backend.WrapError("Unable to lookup operation by slug", err)
	}
	return &opAndData, nil
}

// lookupOperationFinding returns an operation & finding model for the given operation slug / finding uuid
// and ensures that the finding belongs to the specified operation
func lookupOperationFinding(db *database.Connection, operationSlug string, findingUUID string) (*models.Operation, *models.Finding, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, nil, err
	}

	var finding models.Finding
	err = db.Get(&finding, sq.Select("*").From("findings").Where(sq.Eq{"uuid": findingUUID}))
	if err != nil {
		return nil, nil, backend.WrapError("Unable to lookup finding by uuid", err)
	}

	if finding.OperationID != operation.ID {
		return nil, nil, fmt.Errorf("Unable to lookup operation/finding. Finding %d belongs to operation %d not %d", finding.ID, finding.OperationID, operation.ID)
	}

	return operation, &finding, nil
}

// lookupOperationEvidence returns an operation & evidence model for the given operation slug / evidence uuid
// and ensures that the evidence belongs to the specified operation
func lookupOperationEvidence(db *database.Connection, operationSlug string, evidenceUUID string) (*models.Operation, *models.Evidence, error) {
	operation, err := lookupOperation(db, operationSlug)
	if err != nil {
		return nil, nil, err
	}

	var evidence models.Evidence
	err = db.Get(&evidence, sq.Select("*").
		From("evidence").
		Where(sq.Eq{"uuid": evidenceUUID}))
	if err != nil {
		return nil, nil, backend.WrapError("Unable to lookup evidence by uuid", err)
	}

	if evidence.OperationID != operation.ID {
		return nil, nil, fmt.Errorf("Unable to lookup operation/evidence. Evidence %d belongs to operation %d not %d", evidence.ID, evidence.OperationID, operation.ID)
	}

	return operation, &evidence, nil
}

// Returns an error if any specified tagIDs belong to an operation other than the one specified
func ensureTagIDsBelongToOperation(db *database.Connection, tagIDs []int64, operation *models.Operation) error {
	if len(tagIDs) == 0 {
		return nil
	}
	var badTags []dtos.Tag

	err := db.Select(&badTags, sq.Select("id").
		From("tags").
		Where(sq.Eq{"id": tagIDs}).
		Where(sq.NotEq{"operation_id": operation.ID}))

	if err != nil {
		return backend.WrapError("Unable to lookup tags by operation ID", err)
	}
	if len(badTags) > 0 {
		return fmt.Errorf("Unable to verify tags for operation. Tags [%v] do not belong to operation %s", badTags, operation.Slug)
	}
	return nil
}

// policyRequireWithAdminBypass is a small wrapper around policy.Require. In addition to normal policy checks,
// this will also check if the user is an admin. If so, then the admin is permitted to act
//
// Note: this is not always desirable to use, as it will show Admin users non-personalized content (i.e. no filtering)
func policyRequireWithAdminBypass(ctx context.Context, requiredPermissions ...policy.Permission) error {
	if middleware.IsAdmin(ctx) {
		return nil
	}

	return policy.Require(middleware.Policy(ctx), requiredPermissions...)
}

// isAdmin checks if the admin flag is set in the context. If not, then a standard error is returned
func isAdmin(ctx context.Context) error {
	if !middleware.IsAdmin(ctx) {
		return fmt.Errorf("Requesting user is not an admin")
	}
	return nil
}

func userSlugToUserID(db *database.Connection, slug string) (int64, error) {
	var userID int64
	err := db.Get(&userID, sq.Select("id").From("users").Where(sq.Eq{"slug": slug}))
	if err != nil {
		return userID, backend.WrapError("Unable to look up user by slug", err)
	}
	return userID, err
}

func userGroupSlugToUserGroupID(db *database.Connection, slug string) (int64, error) {
	var userGroupID int64
	err := db.Get(&userGroupID, sq.Select("id").From("user_groups").Where(sq.Eq{"slug": slug}))
	if err != nil {
		return userGroupID, backend.WrapError("Unable to look up user group by slug", err)
	}
	return userGroupID, err
}

// lookupUserGroup returns an user group model for the given slug
func lookupUserGroup(db *database.Connection, userGroupSlug string) (*models.UserGroup, error) {
	var userGroup models.UserGroup

	err := db.Get(&userGroup, sq.Select("id", "name", "slug").
		From("user_groups").
		Where(sq.Eq{"slug": userGroupSlug}))
	if err != nil {
		return &userGroup, backend.WrapError("Unable to lookup user group by slug", err)
	}
	return &userGroup, nil
}

func SelfOrSlugToUserID(ctx context.Context, db *database.Connection, slug string) (int64, error) {
	if slug == "" {
		return middleware.UserID(ctx), nil
	}
	return userSlugToUserID(db, slug)
}

func ListActiveServices(ctx context.Context, db *database.Connection) ([]*dtos.ActiveServiceWorker, error) {
	var serviceWorkers []models.ServiceWorker
	err := db.Select(&serviceWorkers, sq.Select("name").
		From("service_workers").
		Where(sq.Eq{"deleted_at": nil}))
	if err != nil {
		return nil, err
	}

	servicesDTO := helpers.Map(serviceWorkers, func(t models.ServiceWorker) *dtos.ActiveServiceWorker {
		return &dtos.ActiveServiceWorker{
			Name: t.Name,
		}
	})
	return servicesDTO, nil
}

var disallowedCharactersRegex = regexp.MustCompile(`[^A-Za-z0-9]+`)

// SanitizeOperationSlug removes objectionable characters from a slug and returns the new slug.
// Current logic: only allow alphanumeric characters and hyphen, with hypen excluded at the start
// and end
func SanitizeSlug(slug string) string {
	return strings.Trim(
		disallowedCharactersRegex.ReplaceAllString(strings.ToLower(slug), "-"),
		"-",
	)
}

func LookupGlobalVar(db *database.Connection, name string) (*models.GlobalVar, error) {
	var globalVar models.GlobalVar

	err := db.Get(&globalVar, sq.Select("*").
		From("global_vars").
		Where(sq.Eq{"name": name}))
	if err != nil {
		return &globalVar, backend.WrapError("Unable to lookup global variable by name", err)
	}
	return &globalVar, nil
}

func LookupOperationVar(db *database.Connection, varSlug string) (*models.OperationVar, error) {
	var operationVar models.OperationVar

	err := db.Get(&operationVar, sq.Select("*").
		From("operation_vars").
		Where(sq.Eq{"slug": varSlug}))
	if err != nil {
		return &operationVar, backend.WrapError("Unable to lookup operation variable by name", err)
	}
	return &operationVar, nil
}
