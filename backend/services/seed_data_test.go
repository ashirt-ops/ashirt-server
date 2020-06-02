// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	localConsts "github.com/theparanoids/ashirt/backend/authschemes/localauth/constants"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/logging"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
)

// TestSeedData is designed to allow a database-to-structure mapping. This is useful either for
// populating/seeding the database (see ApplyTo method), or alternatively, as acting as a source of
// truth for post-db operations.
type TestSeedData struct {
	APIKeys        []models.APIKey
	Findings       []models.Finding
	Evidences      []models.Evidence
	Users          []models.User
	Operations     []models.Operation
	Tags           []models.Tag
	UserOpMap      []models.UserOperationPermission
	TagEviMap      []models.TagEvidenceMap
	EviFindingsMap []models.EvidenceFindingMap
	Queries        []models.Query
}

// AllInitialTagIds is a (convenience) method version of the function TagIDsFromTags
func (seed TestSeedData) AllInitialTagIds() []int64 {
	return TagIDsFromTags(seed.Tags...)
}

// ApplyTo takes the configured TestSeedData and writes these values to the database.
func (seed TestSeedData) ApplyTo(t *testing.T, db *database.Connection) {
	systemLogger := logging.GetSystemLogger()
	systemLogger.Log("msg", "Applying seed data", "firstUser", seed.Users[0].FirstName)
	logging.SetSystemLogger(logging.NewNopLogger())
	defer logging.SetSystemLogger(systemLogger)
	err := db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.BatchInsert("users", len(seed.Users), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.Users[i].ID,
				"slug":       seed.Users[i].Slug,
				"first_name": seed.Users[i].FirstName,
				"last_name":  seed.Users[i].LastName,
				"email":      seed.Users[i].Email,
				"admin":      seed.Users[i].Admin,
				"headless":   seed.Users[i].Headless,
				"disabled":   seed.Users[i].Disabled,
				"created_at": seed.Users[i].CreatedAt,
				"updated_at": seed.Users[i].UpdatedAt,
				"deleted_at": seed.Users[i].DeletedAt,
			}
		})
		tx.BatchInsert("api_keys", len(seed.APIKeys), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.APIKeys[i].ID,
				"access_key": seed.APIKeys[i].AccessKey,
				"secret_key": seed.APIKeys[i].SecretKey,
				"user_id":    seed.APIKeys[i].UserID,
				"created_at": seed.APIKeys[i].CreatedAt,
				"updated_at": seed.APIKeys[i].UpdatedAt,
			}
		})
		tx.BatchInsert("auth_scheme_data", len(seed.Users), func(i int) map[string]interface{} {
			if seed.Users[i].DeletedAt != nil || seed.Users[i].Headless {
				return map[string]interface{}{}
			}
			return map[string]interface{}{
				"id":                 seed.Users[i].ID,
				"auth_scheme":        localConsts.Code,
				"user_key":           seed.Users[i].Slug,
				"user_id":            seed.Users[i].ID,
				"encrypted_password": "$2a$10$MLooRpCcdyyxoXwe3ZiCFuQZfsGeVC7TPCSyYhTs8Bl/sFPd4K67W", //aka "password"
				// "last_login":         nil,
				"created_at": seed.Users[i].CreatedAt,
				"updated_at": seed.Users[i].UpdatedAt,
			}
		})
		tx.BatchInsert("operations", len(seed.Operations), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.Operations[i].ID,
				"slug":       seed.Operations[i].Slug,
				"name":       seed.Operations[i].Name,
				"status":     seed.Operations[i].Status,
				"created_at": seed.Operations[i].CreatedAt,
				"updated_at": seed.Operations[i].UpdatedAt,
			}
		})
		tx.BatchInsert("user_operation_permissions", len(seed.UserOpMap), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"user_id":      seed.UserOpMap[i].UserID,
				"operation_id": seed.UserOpMap[i].OperationID,
				"role":         seed.UserOpMap[i].Role,
				"created_at":   seed.UserOpMap[i].CreatedAt,
				"updated_at":   seed.UserOpMap[i].UpdatedAt,
			}
		})
		tx.BatchInsert("tags", len(seed.Tags), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":           seed.Tags[i].ID,
				"operation_id": seed.Tags[i].OperationID,
				"name":         seed.Tags[i].Name,
				"color_name":   seed.Tags[i].ColorName,
				"created_at":   seed.Tags[i].CreatedAt,
				"updated_at":   seed.Tags[i].UpdatedAt,
			}
		})
		tx.BatchInsert("evidence", len(seed.Evidences), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":              seed.Evidences[i].ID,
				"uuid":            seed.Evidences[i].UUID,
				"operation_id":    seed.Evidences[i].OperationID,
				"operator_id":     seed.Evidences[i].OperatorID,
				"description":     seed.Evidences[i].Description,
				"content_type":    seed.Evidences[i].ContentType,
				"full_image_key":  seed.Evidences[i].FullImageKey,
				"thumb_image_key": seed.Evidences[i].ThumbImageKey,
				"occurred_at":     seed.Evidences[i].OccurredAt,
				"created_at":      seed.Evidences[i].CreatedAt,
				"updated_at":      seed.Evidences[i].UpdatedAt,
			}
		})
		tx.BatchInsert("findings", len(seed.Findings), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":              seed.Findings[i].ID,
				"uuid":            seed.Findings[i].UUID,
				"operation_id":    seed.Findings[i].OperationID,
				"ready_to_report": seed.Findings[i].ReadyToReport,
				"ticket_link":     seed.Findings[i].TicketLink,
				"category":        seed.Findings[i].Category,
				"title":           seed.Findings[i].Title,
				"description":     seed.Findings[i].Description,
				"created_at":      seed.Findings[i].CreatedAt,
				"updated_at":      seed.Findings[i].UpdatedAt,
			}
		})
		tx.BatchInsert("evidence_finding_map", len(seed.EviFindingsMap), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"evidence_id": seed.EviFindingsMap[i].EvidenceID,
				"finding_id":  seed.EviFindingsMap[i].FindingID,
				"created_at":  seed.EviFindingsMap[i].CreatedAt,
				"updated_at":  seed.EviFindingsMap[i].UpdatedAt,
			}
		})
		tx.BatchInsert("tag_evidence_map", len(seed.TagEviMap), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"tag_id":      seed.TagEviMap[i].TagID,
				"evidence_id": seed.TagEviMap[i].EvidenceID,
				"created_at":  seed.TagEviMap[i].CreatedAt,
				"updated_at":  seed.TagEviMap[i].UpdatedAt,
			}
		})
		tx.BatchInsert("queries", len(seed.Queries), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":           seed.Queries[i].ID,
				"operation_id": seed.Queries[i].OperationID,
				"name":         seed.Queries[i].Name,
				"query":        seed.Queries[i].Query,
				"type":         seed.Queries[i].Type,
				"created_at":   seed.Queries[i].CreatedAt,
				"updated_at":   seed.Queries[i].UpdatedAt,
			}
		})
	})

	require.Nil(t, err)
}

