// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package seeding

import (
	"context"
	"os"
	"sort"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jonboulle/clockwork"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"
)

var internalClock clockwork.Clock = clockwork.NewFakeClock()

// SetInternalClock provides a mechanism to override the default/testing clock and use a custom
// clock. There are two options: clockwork.NewRealClock() (for current time) or
// clockwork.NewFakeClockAt() for a fake clock at a given time.
//
// Note: this should be called _before_ seeding the database
func SetInternalClock(newClock clockwork.Clock) {
	internalClock = newClock
}

func GetInternalClock() clockwork.Clock {
	return internalClock
}

// ContextForUser genereates a user's context as if they had just logged in. All settings are set,
// except for NeedsReset, which is always false
func ContextForUser(my models.User, db *database.Connection) context.Context {
	ctx := context.Background()
	return middleware.BuildContextForUser(ctx, db, my.ID, my.Admin)
}

// IsSeeded does a check against the database to see if any users are registered. if no users are
// registered, then it is assumed that the database has not been seeded.
func IsSeeded(db *database.Connection) (bool, error) {
	var count int64
	err := db.Get(&count, sq.Select("count(id)").From("users"))

	return count > 0, err
}

// TagIDsFromTags maps over models.Tags to come up with a collection of IDs for those tags
// equivalent js: tags.map( i => i.ID)
func TagIDsFromTags(tags ...models.Tag) []int64 {
	ids := make([]int64, len(tags))
	for i, t := range tags {
		ids[i] = t.ID
	}
	return ids
}

func newAPIKeyGen(first int64) func(int64, string, []byte) models.APIKey {
	id := iotaLike(first)
	return func(userID int64, accessKey string, secretKey []byte) models.APIKey {
		return models.APIKey{
			ID:        id(),
			UserID:    userID,
			AccessKey: accessKey,
			SecretKey: secretKey,
			CreatedAt: internalClock.Now(),
		}
	}
}

type newUserInput struct {
	FirstName      string
	LastName       string
	Birthday       time.Time
	SetLastUpdated bool
	IsAdmin        bool
	Disabled       bool
	Deleted        bool
	Headless       bool
}

func newUserGen(first int64, toSlug func(f, l string) string) func(input newUserInput) models.User {
	id := iotaLike(first)
	return func(input newUserInput) models.User {
		user := models.User{
			ID:        id(),
			Slug:      toSlug(input.FirstName, input.LastName),
			FirstName: strings.Title(input.FirstName),
			LastName:  strings.Title(input.LastName),
			Email:     toSlug(input.FirstName, input.LastName),
			Admin:     input.IsAdmin,
			Disabled:  input.Disabled,
			CreatedAt: input.Birthday,
			Headless:  input.Headless,
		}
		if input.SetLastUpdated {
			updatedDate := time.Date(input.Birthday.Year()+1, input.Birthday.Month(), input.Birthday.Day(), 0, 0, 0, 0, time.UTC)
			user.UpdatedAt = &updatedDate
		}
		if input.Deleted {
			deletedDate := time.Date(input.Birthday.Year()+1, input.Birthday.Month(), input.Birthday.Day(), 0, 0, 0, 0, time.UTC)
			user.DeletedAt = &deletedDate
		}
		return user
	}
}

func newTagGen(first int64) func(opID int64, name, colorName string) models.Tag {
	id := iotaLike(first)
	return func(opID int64, name, colorName string) models.Tag {
		return models.Tag{
			ID:          id(),
			OperationID: opID,
			Name:        name,
			ColorName:   colorName,
			CreatedAt:   internalClock.Now(),
		}
	}
}

func newOperationGen(first int64) func(slug, fullName string) models.Operation {
	id := iotaLike(first)
	return func(slug, fullName string) models.Operation {
		return models.Operation{
			ID:        id(),
			Slug:      slug,
			Name:      fullName,
			Status:    models.OperationStatusPlanning,
			CreatedAt: internalClock.Now(),
		}
	}
}

func newEvidenceGen(first int64) func(opID, ownerID int64, uuid, desc, contentType string, clockDayOffset int) models.Evidence {
	id := iotaLike(first)
	return func(opID, ownerID int64, uuid, desc, contentType string, clockDayOffset int) models.Evidence {
		return models.Evidence{
			ID:            id(),
			UUID:          uuid,
			OperationID:   opID,
			OperatorID:    ownerID,
			Description:   desc,
			ContentType:   contentType,
			FullImageKey:  uuid,
			ThumbImageKey: uuid,
			OccurredAt:    internalClock.Now().AddDate(0, 0, clockDayOffset),
			CreatedAt:     internalClock.Now(),
		}
	}
}

func newFindingCategoryGen(first int64) func(category string, deleted bool) models.FindingCategory {
	id := iotaLike(first)
	return func(category string, deleted bool) models.FindingCategory {
		findingCategory := models.FindingCategory{
			ID:        id(),
			Category:  category,
			CreatedAt: internalClock.Now(),
		}

		if deleted {
			deletedDate := internalClock.Now()
			findingCategory.DeletedAt = &deletedDate
		}
		return findingCategory
	}
}

