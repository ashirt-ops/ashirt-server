package seeding

import (
	"bytes"
	"context"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/policy"
	"github.com/ashirt-ops/ashirt-server/backend/server/middleware"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"
)

// TinyImg is the smallest png. Used for testing. Reference: https://github.com/mathiasbynens/small
var TinyImg []byte = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
	0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
	0x42, 0x60, 0x82,
}

// TinyCodeblock is a minimal codeblock. Used for testing.
var TinyCodeblock []byte = []byte(`{"contentSubtype": "python", "content": "print(\"Hello World!\")"}`)

// TinyTermRec is a minimal terminal recording. Used for testing.
var TinyTermRec []byte = []byte(
	`{"version":2,"width":75,"height":18,"timestamp":1593020879,"title":"1593020879","env":{"SHELL":"/bin/bash","TERM":"xterm-256color"}}` +
		"\n" + `[0.188801409,"o","\u001b]0;user@localhost:~\u0007"]` +
		"\n" + `[0.189032775,"o","[user@localhost ~]$ "]` +
		"\n" + `[0.716089612,"o","ll\r\n"]` +
		"\n" + `[1.061539838,"o","total 10652\r\n"]` +
		"\n" + `[1.061654704,"o","-rwxrwxr-x. 1 user user 10905365 Jun 24 10:41 \u001b[0m\u001b[38;5;40mtermrec\u001b[0m\r\n"]` +
		"\n" + `[1.062881589,"o","\u001b]0;user@localhost:~\u0007"]` +
		"\n" + `[1.063084503,"o","[user@localhost ~]$ "]` +
		"\n" + `[1.517546751,"o","exit\r\n"]` +
		"\n" + `[2.129344227,"o","exit\r\n"]`,
)

// InitTest creates a connection to the database and provides an established, but otherwise empty
// database. The database is named "service-test-db"
func InitTest(t *testing.T) *database.Connection {
	return InitTestWithOptions(t, TestOptions{})
}

// InitTestWithName creates a connection to the database with an established, but otherwise empty
// database _of the given name_.
func InitTestWithName(t *testing.T, dbname string) *database.Connection {
	return InitTestWithOptions(t, TestOptions{DatabaseName: &dbname})
}

// InitTestWithOptions creates a connection to the database with an established, but otherwise empty
// database, configured with.
func InitTestWithOptions(t *testing.T, options TestOptions) *database.Connection {
	(&options).useDefaults()

	if logging.GetSystemLogger() == nil {
		logging.SetupStdoutLogging()
	}

	return database.NewTestConnectionFromNonStandardMigrationPath(t, *options.DatabaseName, *options.DatabasePath)
}

// ClearDB empties the database of all values. This leaves behind small residue: IDs are already taken,
// so, auto-incremented values will use the next value, not re-use values. However, this is easily
// overcome by specifying what the ID should be -- which is part of each seed anyway.
//
// Note: this should only be done in a testing environment.
func ClearDB(db *database.Connection) error {
	systemLogger := logging.GetSystemLogger()
	systemLogger.Log("msg", "Clearing Database...")
	logging.SetSystemLogger(logging.NewNopLogger())
	defer logging.SetSystemLogger(systemLogger)

	err := db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.Delete(sq.Delete("sessions"))
		tx.Delete(sq.Delete("user_operation_permissions"))
		tx.Delete(sq.Delete("user_group_operation_permissions"))
		tx.Delete(sq.Delete("user_operation_preferences"))
		tx.Delete(sq.Delete("api_keys"))
		tx.Delete(sq.Delete("auth_scheme_data"))
		tx.Delete(sq.Delete("email_queue"))
		tx.Delete(sq.Delete("tag_evidence_map"))
		tx.Delete(sq.Delete("tags"))
		tx.Delete(sq.Delete("default_tags"))
		tx.Delete(sq.Delete("evidence_finding_map"))
		tx.Delete(sq.Delete("evidence_metadata"))
		tx.Delete(sq.Delete("evidence"))
		tx.Delete(sq.Delete("findings"))
		tx.Delete(sq.Delete("finding_categories"))
		tx.Delete(sq.Delete("group_user_map"))
		tx.Delete(sq.Delete("users"))
		tx.Delete(sq.Delete("user_groups"))
		tx.Delete(sq.Delete("queries"))
		tx.Delete(sq.Delete("var_operation_map"))
		tx.Delete(sq.Delete("operations"))
		tx.Delete(sq.Delete("service_workers"))
		tx.Delete(sq.Delete("global_vars"))
		tx.Delete(sq.Delete("operation_vars"))
		tx.Delete(sq.Delete("var_operation_map"))
	})
	return err
}

