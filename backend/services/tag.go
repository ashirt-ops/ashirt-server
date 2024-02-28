package services

import (
	"context"
	"strings"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateTagInput struct {
	Name          string
	ColorName     string
	OperationSlug string
	Description   *string
}

type CreateDefaultTagInput struct {
	Name        string
	ColorName   string
	Description *string
}

type DeleteTagInput struct {
	ID            int64
	OperationSlug string
}

type DeleteDefaultTagInput struct {
	ID int64
}

type ListTagsDifferenceInput struct {
	SourceOperationSlug      string
	DestinationOperationSlug string
}

type ListTagDifferenceForEvidenceInput struct {
	ListTagsDifferenceInput
	SourceEvidenceUUID string
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

type ListTagsForOperationInput struct {
	OperationSlug string
}

type UpdateDefaultTagInput struct {
	ID          int64
	Name        string
	ColorName   string
	Description *string
}

type UpdateTagInput struct {
	ID            int64
	OperationSlug string
	Name          string
	ColorName     string
	Description   *string
}

func CreateTag(ctx context.Context, db *database.Connection, i CreateTagInput) (*dtos.Tag, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to create tag", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyTagsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unable to create tag", backend.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, backend.MissingValueErr("Name")
	}

	tagID, err := db.Insert("tags", map[string]interface{}{
		"name":         i.Name,
		"color_name":   i.ColorName,
		"operation_id": operation.ID,
		"description":  i.Description,
	})
	if err != nil {
		return nil, backend.WrapError("Cannot add new tag", backend.DatabaseErr(err))
	}
	return &dtos.Tag{
		ID:          tagID,
		Name:        i.Name,
		ColorName:   i.ColorName,
		Description: i.Description,
	}, nil
}

// CreateDefaultTag creates a single tag in the default_tags table. Admin only.
func CreateDefaultTag(ctx context.Context, db *database.Connection, i CreateDefaultTagInput) (*dtos.DefaultTag, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.WrapError("Unable to create default tag", backend.UnauthorizedWriteErr(err))
	}

	if i.Name == "" {
		return nil, backend.MissingValueErr("Name")
	}

	tagID, err := db.Insert("default_tags", map[string]interface{}{
		"name":       i.Name,
		"color_name": i.ColorName,
	})
	if err != nil {
		return nil, backend.WrapError("Cannot add new tag", backend.DatabaseErr(err))
	}
	return &dtos.DefaultTag{
		ID:          tagID,
		Name:        i.Name,
		ColorName:   i.ColorName,
		Description: i.Description,
	}, nil
}

func MergeDefaultTags(ctx context.Context, db *database.Connection, i []CreateDefaultTagInput) error {
	if err := policyRequireWithAdminBypass(ctx, policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Unwilling to update default tag", backend.UnauthorizedWriteErr(err))
	}

	tagsToInsert := make([]CreateDefaultTagInput, 0, len(i))
	currentTagNames := make([]string, 0, len(i))

	for _, t := range i {
		if listContainsString(currentTagNames, t.Name) != -1 || t.Name == "" {
			continue // no need to re-process a tag if we've dealt with it -- just use the first instance
		} else {
			currentTagNames = append(currentTagNames, t.Name)
		}

		tagsToInsert = append(tagsToInsert, t)
	}

	err := db.BatchInsert("default_tags", len(tagsToInsert), func(idx int) map[string]interface{} {
		return map[string]interface{}{
			"name":       tagsToInsert[idx].Name,
			"color_name": tagsToInsert[idx].ColorName,
		}
	}, "ON DUPLICATE KEY UPDATE color_name=VALUES(color_name)")

	if err != nil {
		return backend.WrapError("Cannot update default tag", backend.DatabaseErr(err))
	}
	return nil
}

// DeleteTag removes a tag and untags all evidence with the tag
func DeleteTag(ctx context.Context, db *database.Connection, i DeleteTagInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to delete tag", backend.UnauthorizedWriteErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyTagsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to delete tag", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Delete(sq.Delete("tag_evidence_map").Where(sq.Eq{"tag_id": i.ID}))
		tx.Delete(sq.Delete("tags").Where(sq.Eq{"id": i.ID}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete tag", backend.DatabaseErr(err))
	}

	return nil
}

// DeleteDefaultTag removes a single tag in the default_tags table by the tag id. Admin only.
func DeleteDefaultTag(ctx context.Context, db *database.Connection, i DeleteDefaultTagInput) error {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Unwilling to delete default tag", backend.UnauthorizedWriteErr(err))
	}

	err := db.Delete(sq.Delete("default_tags").Where(sq.Eq{"id": i.ID}))
	if err != nil {
		return backend.WrapError("Cannot delete default tag", backend.DatabaseErr(err))
	}

	return nil
}

