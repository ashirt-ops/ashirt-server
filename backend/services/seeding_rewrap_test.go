// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/contentstore"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/database/seeding"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
)

// This file rewraps many of the database seeder helpers.
// seed backend/database/seeding for the real values

// Exported example data + types
var TinyImg = seeding.TinyImg
var TinyCodeblock = seeding.TinyCodeblock
var TinyTermRec = seeding.TinyTermRec

type UserOpPermJoinUser = seeding.UserOpPermJoinUser
type FullEvidence = seeding.FullEvidence

// Exported functions/helpers
var initTest = seeding.InitTest
var getUsersWithRoleForOperationByOperationID = seeding.GetUsersWithRoleForOperationByOperationID
var contextForUser = seeding.ContextForUser
var GetInternalClock = seeding.GetInternalClock

var getFullEvidenceByFindingID = seeding.GetFullEvidenceByFindingID
var getFullEvidenceByOperationID = seeding.GetFullEvidenceByOperationID
var fillEvidenceWithTags = seeding.FillEvidenceWithTags
var getAPIKeysForUserID = seeding.GetAPIKeysForUserID
var getEvidenceIDsFromFinding = seeding.GetEvidenceIDsFromFinding
var getTagIDsFromEvidenceID = seeding.GetTagIDsFromEvidenceID
var getFindingByUUID = seeding.GetFindingByUUID
var getEvidenceForOperation = seeding.GetEvidenceForOperation
var getEvidenceByID = seeding.GetEvidenceByID
var getEvidenceByUUID = seeding.GetEvidenceByUUID
var getEvidenceMetadataByEvidenceID = seeding.GetEvidenceMetadataByEvidenceID
var getFullEvidenceViaSelectBuilder = seeding.GetFullEvidenceViaSelectBuilder
var getOperationFromSlug = seeding.GetOperationFromSlug
var getOperations = seeding.GetOperations
var getOperationsForUser = seeding.GetOperationsForUser
var getUserRolesForOperationByOperationID = seeding.GetUserRolesForOperationByOperationID
var getQueryByID = seeding.GetQueryByID
var getQueriesForOperationID = seeding.GetQueriesForOperationID
var getTagByID = seeding.GetTagByID
var getDefaultTagByID = seeding.GetDefaultTagByID
var getTagFromOperationID = seeding.GetTagFromOperationID
var getDefaultTags = seeding.GetDefaultTags
var getFindingsByOperationID = seeding.GetFindingsByOperationID
var getUserProfile = seeding.GetUserProfile
var getUserBySlug = seeding.GetUserBySlug
var getAllUsers = seeding.GetAllUsers
var getAllDeletedUsers = seeding.GetAllDeletedUsers
var getAuthsForUser = seeding.GetAuthsForUser
var getUsersForAuth = seeding.GetUsersForAuth
var getRealUsers = seeding.GetRealUsers
var getTagUsage = seeding.GetTagUsage

var getServiceWorkerByName = seeding.GetServiceWorkerByName
var getServiceWorkerByID = seeding.GetServiceWorkerByID
var listServiceWorkers = seeding.ListServiceWorkers

var TagIDsFromTags = seeding.TagIDsFromTags
var sorted = seeding.Sorted
var makeDBRowCounter = seeding.MkDBRowCounter
var countRows = seeding.CountRows

func createPopulatedMemStore(seed TestSeedData) *contentstore.MemStore {
	return seeding.CreatePopulatedMemStore(seed.Seeder)
}

// Exported seed data
var HarryPotterSeedData = TestSeedData{
	seeding.HarryPotterSeedData,
}
var NoSeedData = TestSeedData {}