// SimpleFullContext returns back a context with a proper authenticated policy
func SimpleFullContext(my models.User) context.Context {
	ctx := context.Background()
	p := policy.NewAuthenticatedPolicy(my.ID, my.Admin)

	return middleware.InjectIntoContext(ctx, middleware.InjectIntoContextInput{
		UserID:       p.UserID,
		IsSuperAdmin: p.IsSuperAdmin,
		UserPolicy:   p,
	})
}

// CreatePopulatedMemStore generates an in-memory content store with all evidence of the given seed
// populated with tiny versions (see TinyTermRec, TinyImg, and TinyCodeblock). Useful for delete
// tests
func CreatePopulatedMemStore(seed Seeder) *contentstore.MemStore {
	store, _ := contentstore.NewMemStore()
	upload := func(uuid string, data []byte) error { return store.UploadWithName(uuid, bytes.NewReader(data)) }
	for _, evi := range seed.Evidences {
		switch evi.ContentType {
		case "codeblock":
			upload(evi.UUID, TinyCodeblock)
		case "image":
			upload(evi.UUID, TinyImg)
		case "terminal-recording":
			upload(evi.UUID, TinyTermRec)
		case "http-request-cycle":
			// TODO
		}
	}
	return store
}

func MkDBRowCounter(t *testing.T, db *database.Connection, tablename, where string, values ...interface{}) func() int64 {
	return func() int64 {
		return CountRows(t, db, tablename, where, values)
	}
}

func CountRows(t *testing.T, db *database.Connection, tablename, where string, values ...interface{}) int64 {
	var dbQueryCount int64 = -1 // preinitializing to a value we can't get via the query

	err := db.Get(&dbQueryCount, sq.Select("count(*)").
		From(tablename).
		Where(where, values...))
	require.NoError(t, err)
	return dbQueryCount
}

// db queries
func GetAPIKeysForUserID(t *testing.T, db *database.Connection, userID int64) []models.APIKey {
	var apiKeys []models.APIKey
	err := db.Select(&apiKeys, sq.Select("*").
		From("api_keys").
		Where(sq.Eq{"user_id": userID}))
	require.NoError(t, err)
	return apiKeys
}

func GetEvidenceIDsFromFinding(t *testing.T, db *database.Connection, findingID int64) []int64 {
	var list []int64
	err := db.Select(&list, sq.Select("evidence_id").
		From("evidence_finding_map").
		Where(sq.Eq{"finding_id": findingID}).
		OrderBy("evidence_id ASC"))
	require.NoError(t, err)
	return list
}

func GetTagIDsFromEvidenceID(t *testing.T, db *database.Connection, evidenceID int64) []int64 {
	var tagIDs []int64
	err := db.Select(&tagIDs, sq.Select("tag_id").
		From("tag_evidence_map").
		Where(sq.Eq{"evidence_id": evidenceID}))
	require.NoError(t, err)
	return tagIDs
}

func GetFindingByUUID(t *testing.T, db *database.Connection, uuid string) models.Finding {
	var fullFinding models.Finding
	err := db.Get(&fullFinding, sq.Select("*").
		From("findings").
		Where(sq.Eq{"uuid": uuid}))
	require.NoError(t, err)
	return fullFinding
}

func GetEvidenceForOperation(t *testing.T, db *database.Connection, operationID int64) []models.Evidence {
	var evidence []models.Evidence
	err := db.Select(&evidence, sq.Select("*").From("evidence").Where(sq.Eq{"operation_id": operationID}))
	require.NoError(t, err)
	return evidence
}

func GetEvidenceByID(t *testing.T, db *database.Connection, id int64) models.Evidence {
	return GetFullEvidenceViaSelectBuilder(t, db, sq.Eq{"id": id})
}

