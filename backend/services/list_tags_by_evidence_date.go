// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type ListTagsByEvidenceDateInput struct {
	OperationSlug string
}

type TagUsageItem struct {
	TagID      int64     `db:"id"`
	OccurredAt time.Time `db:"occurred_at"`
}

type ExpandedTagUsageData struct {
	TagID      int64
	TagName    string
	ColorName  string
	UsageDates []time.Time
}

func ListTagsByEvidenceDate(ctx context.Context, db *database.Connection, i ListTagsByEvidenceDateInput) ([]*dtos.TagByEvidenceDate, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list tags by usage date", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list tags by usage date", backend.UnauthorizedReadErr(err))
	}

	var fullTagUsage []ExpandedTagUsageData

	// get tags by evidence occurred_at date
	err = db.WithTx(ctx, func(tx *database.Transactable) {
		var dbData []TagUsageItem
		tx.Select(&dbData,
			sq.Select("tags.id", "evidence.occurred_at").
				From("operations").
				LeftJoin("evidence ON operations.id = evidence.operation_id").
				LeftJoin("tag_evidence_map ON evidence.id = tag_evidence_map.evidence_id").
				LeftJoin("tags ON tags.id = tag_evidence_map.tag_id").
				Where(sq.Eq{"slug": i.OperationSlug}).
				Where(sq.NotEq{"tags.name": nil}).
				OrderBy("tags.id", "evidence.occurred_at")) // note sorting by tags.id is critical here

		fullTagUsage = foldTagUsageItems(dbData)
		tagIds := make([]int64, len(fullTagUsage))
		for i := range tagIds {
			tagIds[i] = fullTagUsage[i].TagID
		}

		var tags []models.Tag
		tx.Select(&tags,
			sq.Select("id", "name", "color_name").
				From("tags").
				Where(sq.Eq{"id": tagIds}).
				OrderBy("id")) // note sorting by id is critical here

		// since both queries sort by tags.id and the tags structure only has tag ids from the full set
		// we can do a simple merge here
		for i := range fullTagUsage {
			fullTagUsage[i].TagName = tags[i].Name
			fullTagUsage[i].ColorName = tags[i].ColorName
		}
	})

	if err != nil {
		return nil, backend.WrapError("Cannot get tags by usage date", backend.DatabaseErr(err))
	}

	tagDateDTO := make([]*dtos.TagByEvidenceDate, len(fullTagUsage))
	for idx, tag := range fullTagUsage {
		tagDateDTO[idx] = &dtos.TagByEvidenceDate{
			Tag: dtos.Tag{
				ID:        tag.TagID,
				Name:      tag.TagName,
				ColorName: tag.ColorName,
			},
			UsageDates: tag.UsageDates,
		}
	}
	return tagDateDTO, nil
}

func foldTagUsageItems(data []TagUsageItem) []ExpandedTagUsageData {
	tagData := []ExpandedTagUsageData{}

	currentTagID := int64(0)

	for _, tag := range data {
		if tag.TagID != currentTagID {
			tagData = append(tagData, ExpandedTagUsageData{TagID: tag.TagID, UsageDates: []time.Time{}})
			currentTagID = tag.TagID
		}
		lastItem := &tagData[len(tagData)-1]
		lastItem.UsageDates = append(lastItem.UsageDates, tag.OccurredAt)
	}

	return tagData
}
