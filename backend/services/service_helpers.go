// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"

	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

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

	err := db.Get(&operation, sq.Select("id", "name", "status").
		From("operations").
		Where(sq.Eq{"slug": operationSlug}))
	return &operation, err
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
		return nil, nil, err
	}

	if finding.OperationID != operation.ID {
		return nil, nil, fmt.Errorf("Finding %d belongs to operation %d not %d", finding.ID, finding.OperationID, operation.ID)
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
		return nil, nil, err
	}

	if evidence.OperationID != operation.ID {
		return nil, nil, fmt.Errorf("Evidence %d belongs to operation %d not %d", evidence.ID, evidence.OperationID, operation.ID)
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
		return err
	}
	if len(badTags) > 0 {
		return fmt.Errorf("Tags [%v] do not belong to operation %s", badTags, operation.Slug)
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
	return userID, err
}

func selfOrSlugToUserID(ctx context.Context, db *database.Connection, slug string) (int64, error) {
	if slug == "" {
		return middleware.UserID(ctx), nil
	}
	var userID int64
	err := db.Get(&userID, sq.Select("id").From("users").Where(sq.Eq{"slug": slug}))
	return userID, err
}

func userIDFromSlugTx(tx *database.Transactable, slug string) int64 {
	var user models.User = models.User{ID: -1} // providing some default value in case tx.Get fails
	tx.Get(&user, sq.Select("id").From("users").Where(sq.Eq{"slug": slug}))
	return user.ID
}
