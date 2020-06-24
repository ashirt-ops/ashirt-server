// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"bytes"
	"context"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt/backend/contentstore"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/logging"
	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
	"github.com/theparanoids/ashirt/backend/server/middleware"
)

var internalClock = clockwork.NewFakeClock()

var TinyImg []byte = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
	0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
	0x42, 0x60, 0x82,
} // tiniest png https://github.com/mathiasbynens/small

var TinyCodeblock []byte = []byte(`{"contentType": "codeblock", "contentSubtype": "python", "content": "print(\"Hello World!\")"}`)

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

func simpleFullContext(my models.User) context.Context {
	ctx := context.Background()
	p := policy.NewAuthenticatedPolicy(my.ID, my.Admin)

	return middleware.InjectIntoContext(ctx, middleware.InjectIntoContextInput{
		UserID:       p.UserID,
		IsSuperAdmin: p.IsSuperAdmin,
		UserPolicy:   p,
	})
}

func fullContext(userid int64, p policy.Policy) context.Context {
	ctx := context.Background()
	return middleware.InjectIntoContext(ctx, middleware.InjectIntoContextInput{
		UserID:       userid,
		IsSuperAdmin: false,
		UserPolicy:   p,
	})
}

func fullContextAsAdmin(userid int64, p policy.Policy) context.Context {
	ctx := context.Background()
	return middleware.InjectIntoContext(ctx, middleware.InjectIntoContextInput{
		UserID:       userid,
		IsSuperAdmin: true,
		UserPolicy:   p,
	})
}

// initTest creates a connection to the database and provides an established, but otherwise empty
// database.
func initTest(t *testing.T) *database.Connection {
	logging.SetupStdoutLogging()
	return database.NewTestConnection(t, "service-test-db")
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

func createPopulatedMemStore(seed TestSeedData) *contentstore.MemStore {
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
		}
	}
	return store
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

func newEvidenceGen(first int64) func(opID, ownerID int64, uuid, desc, contentType string) models.Evidence {
	id := iotaLike(first)
	return func(opID, ownerID int64, uuid, desc, contentType string) models.Evidence {
		return models.Evidence{
			ID:            id(),
			UUID:          uuid,
			OperationID:   opID,
			OperatorID:    ownerID,
			Description:   desc,
			ContentType:   contentType,
			FullImageKey:  uuid,
			ThumbImageKey: uuid,
			OccurredAt:    internalClock.Now(),
			CreatedAt:     internalClock.Now(),
		}
	}
}

