// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"strings"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type ListTagsByEvidenceDateInput struct {
	OperationSlug string
}

func ListTagsByEvidenceDate(ctx context.Context, db *database.Connection, i ListTagsByEvidenceDateInput) ([]*dtos.TagByEvidenceDate, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	type TagDateData struct {
		TagID      int64  `db:"id"`
		TagName    string `db:"name"`
		ColorName  string `db:"color_name"`
		UsageDates string `db:"usage_dates"`
	}

	var dbData []TagDateData

	err = db.Select(&dbData,
		sq.Select("tags.id", "tags.name", "tags.color_name").
			Column("group_concat(DISTINCT evidence.occurred_at ORDER BY evidence.occurred_at ASC SEPARATOR ',') AS usage_dates").
			From("operations").
			LeftJoin("evidence ON operations.id = evidence.operation_id").
			LeftJoin("tag_evidence_map ON evidence.id = tag_evidence_map.evidence_id").
			LeftJoin("tags ON tags.id = tag_evidence_map.tag_id").
			Where(sq.Eq{"slug": i.OperationSlug}).
			GroupBy("tags.id"))

	if err != nil {
		return nil, backend.DatabaseErr(err)
	}

	tagDateDTO := make([]*dtos.TagByEvidenceDate, len(dbData))
	for idx, tag := range dbData {
		usage, err := sliceStrDatesToSliceDates(strings.Split(tag.UsageDates, ","))
		if err != nil {
			return nil, backend.DatabaseErr(err)
		}

		tagDateDTO[idx] = &dtos.TagByEvidenceDate{
			Tag: dtos.Tag{
				ID:        tag.TagID,
				Name:      tag.TagName,
				ColorName: tag.ColorName,
			},
			UsageDates: usage,
		}
	}
	return tagDateDTO, nil
}

func sliceStrDatesToSliceDates(dates []string) ([]time.Time, error) {
	times := make([]time.Time, len(dates))

	for i, strTime := range dates {
		t, err := time.Parse("2006-01-02 15:04:05", strTime) // parse dates as YYYY-MM-DD HH:mm:ss (with 0-prefixed times)
		if err != nil {
			return []time.Time{}, err
		}
		times[i] = t
	}

	return times, nil
}
