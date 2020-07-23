// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"io"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
)

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

func ReadEvidence(ctx context.Context, db *database.Connection, contentStore contentstore.Store, i ReadEvidenceInput) (*ReadEvidenceOutput, error) {
	operation, evidence, err := lookupOperationEvidence(db, i.OperationSlug, i.EvidenceUUID)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}
	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	var media io.Reader
	var preview io.Reader
	if i.LoadPreview {
		preview, err = contentStore.Read(evidence.ThumbImageKey)
		if err != nil {
			return nil, backend.WrapError("Unable to read evidence preview", err)
		}
	}

	if i.LoadMedia {
		media, err = contentStore.Read(evidence.FullImageKey)
		if err != nil {
			return nil, backend.WrapError("Unable to read evidence media", err)
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