func newFindingGen(first int64) func(opID int64, uuid, category, title, desc string, ticketLink *string) models.Finding {
	id := iotaLike(first)
	return func(opID int64, uuid, category, title, desc string, ticketLink *string) models.Finding {
		finding := models.Finding{
			ID:            id(),
			OperationID:   opID,
			UUID:          uuid,
			Category:      category,
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

func makeDBRowCounter(t *testing.T, db *database.Connection, tablename, where string, values ...interface{}) func() int64 {
	return func() int64 {
		return countRows(t, db, tablename, where, values)
	}
}

func countRows(t *testing.T, db *database.Connection, tablename, where string, values ...interface{}) int64 {
	var dbQueryCount int64 = -1 // preinitializing to a value we can't get via the query

	err := db.Get(&dbQueryCount, sq.Select("count(*)").
		From(tablename).
		Where(where, values...))
	require.NoError(t, err)
	return dbQueryCount
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

// sorted orders an int slice in asc order, then returns back a copy of the sorted list
// note: underlying list is
func sorted(slice []int64) []int64 {
	clone := make([]int64, len(slice))
	copy(clone, slice)
	sort.Slice(clone, func(i, j int) bool { return clone[i] < clone[j] })
	return clone
}

// db queries
func getAPIKeysForUserID(t *testing.T, db *database.Connection, userID int64) []models.APIKey {
	var apiKeys []models.APIKey
	err := db.Select(&apiKeys, sq.Select("*").
		From("api_keys").
		Where(sq.Eq{"user_id": userID}))
	require.NoError(t, err)
	return apiKeys
}

func getEvidenceIDsFromFinding(t *testing.T, db *database.Connection, findingID int64) []int64 {
	var list []int64
	err := db.Select(&list, sq.Select("evidence_id").
		From("evidence_finding_map").
		Where(sq.Eq{"finding_id": findingID}).
		OrderBy("evidence_id ASC"))
	require.NoError(t, err)
	return list
}

func getTagIDsFromEvidenceID(t *testing.T, db *database.Connection, evidenceID int64) []int64 {
	var tagIDs []int64
	err := db.Select(&tagIDs, sq.Select("tag_id").
		From("tag_evidence_map").
		Where(sq.Eq{"evidence_id": evidenceID}))
	require.NoError(t, err)
	return tagIDs
}

func getFindingByUUID(t *testing.T, db *database.Connection, uuid string) models.Finding {
	var fullFinding models.Finding
	err := db.Get(&fullFinding, sq.Select("*").
		From("findings").
		Where(sq.Eq{"uuid": uuid}))
	require.NoError(t, err)
	return fullFinding
}

func getEvidenceForOperation(t *testing.T, db *database.Connection, operationID int64) []models.Evidence{
	var evidence []models.Evidence
	err := db.Select(&evidence, sq.Select("*").From("evidence").Where(sq.Eq{"operation_id": operationID}))
	require.NoError(t, err)
	return evidence
}

func getEvidenceByID(t *testing.T, db *database.Connection, id int64) models.Evidence {
	return getFullEvidenceViaSelectBuilder(t, db, sq.Eq{"id": id})
}

func getEvidenceByUUID(t *testing.T, db *database.Connection, uuid string) models.Evidence {
	return getFullEvidenceViaSelectBuilder(t, db, sq.Eq{"uuid": uuid})
}

func getFullEvidenceViaSelectBuilder(t *testing.T, db *database.Connection, condition sq.Eq) models.Evidence {
	var evidence models.Evidence
	err := db.Get(&evidence, sq.Select("*").
		From("evidence").
		Where(condition))
	require.NoError(t, err)
	return evidence
}

func getOperationFromSlug(t *testing.T, db *database.Connection, slug string) models.Operation {
	var fullOp models.Operation
	err := db.Get(&fullOp, sq.Select("id", "slug", "name", "status").
		From("operations").
		Where(sq.Eq{"slug": slug}))
	require.NoError(t, err)
	return fullOp
}

func getOperations(t *testing.T, db *database.Connection) []models.Operation {
	var fullOps []models.Operation
	err := db.Select(&fullOps, sq.Select("id", "slug", "name", "status").
		From("operations"))
	require.NoError(t, err)
	return fullOps
}

func getUserRolesForOperationByOperationID(t *testing.T, db *database.Connection, id int64) []models.UserOperationPermission {
	var userRoles []models.UserOperationPermission
	err := db.Select(&userRoles, sq.Select("*").
		From("user_operation_permissions").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return userRoles
}

func getQueryByID(t *testing.T, db *database.Connection, id int64) models.Query {
	var fullQuery models.Query
	err := db.Get(&fullQuery, sq.Select("*").
		From("queries").
		Where(sq.Eq{"id": id}))
	require.NoError(t, err)
	return fullQuery
}

func getQueriesForOperationID(t *testing.T, db *database.Connection, id int64) []models.Query {
	var allQueries []models.Query
	err := db.Select(&allQueries, sq.Select("*").
		From("queries").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return allQueries
}

func getTagByID(t *testing.T, db *database.Connection, id int64) models.Tag {
	var tag models.Tag
	err := db.Get(&tag, sq.Select("*").
		From("tags").
		Where(sq.Eq{"id": id}))
	require.NoError(t, err)
	return tag
}

func getTagFromOperationID(t *testing.T, db *database.Connection, id int64) []models.Tag {
	var allTags []models.Tag
	err := db.Select(&allTags, sq.Select("*").
		From("tags").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return allTags
}

func getFindingsByOperationID(t *testing.T, db *database.Connection, id int64) []models.Finding {
	var findings []models.Finding
	err := db.Select(&findings, sq.Select("*").
		From("findings").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return findings
}

func getUserProfile(t *testing.T, db *database.Connection, id int64) models.User {
	var user models.User
	err := db.Get(&user, sq.Select("id", "slug", "first_name", "last_name", "email", "admin", "disabled").
		From("users").
		Where(sq.Eq{"id": id}))
	require.NoError(t, err)
	return user
}

func getUserBySlug(t *testing.T, db *database.Connection, slug string) models.User {
	user, err := db.RetrieveUserBySlug(slug)
	require.NoError(t, err)
	return user
}

func getAllUsers(t *testing.T, db *database.Connection) []models.User {
	var users []models.User
	err := db.Select(&users, sq.Select("*").From("users"))
	require.NoError(t, err)
	return users
}

func getAllDeletedUsers(t *testing.T, db *database.Connection) []models.User {
	var users []models.User
	err := db.Select(&users, sq.Select("*").From("users").Where(sq.NotEq{"deleted_at": nil}))
	require.Nil(t, err)
	return users
}

func getAuthsForUser(t *testing.T, db *database.Connection, userID int64) []models.AuthSchemeData {
	var schemes []models.AuthSchemeData
	err := db.Select(&schemes, sq.Select("*").From("auth_scheme_data").
		Where(sq.Eq{"user_id": userID}))
	require.Nil(t, err)
	return schemes
}

func getUsersForAuth(t *testing.T, db *database.Connection, authName string) []models.User {
	// return a list of users that: 1. aren't deleted 2. aren't headless 3. have the given auth scheme
	var users []models.User
	err := db.Select(&users, sq.Select("distinctrow users.*").From("users").
		Join("auth_scheme_data ON user_id = users.id").Where(sq.Eq{"users.deleted_at": nil, "auth_scheme": authName}))
	require.Nil(t, err)
	return users
}

func getRealUsers(t *testing.T, db *database.Connection) []models.User {
	// return a list of users that: 1. aren't deleted 2. aren't headless
	var users []models.User
	err := db.Select(&users, sq.Select("distinctrow users.*").From("users").
		Join("auth_scheme_data ON user_id = users.id").Where(sq.Eq{"users.deleted_at": nil}))
	require.Nil(t, err)
	return users
}

type FullEvidence struct {
	models.Evidence
	// copied from models.User
	Slug      string `db:"slug"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Tags      []models.Tag
}

func getFullEvidenceByFindingID(t *testing.T, db *database.Connection, findingID int64) []FullEvidence {
	var allFullEvidence []FullEvidence
	err := db.Select(&allFullEvidence, sq.Select("evidence.*", "users.first_name", "users.last_name", "users.slug").
		From("evidence_finding_map").
		LeftJoin("evidence ON evidence_finding_map.evidence_id = evidence.id").
		LeftJoin("users on evidence.operator_id = users.id").
		Where(sq.Eq{"finding_id": findingID}))
	require.NoError(t, err)
	fillEvidenceWithTags(t, db, &allFullEvidence)

	return allFullEvidence
}

func getFullEvidenceByOperationID(t *testing.T, db *database.Connection, operationID int64) []FullEvidence {
	var allFullEvidence []FullEvidence
	err := db.Select(&allFullEvidence, sq.Select("evidence.*", "users.first_name", "users.last_name", "users.slug").
		From("evidence").
		LeftJoin("users on evidence.operator_id = users.id").
		Where(sq.Eq{"operation_id": operationID}))
	require.NoError(t, err)

	fillEvidenceWithTags(t, db, &allFullEvidence)

	return allFullEvidence
}

func fillEvidenceWithTags(t *testing.T, db *database.Connection, evidence *[]FullEvidence) {
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

func getUsersWithRoleForOperationByOperationID(t *testing.T, db *database.Connection, id int64) []UserOpPermJoinUser {
	var allUserOpRoles []UserOpPermJoinUser
	err := db.Select(&allUserOpRoles, sq.Select("user_operation_permissions.role", "users.first_name", "users.last_name", "users.slug").
		From("user_operation_permissions").
		LeftJoin("users ON users.id = user_operation_permissions.user_id").
		Where(sq.Eq{"operation_id": id}))
	require.NoError(t, err)
	return allUserOpRoles
}
