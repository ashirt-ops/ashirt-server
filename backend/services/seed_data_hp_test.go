// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"strings"

	"github.com/theparanoids/ashirt/backend/models"
	"github.com/theparanoids/ashirt/backend/policy"
)

var HarryPotterSeedData = TestSeedData{
	Users:      []models.User{UserHarry, UserRon, UserGinny, UserHermione, UserNeville, UserDraco, UserSnape, UserDumbledore, UserHagrid, UserTomRiddle, UserHeadlessNick},
	Operations: []models.Operation{OpSorcerersStone, OpChamberOfSecrets, OpPrisonerOfAzkaban, OpGobletOfFire, OpOrderOfThePhoenix, OpHalfBloodPrince, OpDeathlyHallows},
	Tags: []models.Tag{
		TagFamily, TagFriendship, TagHome, TagLoyalty, TagCourage, TagGoodVsEvil, TagSupernatural,
		TagMercury, TagVenus, TagEarth, TagMars, TagJupiter, TagSaturn, TagNeptune,
	},
	APIKeys: []models.APIKey{
		APIKeyHarry1, APIKeyHarry2,
		APIKeyRon1, APIKeyRon2,
	},
	UserOpMap: []models.UserOperationPermission{
		newUserOpPermission(UserHarry, OpSorcerersStone, policy.OperationRoleAdmin),
		newUserOpPermission(UserHarry, OpChamberOfSecrets, policy.OperationRoleWrite),
		newUserOpPermission(UserRon, OpChamberOfSecrets, policy.OperationRoleAdmin),

		newUserOpPermission(UserDumbledore, OpSorcerersStone, policy.OperationRoleAdmin),
		newUserOpPermission(UserDumbledore, OpChamberOfSecrets, policy.OperationRoleAdmin),
	},
	Findings: []models.Finding{
		FindingBook2Magic, FindingBook2CGI, FindingBook2SpiderFear,
	},
	Evidences: []models.Evidence{
		EviDursleys, EviMirrorOfErised,
		EviFlyingCar, EviDobby, EviSpiderAragog, EviMoaningMyrtle, EviWhompingWillow, EviTomRiddlesDiary,
	},
	TagEviMap: unionTagEviMap(
		associateTagsToEvidence(EviDursleys, TagFamily, TagHome),
		associateTagsToEvidence(EviFlyingCar, TagEarth, TagSaturn),
		associateTagsToEvidence(EviDobby, TagMars, TagJupiter, TagMercury),
	),
	EviFindingsMap: unionEviFindingMap(
		associateEvidenceToFinding(FindingBook2Magic, EviDobby, EviFlyingCar, EviWhompingWillow),
		associateEvidenceToFinding(FindingBook2CGI, EviDobby, EviSpiderAragog),
	),
	Queries: []models.Query{
		QuerySalazarsHier,
		QueryWhereIsTheChamberOfSecrets,
	},
}

var newHPUser = newUserGen(1, func(f, l string) string { return strings.ToLower(f[:1] + l + "@hogwarts.edu") })
var UserDumbledore = newHPUser(newUserInput{FirstName: "Albus", LastName: "Dumbledore", Birthday: date(1970, 8, 1), SetLastUpdated: true, IsAdmin: true}) // birthday should be in 1881, but timestamp range is 1970-2038

