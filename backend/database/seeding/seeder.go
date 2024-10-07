package seeding

import (
	"context"
	"strings"
	"time"

	localConsts "github.com/ashirt-ops/ashirt-server/backend/authschemes/localauth/constants"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"golang.org/x/crypto/bcrypt"
)

// Seeder is designed to allow a database-to-structure mapping. This is useful either for
// populating/seeding the database (see ApplyTo method), or alternatively, as acting as a source of
// truth for post-db operations.
type Seeder struct {
	FindingCategories []models.FindingCategory
	APIKeys           []models.APIKey
	Findings          []models.Finding
	Evidences         []models.Evidence
	EvidenceMetadatas []models.EvidenceMetadata
	Users             []models.User
	UserGroups        []models.UserGroup
	UserGroupMaps     []models.UserGroupMap
	Operations        []models.Operation
	DefaultTags       []models.DefaultTag
	Tags              []models.Tag
	UserOpMap         []models.UserOperationPermission
	UserGroupOpMap    []models.UserGroupOperationPermission
	UserOpPrefMap     []models.UserOperationPreferences
	TagEviMap         []models.TagEvidenceMap
	EviFindingsMap    []models.EvidenceFindingMap
	Queries           []models.Query
	ServiceWorkers    []models.ServiceWorker
	GlobalVars        []models.GlobalVar
	OperationVars     []models.OperationVar
	VarOperationMap   []models.VarOperationMap
}

// AllInitialTagIds is a (convenience) method version of the function TagIDsFromTags
func (seed Seeder) AllInitialTagIds() []int64 {
	return TagIDsFromTags(seed.Tags...)
}

// AllInitialDefaultTagIds is a (convenience) method version of the function DefaultTagIDsFromTags
func (seed Seeder) AllInitialDefaultTagIds() []int64 {
	return DefaultTagIDsFromTags(seed.DefaultTags...)
}

