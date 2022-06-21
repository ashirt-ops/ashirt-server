// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/enhancementservices"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
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

// cloneContext makes a copy of a context that will persist beyond the given context's lifespan.
// Note: Most of the time, you probably don't want to use this. This is really only useful when
// you need to ensure that an operation persists longer than the request
func cloneContext(ctx context.Context) context.Context {
	// copy the logger
	rtnCtx, _ := logging.AddRequestLogger(context.Background(), logging.ReqLogger(ctx))

	return rtnCtx
}