var UserHarry = newHPUser(newUserInput{FirstName: "Harry", LastName: "Potter", Birthday: date(1980, 7, 31), SetLastUpdated: true})
var UserRon = newHPUser(newUserInput{FirstName: "Ronald", LastName: "Weasley", Birthday: date(1980, 3, 1), SetLastUpdated: true})
var UserGinny = newHPUser(newUserInput{FirstName: "Ginny", LastName: "Weasley", Birthday: date(1981, 3, 1), SetLastUpdated: true})
var UserHermione = newHPUser(newUserInput{FirstName: "Hermione", LastName: "Granger", Birthday: date(1979, 9, 19), SetLastUpdated: true})
var UserNeville = newHPUser(newUserInput{FirstName: "Neville", LastName: "Longbottom", Birthday: date(1979, 9, 19), SetLastUpdated: true})
var UserDraco = newHPUser(newUserInput{FirstName: "Draco", LastName: "Malfoy", Birthday: date(1980, 6, 5), SetLastUpdated: true})
var UserSnape = newHPUser(newUserInput{FirstName: "Serverus", LastName: "Snape", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserHagrid = newHPUser(newUserInput{FirstName: "Rubeus", LastName: "Hagrid", Birthday: date(1980, 1, 1), SetLastUpdated: true, Disabled: true})
var UserTomRiddle = newHPUser(newUserInput{FirstName: "Tom", LastName: "Riddle", Birthday: date(1980, 1, 1), SetLastUpdated: true, Deleted: true})
var UserHeadlessNick = newHPUser(newUserInput{FirstName: "Nicholas", LastName: "de Mimsy-Porpington", Birthday: date(1980, 1, 1), SetLastUpdated: true, Headless: true})

var newAPIKey = newAPIKeyGen(1)
var APIKeyHarry1 = newAPIKey(UserHarry.ID, "harry-abc", []byte{0x01, 0x02, 0x03})
var APIKeyHarry2 = newAPIKey(UserHarry.ID, "harry-123", []byte{0x11, 0x12, 0x13})
var APIKeyRon1 = newAPIKey(UserRon.ID, "ron-abc", []byte{0x01, 0x02, 0x03})
var APIKeyRon2 = newAPIKey(UserRon.ID, "ron-123", []byte{0x11, 0x12, 0x13})

var newHPOp = newOperationGen(1)
var OpSorcerersStone = newHPOp("HPSS", "Harry Potter and The Sorcerer's Stone")
var OpChamberOfSecrets = newHPOp("HPCoS", "Harry Potter and The Chamber Of Secrets")
var OpPrisonerOfAzkaban = newHPOp("HPPoA", "Harry Potter and The Prisoner Of Azkaban")
var OpGobletOfFire = newHPOp("HPGoF", "Harry Potter and The Goblet Of Fire")
var OpOrderOfThePhoenix = newHPOp("HPOotP", "Harry Potter and The Order Of The Phoenix")
var OpHalfBloodPrince = newHPOp("HPHBP", "Harry Potter and The Half Blood Prince")
var OpDeathlyHallows = newHPOp("HPDH", "Harry Potter and The Deathly Hallows")

var newHPTag = newTagGen(1)
var TagFamily = newHPTag(OpSorcerersStone.ID, "Family", "red")
var TagFriendship = newHPTag(OpSorcerersStone.ID, "Friendship", "orange")
var TagHome = newHPTag(OpSorcerersStone.ID, "Home", "yellow")
var TagLoyalty = newHPTag(OpSorcerersStone.ID, "Loyalty", "green")
var TagCourage = newHPTag(OpSorcerersStone.ID, "Courage", "blue")
var TagGoodVsEvil = newHPTag(OpSorcerersStone.ID, "Good vs. Evil", "indigo")
var TagSupernatural = newHPTag(OpSorcerersStone.ID, "Super Natural", "violet")

var TagMercury = newHPTag(OpChamberOfSecrets.ID, "Mercury", "violet")
var TagVenus = newHPTag(OpChamberOfSecrets.ID, "Venus", "red")
var TagEarth = newHPTag(OpChamberOfSecrets.ID, "Earth", "orange")
var TagMars = newHPTag(OpChamberOfSecrets.ID, "Mars", "yellow")
var TagJupiter = newHPTag(OpChamberOfSecrets.ID, "Jupiter", "green")
var TagSaturn = newHPTag(OpChamberOfSecrets.ID, "Saturn", "blue")
var TagNeptune = newHPTag(OpChamberOfSecrets.ID, "Neptune", "indigo")

var newHPEvidence = newEvidenceGen(1)
var EviDursleys = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "evi-uuid-dursleys", "Dursleys take care of young harry", "image")
var EviMirrorOfErised = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "evi-uuid-mirror", "Harry see parents in the mirror", "image")

var EviFlyingCar = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-flyingcar", "A Car that flies", "image")
var EviDobby = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-dobby", "an elf?", "image")
var EviSpiderAragog = newHPEvidence(OpChamberOfSecrets.ID, UserHagrid.ID, "evi-uuid-spider", "Just a big spider", "image")
var EviMoaningMyrtle = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-myrtle", "She's very sad", "image")
var EviWhompingWillow = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-willow", "Don't get too close", "image")
var EviTomRiddlesDiary = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-diary", "What's a Horcrux?", "codeblock")

var newHPQuery = newQueryGen(1)
var QuerySalazarsHier = newHPQuery(OpChamberOfSecrets.ID, "Find Heir", "Magic Query String", "findings")
var QueryWhereIsTheChamberOfSecrets = newHPQuery(OpChamberOfSecrets.ID, "Locate Chamber", "Fancy Query", "evidence")

var newHPFinding = newFindingGen(1)
var noLink = ""
var spiderLink = "https://www.google.com/search?q=spider+predators"
var FindingBook2Magic = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-b2magic", "some-category", "lots o' magic", "Magic plagues Harry's life", nil)
var FindingBook2CGI = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-cgi", "alt-category", "this looks fake", "I'm not entirely sure this is all above board", &noLink)
var FindingBook2SpiderFear = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-spider", "some-category", "how to scare spiders", "Who would have thought?", &spiderLink)