// ApplyTo takes the configured Seeder and writes these values to the database.
func (seed Seeder) ApplyTo(db *database.Connection) error {
	if len(seed.Users) == 0 {
		logging.GetSystemLogger().Warn("Applying NO seed data")
		return nil
	}

	systemLogger := logging.GetSystemLogger()
	systemLogger.Info("Applying seed data", "firstUser", seed.Users[0].FirstName)
	logging.SetSystemLogger(logging.NewNopLogger())
	defer logging.SetSystemLogger(systemLogger)
	err := db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.BatchInsert("finding_categories", len(seed.FindingCategories), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.FindingCategories[i].ID,
				"category":   seed.FindingCategories[i].Category,
				"created_at": seed.FindingCategories[i].CreatedAt,
				"updated_at": seed.FindingCategories[i].UpdatedAt,
				"deleted_at": seed.FindingCategories[i].DeletedAt,
			}
		})
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
		tx.BatchInsert("user_groups", len(seed.UserGroups), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.UserGroups[i].ID,
				"slug":       seed.UserGroups[i].Slug,
				"name":       seed.UserGroups[i].Name,
				"created_at": seed.UserGroups[i].CreatedAt,
				"updated_at": seed.UserGroups[i].UpdatedAt,
				"deleted_at": seed.UserGroups[i].DeletedAt,
			}
		})
		tx.BatchInsert("group_user_map", len(seed.UserGroupMaps), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"group_id":   seed.UserGroupMaps[i].GroupID,
				"user_id":    seed.UserGroupMaps[i].UserID,
				"created_at": seed.UserGroupMaps[i].CreatedAt,
				"updated_at": seed.UserGroupMaps[i].UpdatedAt,
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
			encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(strings.ToLower(seed.Users[i].FirstName)), bcrypt.DefaultCost)
			return map[string]interface{}{
				"id":                 seed.Users[i].ID,
				"auth_scheme":        localConsts.Code,
				"auth_type":          localConsts.Code,
				"username":           seed.Users[i].FirstName,
				"user_id":            seed.Users[i].ID,
				"encrypted_password": encryptedPassword, //the user's first name, lowercased
				"created_at":         seed.Users[i].CreatedAt,
				"updated_at":         seed.Users[i].UpdatedAt,
			}
		})
		tx.BatchInsert("operations", len(seed.Operations), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.Operations[i].ID,
				"slug":       seed.Operations[i].Slug,
				"name":       seed.Operations[i].Name,
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
		tx.BatchInsert("user_group_operation_permissions", len(seed.UserGroupOpMap), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"group_id":     seed.UserGroupOpMap[i].UserGroupID,
				"operation_id": seed.UserGroupOpMap[i].OperationID,
				"role":         seed.UserGroupOpMap[i].Role,
				"created_at":   seed.UserGroupOpMap[i].CreatedAt,
				"updated_at":   seed.UserGroupOpMap[i].UpdatedAt,
			}
		})
		tx.BatchInsert("user_operation_preferences", len(seed.UserOpPrefMap), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"user_id":      seed.UserOpPrefMap[i].UserID,
				"operation_id": seed.UserOpPrefMap[i].OperationID,
				"created_at":   seed.UserOpPrefMap[i].CreatedAt,
				"updated_at":   seed.UserOpPrefMap[i].UpdatedAt,
				"is_favorite":  seed.UserOpPrefMap[i].IsFavorite,
			}
		})
		tx.BatchInsert("default_tags", len(seed.DefaultTags), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.DefaultTags[i].ID,
				"name":       seed.DefaultTags[i].Name,
				"color_name": seed.DefaultTags[i].ColorName,
				"created_at": seed.DefaultTags[i].CreatedAt,
				"updated_at": seed.DefaultTags[i].UpdatedAt,
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
		tx.BatchInsert("evidence_metadata", len(seed.EvidenceMetadatas), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":          seed.EvidenceMetadatas[i].ID,
				"evidence_id": seed.EvidenceMetadatas[i].EvidenceID,
				"source":      seed.EvidenceMetadatas[i].Source,
				"can_process": seed.EvidenceMetadatas[i].CanProcess,
				"body":        seed.EvidenceMetadatas[i].Body,
				"created_at":  seed.EvidenceMetadatas[i].CreatedAt,
				"updated_at":  seed.EvidenceMetadatas[i].UpdatedAt,
			}
		})
		tx.BatchInsert("findings", len(seed.Findings), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":              seed.Findings[i].ID,
				"uuid":            seed.Findings[i].UUID,
				"operation_id":    seed.Findings[i].OperationID,
				"ready_to_report": seed.Findings[i].ReadyToReport,
				"ticket_link":     seed.Findings[i].TicketLink,
				"category_id":     seed.Findings[i].CategoryID,
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
		tx.BatchInsert("service_workers", len(seed.ServiceWorkers), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.ServiceWorkers[i].ID,
				"name":       seed.ServiceWorkers[i].Name,
				"config":     seed.ServiceWorkers[i].Config,
				"created_at": seed.ServiceWorkers[i].CreatedAt,
				"updated_at": seed.ServiceWorkers[i].UpdatedAt,
				"deleted_at": seed.ServiceWorkers[i].DeletedAt,
			}
		})
		tx.BatchInsert("global_vars", len(seed.GlobalVars), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.GlobalVars[i].ID,
				"name":       seed.GlobalVars[i].Name,
				"value":      seed.GlobalVars[i].Value,
				"created_at": seed.GlobalVars[i].CreatedAt,
				"updated_at": seed.GlobalVars[i].UpdatedAt,
			}
		})
		tx.BatchInsert("operation_vars", len(seed.OperationVars), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"id":         seed.OperationVars[i].ID,
				"slug":       seed.OperationVars[i].Slug,
				"name":       seed.OperationVars[i].Name,
				"value":      seed.OperationVars[i].Value,
				"created_at": seed.OperationVars[i].CreatedAt,
				"updated_at": seed.OperationVars[i].UpdatedAt,
			}
		})
		tx.BatchInsert("var_operation_map", len(seed.VarOperationMap), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"var_id":       seed.VarOperationMap[i].VarID,
				"operation_id": seed.VarOperationMap[i].OperationID,
				"created_at":   seed.VarOperationMap[i].CreatedAt,
				"updated_at":   seed.VarOperationMap[i].UpdatedAt,
			}
		})
	})

	return err
}

func (seed Seeder) Reset(db *database.Connection) error {
	err := ClearDB(db)
	if err != nil {
		return err
	}
	return seed.ApplyTo(db)
}