func GetEvidenceByUUID(t *testing.T, db *database.Connection, uuid string) models.Evidence {
	return GetFullEvidenceViaSelectBuilder(t, db, sq.Eq{"uuid": uuid})
}

func GetEvidenceMetadataByEvidenceID(t *testing.T, db *database.Connection, id int64) []models.EvidenceMetadata {
	var evidenceMetadata []models.EvidenceMetadata
	err := db.Select(&evidenceMetadata, sq.Select("*").
		From("evidence_metadata").
		Where(sq.Eq{"evidence_id": id}))
	require.NoError(t, err)
	return evidenceMetadata
}

func GetFullEvidenceViaSelectBuilder(t *testing.T, db *database.Connection, condition sq.Eq) models.Evidence {
	var evidence models.Evidence
	err := db.Get(&evidence, sq.Select("*").
		From("evidence").
		Where(condition))
	require.NoError(t, err)
	return evidence
}

func GetOperationFromSlug(t *testing.T, db *database.Connection, slug string) models.Operation {
	var fullOp models.Operation
	err := db.Get(&fullOp, sq.Select("id", "slug", "name").
		From("operations").
		Where(sq.Eq{"slug": slug}))
	require.NoError(t, err)
	return fullOp
}

func GetOperations(t *testing.T, db *database.Connection) []services.OperationWithID {
	var operations []struct {
		models.Operation
		NumUsers    int `db:"num_users"`
		NumEvidence int `db:"num_evidence"`
		NumTags     int `db:"num_tags"`
	}

	var topContribs []services.TopContribWithID

	var evidenceCount []services.EvidenceCountWithID

	err := db.Select(&operations, sq.Select("operations.id", "slug", "operations.name", "count(distinct(user_operation_permissions.user_id)) AS num_users", "count(distinct(evidence.id)) AS num_evidence", "count(distinct(tags.id)) AS num_tags").
		From("operations").
		LeftJoin("user_operation_permissions ON user_operation_permissions.operation_id = operations.id").
		LeftJoin("evidence ON evidence.operation_id = operations.id").
		LeftJoin("tags ON tags.operation_id = operations.id").
		GroupBy("operations.id").
		OrderBy("operations.created_at DESC"))

	require.NoError(t, err)

	err = db.SelectRaw(&topContribs, services.GetTopContributorsForEachOperation)

	require.NoError(t, err)

	err = db.SelectRaw(&evidenceCount, services.EvidenceCountForAllOperations)

	require.NoError(t, err)

	operationsWithID := []services.OperationWithID{}
	for _, operation := range operations {

		filteredTopContribs := helpers.Filter(topContribs, func(contributor services.TopContribWithID) bool {
			return contributor.OperationID == operation.ID
		})

		topContribsForOp := make([]dtos.TopContrib, 0, len(filteredTopContribs))
		for i := range filteredTopContribs {
			var topContrib dtos.TopContrib
			topContrib.Slug = filteredTopContribs[i].Slug
			topContrib.Count = filteredTopContribs[i].Count
			topContribsForOp = append(topContribsForOp, topContrib)
		}

		var evidenceCountForOp dtos.EvidenceCount
		idx, _ := helpers.Find(evidenceCount, func(item services.EvidenceCountWithID) bool {
			return item.OperationID == operation.ID
		})
		if idx > -1 {
			evidenceCountForOp.CodeblockCount = evidenceCount[idx].CodeblockCount
			evidenceCountForOp.ImageCount = evidenceCount[idx].ImageCount
			evidenceCountForOp.HarCount = evidenceCount[idx].HarCount
			evidenceCountForOp.EventCount = evidenceCount[idx].EventCount
			evidenceCountForOp.RecordingCount = evidenceCount[idx].RecordingCount
		}

		operationsWithID = append(operationsWithID, services.OperationWithID{
			ID: operation.ID,
			Op: &dtos.Operation{
				Slug:          operation.Slug,
				Name:          operation.Name,
				NumUsers:      operation.NumUsers,
				NumEvidence:   operation.NumEvidence,
				NumTags:       operation.NumTags,
				TopContribs:   topContribsForOp,
				EvidenceCount: evidenceCountForOp,
			},
		})
	}
	return operationsWithID
}