var UserDumbledore = seeding.UserDumbledore
var UserHarry = seeding.UserHarry
var UserRon = seeding.UserRon
var UserGinny = seeding.UserGinny
var UserHermione = seeding.UserHermione
var UserNeville = seeding.UserNeville
var UserSeamus = seeding.UserSeamus
var UserDraco = seeding.UserDraco
var UserSnape = seeding.UserSnape
var UserHagrid = seeding.UserHagrid
var UserTomRiddle = seeding.UserTomRiddle
var UserHeadlessNick = seeding.UserHeadlessNick
var UserCedric = seeding.UserCedric
var UserFleur = seeding.UserFleur
var UserViktor = seeding.UserViktor
var UserAlastor = seeding.UserAlastor
var UserMinerva = seeding.UserMinerva
var UserLucius = seeding.UserLucius
var UserSirius = seeding.UserSirius
var UserPeter = seeding.UserPeter
var UserParvati = seeding.UserParvati
var UserPadma = seeding.UserPadma
var UserCho = seeding.UserCho

var APIKeyHarry1 = seeding.APIKeyHarry1
var APIKeyHarry2 = seeding.APIKeyHarry2
var APIKeyRon1 = seeding.APIKeyRon1
var APIKeyRon2 = seeding.APIKeyRon2
var APIKeyNick = seeding.APIKeyNick

var OpSorcerersStone = seeding.OpSorcerersStone
var OpChamberOfSecrets = seeding.OpChamberOfSecrets
var OpPrisonerOfAzkaban = seeding.OpPrisonerOfAzkaban
var OpGobletOfFire = seeding.OpGobletOfFire
var OpOrderOfThePhoenix = seeding.OpOrderOfThePhoenix
var OpHalfBloodPrince = seeding.OpHalfBloodPrince
var OpDeathlyHallows = seeding.OpDeathlyHallows
var OpGanttChart = seeding.OpGanttChart

var DefaultTagWho = seeding.DefaultTagWho
var DefaultTagWhat = seeding.DefaultTagWhat
var DefaultTagWhere = seeding.DefaultTagWhere
var DefaultTagWhen = seeding.DefaultTagWhen
var DefaultTagWhy = seeding.DefaultTagWhy

var TagFamily = seeding.TagFamily
var TagFriendship = seeding.TagFriendship
var TagHome = seeding.TagHome
var TagLoyalty = seeding.TagLoyalty
var TagCourage = seeding.TagCourage
var TagGoodVsEvil = seeding.TagGoodVsEvil
var TagSupernatural = seeding.TagSupernatural

var TagMercury = seeding.TagMercury
var TagVenus = seeding.TagVenus
var TagEarth = seeding.TagEarth
var TagMars = seeding.TagMars
var TagJupiter = seeding.TagJupiter
var TagSaturn = seeding.TagSaturn
var TagNeptune = seeding.TagNeptune

var CommonTagWhoSS = seeding.CommonTagWhoSS
var CommonTagWhatSS = seeding.CommonTagWhatSS
var CommonTagWhereSS = seeding.CommonTagWhereSS
var CommonTagWhenSS = seeding.CommonTagWhenSS
var CommonTagWhySS = seeding.CommonTagWhySS

var CommonTagWhoCoS = seeding.CommonTagWhoCoS
var CommonTagWhatCoS = seeding.CommonTagWhatCoS
var CommonTagWhereCoS = seeding.CommonTagWhereCoS
var CommonTagWhenCoS = seeding.CommonTagWhenCoS
var CommonTagWhyCoS = seeding.CommonTagWhyCoS

var CommonTagWhoGantt = seeding.CommonTagWhoGantt
var CommonTagWhatGantt = seeding.CommonTagWhatGantt
var CommonTagWhereGantt = seeding.CommonTagWhereGantt
var CommonTagWhenGantt = seeding.CommonTagWhenGantt
var CommonTagWhyGantt = seeding.CommonTagWhyGantt

var EviDursleys = seeding.EviDursleys
var EviMirrorOfErised = seeding.EviMirrorOfErised

var EviFlyingCar = seeding.EviFlyingCar
var EviDobby = seeding.EviDobby
var EviSpiderAragog = seeding.EviSpiderAragog
var EviMoaningMyrtle = seeding.EviMoaningMyrtle
var EviWhompingWillow = seeding.EviWhompingWillow
var EviTomRiddlesDiary = seeding.EviTomRiddlesDiary