// ListTagDifference determines which tag values are common between two operations. This is done via
// relative comparison. That is, all of the tags in the "source" are compared against the "destination"
// returning only tags that are common, and tags that are in the source, but not in the destination.
// The opposite list (tags that exist in the destination, but not the source) is not generated.
func ListTagDifference(ctx context.Context, db *database.Connection, i ListTagsDifferenceInput) (*dtos.TagDifference, error) {
	sourceOperation, err := lookupOperation(db, i.SourceOperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list tag differences", err)
	}
	destinationOperation, err := lookupOperation(db, i.DestinationOperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list tag differences", err)
	}

	if err := policyRequireWithAdminBypass(ctx,
		policy.CanReadOperation{OperationID: sourceOperation.ID},
		policy.CanReadOperation{OperationID: destinationOperation.ID},
	); err != nil {
		return nil, backend.WrapError("Unwilling to list tag differences", backend.UnauthorizedReadErr(err))
	}

	sourceTags, err := listTagsForOperation(db, sourceOperation.ID)
	if err != nil {
		return nil, backend.WrapError("Cannot list tag differences", err)
	}
	destinationTags, err := listTagsForOperation(db, destinationOperation.ID)
	if err != nil {
		return nil, backend.WrapError("Cannot list tag differences", err)
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
		return nil, backend.WrapError("Unable to list tag difference", err)
	}

	_, evidence, err := lookupOperationEvidence(db, input.SourceOperationSlug, input.SourceEvidenceUUID)
	if err != nil {
		return nil, backend.WrapError("Unable to get evidence for tagdiff", err)
	}

	tagMap, _, err := tagsForEvidenceByID(db, []int64{evidence.ID})
	if err != nil {
		return nil, backend.WrapError("Cannot get evidence tags", err)
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

func standardizeTagName(tags []*dtos.TagWithUsage) map[string]*dtos.Tag {
	names := make(map[string]*dtos.Tag)
	for _, tag := range tags {
		standardName := strings.ToLower(strings.TrimSpace(tag.Name))
		names[standardName] = &tag.Tag
	}
	return names
}

func listContainsString(haystack []string, needle string) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}

func ListTagsForOperation(ctx context.Context, db *database.Connection, i ListTagsForOperationInput) ([]*dtos.TagWithUsage, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list tags for operation", backend.UnauthorizedReadErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list tags for operation", backend.UnauthorizedReadErr(err))
	}

	return listTagsForOperation(db, operation.ID)
}

// listTagsForOperation generates a list tags associted with a given operation. This does not
// check permission, and so is not exported, and is intended to only be used as a helper method
// for other services
func listTagsForOperation(db *database.Connection, operationID int64) ([]*dtos.TagWithUsage, error) {
	type DBTag struct {
		models.Tag
		TagCount int64 `db:"tag_count"`
	}
	var tags []DBTag
	err := db.Select(&tags, sq.Select("tags.*").Column("count(tag_id) AS tag_count").
		From("tags").
		LeftJoin("tag_evidence_map ON tag_evidence_map.tag_id = tags.id").
		Where(sq.Eq{"operation_id": operationID}).
		GroupBy("tags.id").
		OrderBy("tags.id ASC"))
	if err != nil {
		return nil, backend.WrapError("Cannot get tags for operation", backend.DatabaseErr(err))
	}

	tagsDTO := make([]*dtos.TagWithUsage, len(tags))
	for idx, tag := range tags {
		tagsDTO[idx] = &dtos.TagWithUsage{
			Tag: dtos.Tag{
				ID:          tag.Tag.ID,
				Name:        tag.Tag.Name,
				ColorName:   tag.Tag.ColorName,
				Description: tag.Tag.Description,
			},
			EvidenceCount: tag.TagCount,
		}
	}
	return tagsDTO, nil
}

// ListDefaultTags provides a list of all of the tags in the default_tags table. Admin only.
func ListDefaultTags(ctx context.Context, db *database.Connection) ([]*dtos.DefaultTag, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.WrapError("Unwilling to list default tags", backend.UnauthorizedReadErr(err))
	}

	var tags []models.Tag
	err := db.Select(&tags, sq.Select("id", "name", "color_name").From("default_tags"))

	if err != nil {
		return nil, backend.WrapError("Cannot get default tags", backend.DatabaseErr(err))
	}

	tagsDTO := make([]*dtos.DefaultTag, len(tags))
	for idx, tag := range tags {
		tagsDTO[idx] = &dtos.DefaultTag{
			ID:          tag.ID,
			Name:        tag.Name,
			ColorName:   tag.ColorName,
			Description: tag.Description,
		}
	}
	return tagsDTO, nil
}

// UpdateTag updates a tag's name and color
func UpdateTag(ctx context.Context, db *database.Connection, i UpdateTagInput) error {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return backend.WrapError("Unable to update tag", backend.UnauthorizedWriteErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanModifyTagsOfOperation{OperationID: operation.ID}); err != nil {
		return backend.WrapError("Unwilling to update tag", backend.UnauthorizedWriteErr(err))
	}

	err = db.Update(sq.Update("tags").
		SetMap(map[string]interface{}{
			"name":        i.Name,
			"color_name":  i.ColorName,
			"description": i.Description,
		}).
		Where(sq.Eq{"id": i.ID}))

	if err != nil {
		return backend.WrapError("Cannot update tag", backend.DatabaseErr(err))
	}
	return nil
}

func UpdateDefaultTag(ctx context.Context, db *database.Connection, i UpdateDefaultTagInput) error {
	if err := policyRequireWithAdminBypass(ctx, policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Unwilling to update default tag", backend.UnauthorizedWriteErr(err))
	}

	err := db.Update(sq.Update("default_tags").
		SetMap(map[string]interface{}{
			"name":        i.Name,
			"color_name":  i.ColorName,
			"description": i.Description,
		}).
		Where(sq.Eq{"id": i.ID}))

	if err != nil {
		return backend.WrapError("Cannot update default tag", backend.DatabaseErr(err))
	}
	return nil
}