func newFindingGen(first int64) func(opID int64, uuid string, category *int64, title, desc string, ticketLink *string) models.Finding {
	id := iotaLike(first)
	return func(opID int64, uuid string, category *int64, title, desc string, ticketLink *string) models.Finding {
		finding := models.Finding{
			ID:            id(),
			OperationID:   opID,
			UUID:          uuid,
			CategoryID:    category,
			Title:         title,
			Description:   desc,
			ReadyToReport: (ticketLink != nil),
			CreatedAt:     internalClock.Now(),
		}
		if finding.ReadyToReport && *ticketLink != "" {
			finding.TicketLink = ticketLink
		}
		return finding
	}
}

func newUserOpPermission(user models.User, op models.Operation, role policy.OperationRole) models.UserOperationPermission {
	return models.UserOperationPermission{
		UserID:      user.ID,
		OperationID: op.ID,
		Role:        role,
		CreatedAt:   internalClock.Now(),
	}
}

func newQueryGen(first int64) func(opID int64, name, query, qType string) models.Query {
	id := iotaLike(first)
	return func(opID int64, name, query, qType string) models.Query {
		return models.Query{
			ID:          id(),
			OperationID: opID,
			Name:        name,
			Query:       query,
			Type:        qType,
			CreatedAt:   internalClock.Now(),
		}
	}
}

// associateEvidenceToTag mirrors associateTagsToEvidence. Rather than associating multiple tags
// with a single piece of evidence this will instead associate a single tag to multiple evidence.
func associateEvidenceToTag(tag models.Tag, evis ...models.Evidence) []models.TagEvidenceMap {
	mappings := make([]models.TagEvidenceMap, 0, len(evis))
	for _, evi := range evis {
		if evi.OperationID == tag.OperationID {
			mappings = append(mappings, models.TagEvidenceMap{TagID: tag.ID, EvidenceID: evi.ID, CreatedAt: internalClock.Now()})
		} else {
			// will likely be ignored, but helpful in constructing new sets
			os.Stderr.WriteString("[Testing - WARNING] Trying to associate tag(" + tag.Name + ") with evidence(" + evi.UUID + ") in differeing operations\n")
		}
	}
	return mappings
}

func associateTagsToEvidence(evi models.Evidence, tags ...models.Tag) []models.TagEvidenceMap {
	mappings := make([]models.TagEvidenceMap, 0, len(tags))

	for _, t := range tags {
		if t.OperationID == evi.OperationID {
			mappings = append(mappings, models.TagEvidenceMap{TagID: t.ID, EvidenceID: evi.ID, CreatedAt: internalClock.Now()})
		} else {
			// will likely be ignored, but helpful in constructing new sets
			os.Stderr.WriteString("[Testing - WARNING] Trying to associate tag(" + t.Name + ") with evidence(" + evi.UUID + ") in differeing operations\n")
		}
	}
	return mappings
}

func associateEvidenceToFinding(finding models.Finding, evi ...models.Evidence) []models.EvidenceFindingMap {
	mappings := make([]models.EvidenceFindingMap, 0, len(evi))

	for _, e := range evi {
		if e.OperationID == finding.OperationID {
			mappings = append(mappings, models.EvidenceFindingMap{EvidenceID: e.ID, FindingID: finding.ID, CreatedAt: internalClock.Now()})
		} else {
			// will likely be ignored, but helpful in constructing new sets
			os.Stderr.WriteString("[Testing - WARNING] Trying to associate evidence(" + e.UUID + ") with finding(" + finding.Title + ") in differeing operations\n")
		}
	}
	return mappings
}

func unionTagEviMap(parts ...[]models.TagEvidenceMap) []models.TagEvidenceMap {
	totalLength := 0
	for _, p := range parts {
		totalLength += len(p)
	}
	result := make([]models.TagEvidenceMap, totalLength)
	copied := 0
	for _, part := range parts {
		copied += copy(result[copied:], part)
	}
	return result
}

func unionEviFindingMap(parts ...[]models.EvidenceFindingMap) []models.EvidenceFindingMap {
	totalLength := 0
	for _, p := range parts {
		totalLength += len(p)
	}
	result := make([]models.EvidenceFindingMap, totalLength)
	copied := 0
	for _, part := range parts {
		copied += copy(result[copied:], part)
	}
	return result
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// iotaLike produces an integer iterator.
func iotaLike(start int64) func() int64 {
	counter := start
	return func() int64 {
		rtn := counter
		counter++
		return rtn
	}
}

// Sorted orders an int slice in asc order, then returns back a copy of the sorted list
// note: underlying list is unedited
func Sorted(slice []int64) []int64 {
	clone := make([]int64, len(slice))
	copy(clone, slice)
	sort.Slice(clone, func(i, j int) bool { return clone[i] < clone[j] })
	return clone
}