var EviPetrifiedHermione = seeding.EviPetrifiedHermione

var EviGantt01 = seeding.EviGantt01
var EviGantt02 = seeding.EviGantt02
var EviGantt03 = seeding.EviGantt03
var EviGantt04 = seeding.EviGantt04
var EviGantt05 = seeding.EviGantt05
var EviGantt06 = seeding.EviGantt06
var EviGantt07 = seeding.EviGantt07
var EviGantt08 = seeding.EviGantt08
var EviGantt09 = seeding.EviGantt09
var EviGantt10 = seeding.EviGantt10
var EviGantt11 = seeding.EviGantt11
var EviGantt12 = seeding.EviGantt12
var EviGantt13 = seeding.EviGantt13
var EviGantt14 = seeding.EviGantt14
var EviGantt15 = seeding.EviGantt15
var EviGantt16 = seeding.EviGantt16
var EviGantt17 = seeding.EviGantt17
var EviGantt18 = seeding.EviGantt18
var EviGantt19 = seeding.EviGantt19
var EviGantt20 = seeding.EviGantt20
var EviGanttExtra = seeding.EviGanttExtra

var QuerySalazarsHier = seeding.QuerySalazarsHier
var QueryWhereIsTheChamberOfSecrets = seeding.QueryWhereIsTheChamberOfSecrets

var FindingBook2Magic = seeding.FindingBook2Magic
var FindingBook2CGI = seeding.FindingBook2CGI
var FindingBook2SpiderFear = seeding.FindingBook2SpiderFear

var ProductFindingCategory = seeding.ProductFindingCategory
var NetworkFindingCategory = seeding.NetworkFindingCategory
var EnterpriseFindingCategory = seeding.EnterpriseFindingCategory
var VendorFindingCategory = seeding.VendorFindingCategory
var BehavioralFindingCategory = seeding.BehavioralFindingCategory
var DetectionGapFindingCategory = seeding.DetectionGapFindingCategory
var DeletedCategory = seeding.DeletedCategory

var DemoServiceWorker = seeding.DemoServiceWorker

type TestSeedData struct {
	seeding.Seeder
}

func (seed TestSeedData) ApplyTo(t *testing.T, db *database.Connection) {
	err := seed.Seeder.ApplyTo(db)
	require.NoError(t, err)
}

func (seed TestSeedData) AllInitialTagIds() []int64 {
	return seed.Seeder.AllInitialTagIds()
}

func (seed TestSeedData) EvidenceIDsForFinding(finding models.Finding) []int64 {
	return seed.Seeder.EvidenceIDsForFinding(finding)
}

func (seed TestSeedData) TagsForFinding(finding models.Finding) []models.Tag {
	return seed.Seeder.TagsForFinding(finding)
}

func (seed TestSeedData) TagsForEvidence(evidence models.Evidence) []models.Tag {
	return seed.Seeder.TagsForEvidence(evidence)
}

func (seed TestSeedData) GetTagFromID(id int64) models.Tag {
	return seed.Seeder.GetTagFromID(id)
}

func (seed TestSeedData) GetUserFromID(id int64) models.User {
	return seed.Seeder.GetUserFromID(id)
}

func (seed TestSeedData) UsersForOp(op models.Operation) []models.User {
	return seed.Seeder.UsersForOp(op)
}

func (seed TestSeedData) UserRoleForOp(user models.User, op models.Operation) policy.OperationRole {
	return seed.Seeder.UserRoleForOp(user, op)
}

func (seed TestSeedData) EvidenceForOperation(opID int64) []models.Evidence {
	return seed.Seeder.EvidenceForOperation(opID)
}

func (seed TestSeedData) TagIDsUsageByDate(opID int64) map[int64][]time.Time {
	return seed.Seeder.TagIDsUsageByDate(opID)
}

func (seed TestSeedData) CategoryForFinding(finding models.Finding) string {
	return seed.Seeder.CategoryForFinding(finding)
}