// TODO this is too close to what ListOperations does. We can leave it like this, but should have a
// a test to verify the intial case (matches the seed data), plus a minor addition
func GetOperationsForUser(t *testing.T, db *database.Connection, user models.User) []*dtos.Operation {
	operations := GetOperations(t, db)
	ctx := ContextForUser(user, db)

	var operationPreference []models.UserOperationPreferences
	err := db.Select(&operationPreference, sq.Select("operation_id", "is_favorite").
		From("user_operation_preferences").
		Where(sq.Eq{"user_id": middleware.UserID(ctx)}))

	require.NoError(t, err)

	operationPreferenceMap := make(map[int64]bool)
	for _, op := range operationPreference {
		operationPreferenceMap[op.OperationID] = op.IsFavorite
	}

	filteredOperationsDTO := make([]*dtos.Operation, 0, len(operations))
	for _, operation := range operations {
		if middleware.Policy(ctx).Check(policy.CanReadOperation{OperationID: operation.ID}) {
			operation.Op.Favorite = operationPreferenceMap[operation.ID]
			filteredOperationsDTO = append(filteredOperationsDTO, operation.Op)
		}
	}
	return filteredOperationsDTO
}

func GetUserRolesForOperationByOperationID(t *testing.T, db *database.Connection, id int64) []models.UserOperationPermission {
	var userRoles []models.UserOperationPermission
	err := db.Select(&userRoles, sq.Select("*").
		From("user_operation_permissions").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return userRoles
}

func GetQueryByID(t *testing.T, db *database.Connection, id int64) models.Query {
	var fullQuery models.Query
	err := db.Get(&fullQuery, sq.Select("*").
		From("queries").
		Where(sq.Eq{"id": id}))
	require.NoError(t, err)
	return fullQuery
}

func GetQueriesForOperationID(t *testing.T, db *database.Connection, id int64) []models.Query {
	var allQueries []models.Query
	err := db.Select(&allQueries, sq.Select("*").
		From("queries").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return allQueries
}

func GetTagByID(t *testing.T, db *database.Connection, id int64) models.Tag {
	var tag models.Tag
	err := db.Get(&tag, sq.Select("*").
		From("tags").
		Where(sq.Eq{"id": id}))
	require.NoError(t, err)
	return tag
}

func GetDefaultTagByID(t *testing.T, db *database.Connection, id int64) models.DefaultTag {
	var tag models.DefaultTag
	err := db.Get(&tag, sq.Select("*").
		From("default_tags").
		Where(sq.Eq{"id": id}))
	require.NoError(t, err)
	return tag
}

func GetTagFromOperationID(t *testing.T, db *database.Connection, id int64) []models.Tag {
	var allTags []models.Tag
	err := db.Select(&allTags, sq.Select("*").
		From("tags").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return allTags
}

func GetDefaultTags(t *testing.T, db *database.Connection) []models.DefaultTag {
	var allTags []models.DefaultTag
	err := db.Select(&allTags, sq.Select("*").
		From("default_tags"))
	require.NoError(t, err)
	return allTags
}

func GetFindingsByOperationID(t *testing.T, db *database.Connection, id int64) []models.Finding {
	var findings []models.Finding
	err := db.Select(&findings, sq.Select("*").
		From("findings").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return findings
}

func GetUserProfile(t *testing.T, db *database.Connection, id int64) models.User {
	var user models.User
	err := db.Get(&user, sq.Select("id", "slug", "first_name", "last_name", "email", "admin", "disabled").
		From("users").
		Where(sq.Eq{"id": id}))
	require.NoError(t, err)
	return user
}

func GetUserBySlug(t *testing.T, db *database.Connection, slug string) models.User {
	user, err := db.RetrieveUserBySlug(slug)
	require.NoError(t, err)
	return user
}

func GetAllUsers(t *testing.T, db *database.Connection) []models.User {
	var users []models.User
	err := db.Select(&users, sq.Select("*").From("users"))
	require.NoError(t, err)
	return users
}

func GetAllDeletedUsers(t *testing.T, db *database.Connection) []models.User {
	var users []models.User
	err := db.Select(&users, sq.Select("*").From("users").Where(sq.NotEq{"deleted_at": nil}))
	require.Nil(t, err)
	return users
}

func GetAuthsForUser(t *testing.T, db *database.Connection, userID int64) []models.AuthSchemeData {
	var schemes []models.AuthSchemeData
	err := db.Select(&schemes, sq.Select("*").From("auth_scheme_data").
		Where(sq.Eq{"user_id": userID}))
	require.Nil(t, err)
	return schemes
}

func GetUsersForAuth(t *testing.T, db *database.Connection, authName string) []models.User {
	// return a list of users that: 1. aren't deleted 2. aren't headless 3. have the given auth scheme
	var users []models.User
	err := db.Select(&users, sq.Select("distinctrow users.*").From("users").
		Join("auth_scheme_data ON user_id = users.id").Where(sq.Eq{"users.deleted_at": nil, "auth_scheme": authName}))
	require.Nil(t, err)
	return users
}

func GetRealUsers(t *testing.T, db *database.Connection) []models.User {
	// return a list of users that: 1. aren't deleted 2. aren't headless
	var users []models.User
	err := db.Select(&users, sq.Select("distinctrow users.*").From("users").
		Join("auth_scheme_data ON user_id = users.id").Where(sq.Eq{"users.deleted_at": nil}))
	require.Nil(t, err)
	return users
}

func GetTagUsage(t *testing.T, db *database.Connection, tagID int64) int64 {
	var usageCount int64
	db.Get(&usageCount, sq.Select("count(*)").
		From("tag_evidence_map").
		Where(sq.Eq{"tag_id": tagID}))
	// ignoring the error here -- just return 0, which is appropriate anyway
	return usageCount
}

type FullEvidence struct {
	models.Evidence
	// copied from models.User
	Slug      string `db:"slug"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Tags      []models.Tag
}

func GetFullEvidenceByFindingID(t *testing.T, db *database.Connection, findingID int64) []FullEvidence {
	var allFullEvidence []FullEvidence
	err := db.Select(&allFullEvidence, sq.Select("evidence.*", "users.first_name", "users.last_name", "users.slug").
		From("evidence_finding_map").
		LeftJoin("evidence ON evidence_finding_map.evidence_id = evidence.id").
		LeftJoin("users on evidence.operator_id = users.id").
		Where(sq.Eq{"finding_id": findingID}))
	require.NoError(t, err)
	FillEvidenceWithTags(t, db, &allFullEvidence)

	return allFullEvidence
}

func GetServiceWorkerByName(t *testing.T, db *database.Connection, name string) models.ServiceWorker {
	var worker models.ServiceWorker
	err := db.Get(&worker, sq.Select("*").From("service_workers").Where(sq.Eq{"name": name}))
	require.NoError(t, err)
	return worker
}

func GetServiceWorkerByID(t *testing.T, db *database.Connection, id int64) models.ServiceWorker {
	var worker models.ServiceWorker
	err := db.Get(&worker, sq.Select("*").From("service_workers").Where(sq.Eq{"id": id}))
	require.NoError(t, err)
	return worker
}

func ListServiceWorkers(t *testing.T, db *database.Connection) []models.ServiceWorker {
	var workers []models.ServiceWorker
	err := db.Select(&workers, sq.Select("*").From("service_workers"))
	require.NoError(t, err)
	return workers
}

func GetFullEvidenceByOperationID(t *testing.T, db *database.Connection, operationID int64) []FullEvidence {
	var allFullEvidence []FullEvidence
	err := db.Select(&allFullEvidence, sq.Select("evidence.*", "users.first_name", "users.last_name", "users.slug").
		From("evidence").
		LeftJoin("users on evidence.operator_id = users.id").
		Where(sq.Eq{"operation_id": operationID}))
	require.NoError(t, err)

	FillEvidenceWithTags(t, db, &allFullEvidence)

	return allFullEvidence
}

func FillEvidenceWithTags(t *testing.T, db *database.Connection, evidence *[]FullEvidence) {
	for i, item := range *evidence {
		var tags []models.Tag
		err := db.Select(&tags, sq.Select("tags.*").
			From("tag_evidence_map").
			LeftJoin("tags ON tags.id = tag_evidence_map.tag_id").
			Where(sq.Eq{"evidence_id": item.Evidence.ID}))
		require.NoError(t, err)
		(*evidence)[i].Tags = tags
	}
}

type UserOpPermJoinUser struct {
	models.User
	Role policy.OperationRole `db:"role"`
}

func GetUsersWithRoleForOperationByOperationID(t *testing.T, db *database.Connection, id int64) []UserOpPermJoinUser {
	var allUserOpRoles []UserOpPermJoinUser
	err := db.Select(&allUserOpRoles, sq.Select("user_operation_permissions.role", "users.first_name", "users.last_name", "users.slug").
		From("user_operation_permissions").
		LeftJoin("users ON users.id = user_operation_permissions.user_id").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return allUserOpRoles
}

type UserGroupOpPermJoinUser struct {
	models.UserGroup
	Role policy.OperationRole `db:"role"`
}

func GetUserGroupsWithRoleForOperationByOperationID(t *testing.T, db *database.Connection, id int64) []UserGroupOpPermJoinUser {
	var allUserGroupOpRoles []UserGroupOpPermJoinUser
	err := db.Select(&allUserGroupOpRoles, sq.Select("user_group_operation_permissions.role", "user_groups.name", "user_groups.slug").
		From("user_group_operation_permissions").
		LeftJoin("user_groups ON user_groups.id = user_group_operation_permissions.group_id").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return allUserGroupOpRoles
}

func GetUserGroupFromSlug(t *testing.T, db *database.Connection, slug string) models.UserGroup {
	var fullUserGroup models.UserGroup
	err := db.Get(&fullUserGroup, sq.Select("id", "slug", "name").
		From("user_groups").
		Where(sq.Eq{"slug": slug}))
	require.NoError(t, err)
	return fullUserGroup
}

type PreferencesOperations struct {
	models.UserOperationPreferences
	Slug string `db:"slug"`
}

func GetFavoritesByUserID(t *testing.T, db *database.Connection, id int64) []PreferencesOperations {
	var preferencesOperations []PreferencesOperations
	err := db.Select(&preferencesOperations, sq.Select("user_operation_preferences.*", "operations.slug").
		From("user_operation_preferences").
		Join("operations on user_operation_preferences.operation_id = operations.id").
		Where(sq.Eq{"user_id": id, "is_favorite": true}))
	require.NoError(t, err)
	return preferencesOperations
}

func GetFavoriteForOperation(t *testing.T, db *database.Connection, slug string, id int64) bool {
	var favorites []bool
	err := db.Select(&favorites, sq.Select("is_favorite").
		From("user_operation_preferences").
		Join("operations on user_operation_preferences.operation_id = operations.id").
		Where(sq.Eq{"slug": slug, "user_id": id}))
	require.NoError(t, err)

	isFavorite := len(favorites) > 0 && favorites[0]

	return isFavorite
}

func GetGlobalVarFromName(t *testing.T, db *database.Connection, name string) models.GlobalVar {
	var globalVar models.GlobalVar

	err := db.Get(&globalVar, sq.Select("*").
		From("global_vars").
		Where(sq.Eq{"name": name}))
	require.NoError(t, err)
	return globalVar
}

func GetOperationVarFromSlug(t *testing.T, db *database.Connection, slug string) models.OperationVar {
	var operationVar models.OperationVar

	err := db.Get(&operationVar, sq.Select("*").
		From("operation_vars").
		Where(sq.Eq{"slug": slug}))
	require.NoError(t, err)
	return operationVar
}

type TestOptions struct {
	DatabasePath *string
	DatabaseName *string
}

func (opts *TestOptions) useDefaults() {
	if opts.DatabasePath == nil {
		opts.DatabasePath = helpers.Ptr("../migrations")
	}
	if opts.DatabaseName == nil {
		opts.DatabaseName = helpers.Ptr("service-test-db")
	}
}

func ApplySeeding(t *testing.T, seed Seeder, db *database.Connection) {
	err := seed.ApplyTo(db)
	require.NoError(t, err)
}