func (seed Seeder) CategoryForFinding(finding models.Finding) string {
	if finding.CategoryID == nil {
		return ""
	}
	for _, row := range seed.FindingCategories {
		if row.ID == *finding.CategoryID {
			return row.Category
		}
	}
	return ""
}

func (seed Seeder) EvidenceIDsForFinding(finding models.Finding) []int64 {
	rtn := make([]int64, 0)
	for _, row := range seed.EviFindingsMap {
		if row.FindingID == finding.ID {
			rtn = append(rtn, row.EvidenceID)
		}
	}
	return rtn
}

func (seed Seeder) TagsForFinding(finding models.Finding) []models.Tag {
	evidenceIDs := seed.EvidenceIDsForFinding(finding)
	rtn := make([]models.Tag, 0)
	for _, eviID := range evidenceIDs {
		subset := seed.tagsForEvidenceID(eviID)
		rtn = append(rtn, subset...)
	}
	return rtn
}

func (seed Seeder) TagsForEvidence(evidence models.Evidence) []models.Tag {
	return seed.tagsForEvidenceID(evidence.ID)
}

func (seed Seeder) tagsForEvidenceID(eviID int64) []models.Tag {
	rtn := make([]models.Tag, 0)
	for _, row := range seed.TagEviMap {
		if eviID == row.EvidenceID {
			rtn = append(rtn, seed.GetTagFromID(row.TagID))
		}
	}
	return rtn
}

func (seed Seeder) GetTagFromID(id int64) models.Tag {
	for _, item := range seed.Tags {
		if item.ID == id {
			return item
		}
	}
	return models.Tag{}
}

func (seed Seeder) GetUserFromID(id int64) models.User {
	for _, item := range seed.Users {
		if item.ID == id {
			return item
		}
	}
	return models.User{}
}

func (seed Seeder) GetUserGroupFromID(id int64) models.UserGroup {
	for _, item := range seed.UserGroups {
		if item.ID == id {
			return item
		}
	}
	return models.UserGroup{}
}

func (seed Seeder) UsersForOp(op models.Operation) []models.User {
	rtn := make([]models.User, 0)

	for _, row := range seed.UserOpMap {
		if row.OperationID == op.ID {
			rtn = append(rtn, seed.GetUserFromID(row.UserID))
		}
	}
	return rtn
}

func (seed Seeder) UserGroupsForOp(op models.Operation) []models.UserGroup {
	rtn := make([]models.UserGroup, 0)

	for _, row := range seed.UserGroupOpMap {
		if row.OperationID == op.ID {
			rtn = append(rtn, seed.GetUserGroupFromID(row.UserGroupID))
		}
	}
	return rtn
}

func (seed Seeder) UserRoleForOp(user models.User, op models.Operation) policy.OperationRole {
	for _, row := range seed.UserOpMap {
		if row.OperationID == op.ID && row.UserID == user.ID {
			return row.Role
		}
	}
	return ""
}

func (seed Seeder) UserGroupRoleForOp(userGroup models.UserGroup, op models.Operation) policy.OperationRole {
	for _, row := range seed.UserGroupOpMap {
		if row.OperationID == op.ID && row.UserGroupID == userGroup.ID {
			return row.Role
		}
	}
	return ""
}

func (seed Seeder) EvidenceForOperation(opID int64) []models.Evidence {
	evidence := make([]models.Evidence, 0)
	for _, row := range seed.Evidences {
		if row.OperationID == opID {
			evidence = append(evidence, row)
		}
	}
	return evidence
}

func (seed Seeder) TagIDsUsageByDate(opID int64) map[int64][]time.Time {
	evidence := seed.EvidenceForOperation(opID)
	tagIDUsageMap := make(map[int64][]time.Time)

	for _, evi := range evidence {
		tags := seed.TagsForEvidence(evi)
		for _, t := range tags {
			tmp := tagIDUsageMap[t.ID]
			tmp = append(tmp, evi.OccurredAt)
			tagIDUsageMap[t.ID] = tmp
		}
	}

	return tagIDUsageMap
}

func (seed Seeder) OperationForEvidence(evidence models.Evidence) *models.Operation {
	for _, op := range seed.Operations {
		if op.ID == evidence.OperationID {
			return &op
		}
	}
	return nil
}
