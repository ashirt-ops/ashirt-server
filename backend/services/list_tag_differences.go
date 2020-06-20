// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"strings"

	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
)

type ListTagsDifferenceInput struct {
	SourceOperationSlug      string
	DestinationOperationSlug string
}

// ListTagDifference determines which tag values are common between two operations. This is done via
// relative comparison. That is, all of the tags in the "source" are compared against the "destination"
// returning only tags that are common, and tags that are in the source, but not in the destination.
// The opposite list (tags that exist in the destination, but not the source) is not generated.
func ListTagDifference(ctx context.Context, db *database.Connection, slugs ListTagsDifferenceInput) (*dtos.TagDifference, error) {
	sourceTags, err := ListTagsForOperation(ctx, db, ListTagsForOperationInput{slugs.SourceOperationSlug})
	if err != nil {
		return nil, err
	}
	destinationTags, err := ListTagsForOperation(ctx, db, ListTagsForOperationInput{slugs.DestinationOperationSlug})
	if err != nil {
		return nil, err
	}

	srcTagNames := standardizeTagName(sourceTags)
	dstTagNames := standardizeTagName(destinationTags)

	var diff dtos.TagDifference
	for k, srcTag := range srcTagNames {
		dstTag, ok := dstTagNames[k]
		if ok {
			diff.Included = append(diff.Included, dtos.TagPair{SourceTag: *srcTag, DestinationTag: *dstTag})
		} else {
			diff.Excluded = append(diff.Excluded, *srcTag)
		}
	}

	return &diff, nil
}

func standardizeTagName(tags []*dtos.Tag) map[string]*dtos.Tag {
	names := make(map[string]*dtos.Tag)
	for _, tag := range tags {
		standardName := strings.ToLower(strings.TrimSpace(tag.Name))
		names[standardName] = tag
	}
	return names
}
