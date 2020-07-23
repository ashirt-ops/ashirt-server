// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"
	"io"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type UpdateEvidenceInput struct {
	OperationSlug string
	EvidenceUUID  string
	Description   *string
	TagsToAdd     []int64
	TagsToRemove  []int64
	Content       io.Reader
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
		return backend.WrapError("Cannot updat eevidence", backend.DatabaseErr(err))
	}

	return nil
}
