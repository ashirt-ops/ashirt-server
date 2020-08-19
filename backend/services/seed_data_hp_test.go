// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"strings"

	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
)

var HarryPotterSeedData = TestSeedData{
	Users:      []models.User{UserHarry, UserRon, UserGinny, UserHermione, UserNeville, UserSeamus, UserDraco, UserSnape, UserDumbledore, UserHagrid, UserTomRiddle, UserHeadlessNick},
	Operations: []models.Operation{OpSorcerersStone, OpChamberOfSecrets, OpPrisonerOfAzkaban, OpGobletOfFire, OpOrderOfThePhoenix, OpHalfBloodPrince, OpDeathlyHallows, OpGanttChart},
	Tags: []models.Tag{
		TagFamily, TagFriendship, TagHome, TagLoyalty, TagCourage, TagGoodVsEvil, TagSupernatural,
		TagMercury, TagVenus, TagEarth, TagMars, TagJupiter, TagSaturn, TagNeptune,

		// common tags among all operations
		CommonTagWhoSS, CommonTagWhatSS, CommonTagWhereSS, CommonTagWhenSS, CommonTagWhySS,
		CommonTagWhoCoS, CommonTagWhatCoS, CommonTagWhereCoS, CommonTagWhenCoS, CommonTagWhyCoS,
		CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagHowGantt,
	},
	APIKeys: []models.APIKey{
		APIKeyHarry1, APIKeyHarry2,
		APIKeyRon1, APIKeyRon2,
	},
	UserOpMap: []models.UserOperationPermission{
		// OpSorcerersStone and OpChamberOfSecrets are used to check read/write permissions
		// The following should always remain true:
		// | User     | SS Perm | CoS Perm |
		// | -------- | ------- | -------- |
		// | Harry    | Admin   | Write    |
		// | Ron      | Write   | Admin    |
		// | Hermione | Read    | Write    |
		// | Seamus   | Write   | Read     |
		// | Ginny    | <none>  | Write    |
		// | Neville  | Write   | <none>   |

		newUserOpPermission(UserHarry, OpSorcerersStone, policy.OperationRoleAdmin),
		newUserOpPermission(UserRon, OpSorcerersStone, policy.OperationRoleWrite),
		newUserOpPermission(UserSeamus, OpSorcerersStone, policy.OperationRoleWrite),
		newUserOpPermission(UserHermione, OpSorcerersStone, policy.OperationRoleRead),
		newUserOpPermission(UserNeville, OpSorcerersStone, policy.OperationRoleWrite),

		newUserOpPermission(UserRon, OpChamberOfSecrets, policy.OperationRoleAdmin),
		newUserOpPermission(UserHarry, OpChamberOfSecrets, policy.OperationRoleWrite),
		newUserOpPermission(UserHermione, OpChamberOfSecrets, policy.OperationRoleWrite),
		newUserOpPermission(UserSeamus, OpChamberOfSecrets, policy.OperationRoleRead),
		newUserOpPermission(UserGinny, OpChamberOfSecrets, policy.OperationRoleWrite),

		newUserOpPermission(UserDumbledore, OpSorcerersStone, policy.OperationRoleAdmin),
		newUserOpPermission(UserDumbledore, OpChamberOfSecrets, policy.OperationRoleAdmin),

		newUserOpPermission(UserDumbledore, OpGanttChart, policy.OperationRoleAdmin),
		newUserOpPermission(UserHarry, OpGanttChart, policy.OperationRoleWrite),
		newUserOpPermission(UserGinny, OpGanttChart, policy.OperationRoleRead),
	},
	Findings: []models.Finding{
		FindingBook2Magic, FindingBook2CGI, FindingBook2SpiderFear,
	},
	Evidences: []models.Evidence{
		EviDursleys, EviMirrorOfErised,
		EviFlyingCar, EviDobby, EviSpiderAragog, EviMoaningMyrtle, EviWhompingWillow, EviTomRiddlesDiary, EviPetrifiedHermione,
		EviGanttZero, EviGanttOne, EviGanttTwo, EviGanttExtra, EviGanttThree, EviGanttFour,
		EviGanttLong00, EviGanttLong01, EviGanttLong02, EviGanttLong03, EviGanttLong04, EviGanttLong05, EviGanttLong06, EviGanttLong07, EviGanttLong08, EviGanttLong09,
		EviGanttLong10, EviGanttLong11, EviGanttLong12, EviGanttLong13, EviGanttLong14, EviGanttLong15, EviGanttLong16, EviGanttLong17, EviGanttLong18, EviGanttLong19,
		EviGanttLong20, EviGanttLong21, EviGanttLong22, EviGanttLong23, EviGanttLong24, EviGanttLong25, EviGanttLong26, EviGanttLong27, EviGanttLong28, EviGanttLong29,
		EviGanttLong30, EviGanttLong31, EviGanttLong32, EviGanttLong33, EviGanttLong34, EviGanttLong35, EviGanttLong36, EviGanttLong37, EviGanttLong38, EviGanttLong39,
		EviGanttLong40, EviGanttLong41, EviGanttLong42, EviGanttLong43, EviGanttLong44, EviGanttLong45, EviGanttLong46, EviGanttLong47, EviGanttLong48, EviGanttLong49,
		EviGanttLong50, EviGanttLong51, EviGanttLong52, EviGanttLong53, EviGanttLong54, EviGanttLong55, EviGanttLong56, EviGanttLong57, EviGanttLong58, EviGanttLong59,
	},
	TagEviMap: unionTagEviMap(
		associateTagsToEvidence(EviDursleys, TagFamily, TagHome),
		associateTagsToEvidence(EviFlyingCar, TagEarth, TagSaturn),
		associateTagsToEvidence(EviDobby, TagMars, TagJupiter, TagMercury),
		associateTagsToEvidence(EviPetrifiedHermione, TagMars, CommonTagWhatCoS, CommonTagWhoCoS),

		// tags are in a pattern for easy test verification of operation overview:
		//       01234
		// who   #...#
		// what  .#.#.
		// where #####
		// when  .###.
		// why   ##.##
		// Extra is added in to verify multiple evidence on a single day reflected in count
		associateTagsToEvidence(EviGanttZero, CommonTagWhoGantt, CommonTagWhereGantt, CommonTagWhyGantt),
		associateTagsToEvidence(EviGanttOne, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhenGantt),
		associateTagsToEvidence(EviGanttTwo, CommonTagWhereGantt, CommonTagWhenGantt),
		associateTagsToEvidence(EviGanttExtra, CommonTagWhenGantt),
		associateTagsToEvidence(EviGanttThree, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt),
		associateTagsToEvidence(EviGanttFour, CommonTagWhoGantt, CommonTagWhereGantt, CommonTagWhyGantt),

		associateEvidenceToTag(TagHowGantt,
			EviGanttLong00, EviGanttLong01, EviGanttLong02, EviGanttLong03, EviGanttLong04, EviGanttLong05, EviGanttLong06, EviGanttLong07, EviGanttLong08, EviGanttLong09,
			EviGanttLong10, EviGanttLong11, EviGanttLong12, EviGanttLong13, EviGanttLong14, EviGanttLong15, EviGanttLong16, EviGanttLong17, EviGanttLong18, EviGanttLong19,
			EviGanttLong20, EviGanttLong21, EviGanttLong22, EviGanttLong23, EviGanttLong24, EviGanttLong25, EviGanttLong26, EviGanttLong27, EviGanttLong28, EviGanttLong29,
			EviGanttLong30, EviGanttLong31, EviGanttLong32, EviGanttLong33, EviGanttLong34, EviGanttLong35, EviGanttLong36, EviGanttLong37, EviGanttLong38, EviGanttLong39,
			EviGanttLong40, EviGanttLong41, EviGanttLong42, EviGanttLong43, EviGanttLong44, EviGanttLong45, EviGanttLong46, EviGanttLong47, EviGanttLong48, EviGanttLong49,
			EviGanttLong50, EviGanttLong51, EviGanttLong52, EviGanttLong53, EviGanttLong54, EviGanttLong55, EviGanttLong56, EviGanttLong57, EviGanttLong58, EviGanttLong59,
		),
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
var UserSeamus = newHPUser(newUserInput{FirstName: "Seamus", LastName: "Finnigan", Birthday: date(1980, 9, 1), SetLastUpdated: true})
var UserDraco = newHPUser(newUserInput{FirstName: "Draco", LastName: "Malfoy", Birthday: date(1980, 6, 5), SetLastUpdated: true})
var UserSnape = newHPUser(newUserInput{FirstName: "Serverus", LastName: "Snape", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserHagrid = newHPUser(newUserInput{FirstName: "Rubeus", LastName: "Hagrid", Birthday: date(1980, 1, 1), SetLastUpdated: true, Disabled: true})
var UserTomRiddle = newHPUser(newUserInput{FirstName: "Tom", LastName: "Riddle", Birthday: date(1980, 1, 1), SetLastUpdated: true, Deleted: true})
var UserHeadlessNick = newHPUser(newUserInput{FirstName: "Nicholas", LastName: "de Mimsy-Porpington", Birthday: date(1980, 1, 1), SetLastUpdated: true, Headless: true})

// Reserved users: Luna Lovegood (Create user test)

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
var OpGanttChart = newHPOp("HPGantt", "Harry Potter and The Curse of Admin Oversight")

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

var CommonTagWhoSS = newHPTag(OpSorcerersStone.ID, "Who", "lightRed")
var CommonTagWhatSS = newHPTag(OpSorcerersStone.ID, "What", "lightBlue")
var CommonTagWhereSS = newHPTag(OpSorcerersStone.ID, "Where", "lightGreen")
var CommonTagWhenSS = newHPTag(OpSorcerersStone.ID, "When", "lightIndigo")
var CommonTagWhySS = newHPTag(OpSorcerersStone.ID, "Why", "lightYellow")

var CommonTagWhoCoS = newHPTag(OpChamberOfSecrets.ID, "Who", "lightRed")
var CommonTagWhatCoS = newHPTag(OpChamberOfSecrets.ID, "What", "lightBlue")
var CommonTagWhereCoS = newHPTag(OpChamberOfSecrets.ID, "Where", "lightGreen")
var CommonTagWhenCoS = newHPTag(OpChamberOfSecrets.ID, "When", "lightIndigo")
var CommonTagWhyCoS = newHPTag(OpChamberOfSecrets.ID, "Why", "lightYellow")

var CommonTagWhoGantt = newHPTag(OpGanttChart.ID, "Who", "lightRed")
var CommonTagWhatGantt = newHPTag(OpGanttChart.ID, "What", "lightBlue")
var CommonTagWhereGantt = newHPTag(OpGanttChart.ID, "Where", "lightGreen")
var CommonTagWhenGantt = newHPTag(OpGanttChart.ID, "When", "lightIndigo")
var CommonTagWhyGantt = newHPTag(OpGanttChart.ID, "Why", "lightYellow")
var TagHowGantt = newHPTag(OpGanttChart.ID, "How", "lightTeal")

var newHPEvidence = newEvidenceGen(1)
var EviDursleys = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "evi-uuid-dursleys", "Dursleys take care of young harry", "image", 0)
var EviMirrorOfErised = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "evi-uuid-mirror", "Harry see parents in the mirror", "image", 0)

var EviFlyingCar = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-flyingcar", "A Car that flies", "image", 0)
var EviDobby = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-dobby", "an elf?", "image", 0)
var EviSpiderAragog = newHPEvidence(OpChamberOfSecrets.ID, UserHagrid.ID, "evi-uuid-spider", "Just a big spider", "image", 0)
var EviMoaningMyrtle = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-myrtle", "She's very sad", "image", 0)
var EviWhompingWillow = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-willow", "Don't get too close", "image", 0)
var EviTomRiddlesDiary = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-diary", "What's a Horcrux?", "codeblock", 0)

var EviPetrifiedHermione = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "evi-uuid-rockher", "Strangely real-looking statue", "image", 0)

var EviGanttZero = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-uuid-one", "", "none", -4)
var EviGanttOne = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-uuid-two", "", "none", -3)
var EviGanttTwo = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-uuid-three", "", "none", -2)
var EviGanttThree = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-uuid-four", "", "none", -1)
var EviGanttFour = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-uuid-five", "", "none", 0)
var EviGanttExtra = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-uuid-extra", "", "none", -2)

// These "long" pieces of evidence test a sql issue experienced using group_concat (no longer used)
var EviGanttLong00 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-00", "", "none", 0)
var EviGanttLong01 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-01", "", "none", -1)
var EviGanttLong02 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-02", "", "none", -2)
var EviGanttLong03 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-03", "", "none", -3)
var EviGanttLong04 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-04", "", "none", -4)
var EviGanttLong05 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-05", "", "none", -5)
var EviGanttLong06 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-06", "", "none", -6)
var EviGanttLong07 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-07", "", "none", -7)
var EviGanttLong08 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-08", "", "none", -8)
var EviGanttLong09 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-09", "", "none", -9)
var EviGanttLong10 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-10", "", "none", -10)
var EviGanttLong11 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-11", "", "none", -11)
var EviGanttLong12 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-12", "", "none", -12)
var EviGanttLong13 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-13", "", "none", -13)
var EviGanttLong14 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-14", "", "none", -14)
var EviGanttLong15 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-15", "", "none", -15)
var EviGanttLong16 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-16", "", "none", -16)
var EviGanttLong17 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-17", "", "none", -17)
var EviGanttLong18 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-18", "", "none", -18)
var EviGanttLong19 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-19", "", "none", -19)
var EviGanttLong20 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-20", "", "none", -20)
var EviGanttLong21 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-21", "", "none", -21)
var EviGanttLong22 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-22", "", "none", -22)
var EviGanttLong23 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-23", "", "none", -23)
var EviGanttLong24 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-24", "", "none", -24)
var EviGanttLong25 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-25", "", "none", -25)
var EviGanttLong26 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-26", "", "none", -26)
var EviGanttLong27 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-27", "", "none", -27)
var EviGanttLong28 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-28", "", "none", -28)
var EviGanttLong29 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-29", "", "none", -29)
var EviGanttLong30 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-30", "", "none", -30)
var EviGanttLong31 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-31", "", "none", -31)
var EviGanttLong32 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-32", "", "none", -32)
var EviGanttLong33 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-33", "", "none", -33)
var EviGanttLong34 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-34", "", "none", -34)
var EviGanttLong35 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-35", "", "none", -35)
var EviGanttLong36 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-36", "", "none", -36)
var EviGanttLong37 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-37", "", "none", -37)
var EviGanttLong38 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-38", "", "none", -38)
var EviGanttLong39 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-39", "", "none", -39)
var EviGanttLong40 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-40", "", "none", -40)
var EviGanttLong41 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-41", "", "none", -41)
var EviGanttLong42 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-42", "", "none", -42)
var EviGanttLong43 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-43", "", "none", -43)
var EviGanttLong44 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-44", "", "none", -44)
var EviGanttLong45 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-45", "", "none", -45)
var EviGanttLong46 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-46", "", "none", -46)
var EviGanttLong47 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-47", "", "none", -47)
var EviGanttLong48 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-48", "", "none", -48)
var EviGanttLong49 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-49", "", "none", -49)
var EviGanttLong50 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-50", "", "none", -50)
var EviGanttLong51 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-51", "", "none", -51)
var EviGanttLong52 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-52", "", "none", -52)
var EviGanttLong53 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-53", "", "none", -53)
var EviGanttLong54 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-54", "", "none", -54)
var EviGanttLong55 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-55", "", "none", -55)
var EviGanttLong56 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-56", "", "none", -56)
var EviGanttLong57 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-57", "", "none", -57)
var EviGanttLong58 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-58", "", "none", -58)
var EviGanttLong59 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "evi-gantt-59", "", "none", -59)

var newHPQuery = newQueryGen(1)
var QuerySalazarsHier = newHPQuery(OpChamberOfSecrets.ID, "Find Heir", "Magic Query String", "findings")
var QueryWhereIsTheChamberOfSecrets = newHPQuery(OpChamberOfSecrets.ID, "Locate Chamber", "Fancy Query", "evidence")

var newHPFinding = newFindingGen(1)
var noLink = ""
var spiderLink = "https://www.google.com/search?q=spider+predators"
var FindingBook2Magic = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-b2magic", "some-category", "lots o' magic", "Magic plagues Harry's life", nil)
var FindingBook2CGI = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-cgi", "alt-category", "this looks fake", "I'm not entirely sure this is all above board", &noLink)
var FindingBook2SpiderFear = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-spider", "some-category", "how to scare spiders", "Who would have thought?", &spiderLink)
