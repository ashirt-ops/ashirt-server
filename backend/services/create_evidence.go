// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"
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

func CreateEvidence(ctx context.Context, db *database.Connection, contentStore contentstore.Store, i CreateEvidenceInput) (*dtos.Evidence, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.UnauthorizedWriteErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyEvidenceOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedWriteErr(err)
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
		case "terminal-recording":
			fallthrough
		case "codeblock":
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
			return nil, backend.UploadErr(err)
		}
	}

	evidenceUUID := uuid.New().String()

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		evidenceID, _ := tx.Insert("evidence", map[string]interface{}{
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
		return nil, backend.DatabaseErr(err)
	}

	return &dtos.Evidence{
		UUID:        evidenceUUID,
		Description: i.Description,
		OccurredAt:  i.OccurredAt,
	}, nil
}