func (seed TestSeedData) EvidenceIDsForFinding(finding models.Finding) []int64 {
	rtn := make([]int64, 0)
	for _, row := range seed.EviFindingsMap {
		if row.FindingID == finding.ID {
			rtn = append(rtn, row.EvidenceID)
		}
	}
	return rtn
}

func (seed TestSeedData) TagsForFinding(finding models.Finding) []models.Tag {
	evidenceIDs := seed.EvidenceIDsForFinding(finding)
	rtn := make([]models.Tag, 0)
	for _, eviID := range evidenceIDs {
		subset := seed.tagsForEvidenceID(eviID)
		rtn = append(rtn, subset...)
	}
	return rtn
}

func (seed TestSeedData) TagsForEvidence(evidence models.Evidence) []models.Tag {
	return seed.tagsForEvidenceID(evidence.ID)
}

func (seed TestSeedData) tagsForEvidenceID(eviID int64) []models.Tag {
	rtn := make([]models.Tag, 0)
	for _, row := range seed.TagEviMap {
		if eviID == row.EvidenceID {
			rtn = append(rtn, seed.GetTagFromID(row.TagID))
		}
	}
	return rtn
}

func (seed TestSeedData) GetTagFromID(id int64) models.Tag {
	for _, item := range seed.Tags {
		if item.ID == id {
			return item
		}
	}
	return models.Tag{}
}

func (seed TestSeedData) GetUserFromID(id int64) models.User {
	for _, item := range seed.Users {
		if item.ID == id {
			return item
		}
	}
	return models.User{}
}

func (seed TestSeedData) UsersForOp(op models.Operation) []models.User {
	rtn := make([]models.User, 0)

	for _, row := range seed.UserOpMap {
		if row.OperationID == op.ID {
			rtn = append(rtn, seed.GetUserFromID(row.UserID))
		}
	}
	return rtn
}

func (seed TestSeedData) UserRoleForOp(user models.User, op models.Operation) policy.OperationRole {
	for _, row := range seed.UserOpMap {
		if row.OperationID == op.ID && row.UserID == user.ID {
			return row.Role
		}
	}
	return ""
}
