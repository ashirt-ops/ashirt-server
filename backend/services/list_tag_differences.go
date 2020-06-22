// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"strings"

	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/dtos"
	"github.com/theparanoids/ashirt/backend/policy"
)

type ListTagsDifferenceInput struct {
	SourceOperationSlug      string
	DestinationOperationSlug string
}

type ListTagDifferenceForEvidenceInput struct {
	ListTagsDifferenceInput
	SourceEvidenceUUID string
}

// ListTagDifference determines which tag values are common between two operations. This is done via
// relative comparison. That is, all of the tags in the "source" are compared against the "destination"
// returning only tags that are common, and tags that are in the source, but not in the destination.
// The opposite list (tags that exist in the destination, but not the source) is not generated.
func ListTagDifference(ctx context.Context, db *database.Connection, i ListTagsDifferenceInput) (*dtos.TagDifference, error) {
	sourceOperation, err := lookupOperation(db, i.SourceOperationSlug)
	if err != nil {
		return nil, err
	}
	destinationOperation, err := lookupOperation(db, i.DestinationOperationSlug)
	if err != nil {
		return nil, err
	}

	if err := policyRequireWithAdminBypass(ctx,
		policy.CanReadOperation{OperationID: sourceOperation.ID},
		policy.CanReadOperation{OperationID: destinationOperation.ID},
	); err != nil {
		return nil, backend.UnauthorizedReadErr(err)
	}

	sourceTags, err := listTagsForOperation(db, sourceOperation.ID)
	if err != nil {
		return nil, err
	}
	destinationTags, err := listTagsForOperation(db, destinationOperation.ID)
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

func ListTagDifferenceForEvidence(ctx context.Context, db *database.Connection, input ListTagDifferenceForEvidenceInput) (*dtos.TagDifference, error) {
	diff, err := ListTagDifference(ctx, db, input.ListTagsDifferenceInput)
	if err != nil {
		return nil, err
	}

	_, evidence, err := lookupOperationEvidence(db, input.SourceOperationSlug, input.SourceEvidenceUUID)
	if err != nil {
		return nil, err
	}

	tagMap, _, err := tagsForEvidenceByID(db, []int64{evidence.ID})
	if err != nil {
		return nil, err
	}

	updatedDiff := dtos.TagDifference{}
	for _, mappedTag := range tagMap[evidence.ID] {
		tagID := mappedTag.ID
		for _, tagpair := range diff.Included {
			if tagpair.SourceTag.ID == tagID {
				updatedDiff.Included = append(updatedDiff.Included, tagpair)
			}
		}
		for _, tag := range diff.Excluded {
			if tag.ID == tagID {
				updatedDiff.Excluded = append(updatedDiff.Excluded, tag)
			}
		}
	}

	return &updatedDiff, nil
}

func standardizeTagName(tags []*dtos.Tag) map[string]*dtos.Tag {
	names := make(map[string]*dtos.Tag)
	for _, tag := range tags {
		standardName := strings.ToLower(strings.TrimSpace(tag.Name))
		names[standardName] = tag
	}
	return names
}
