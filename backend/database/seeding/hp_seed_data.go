// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package seeding

import (
	"strings"

	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
)

var HarryPotterSeedData = Seeder{
	Users: []models.User{UserHarry, UserRon, UserGinny, UserHermione, UserNeville, UserSeamus, UserDraco, UserSnape, UserDumbledore, UserHagrid, UserTomRiddle, UserHeadlessNick,
		UserCedric, UserFleur, UserViktor, UserAlastor, UserMinerva, UserLucius, UserSirius, UserPeter, UserParvati, UserPadma, UserCho,
	},
	Operations: []models.Operation{OpSorcerersStone, OpChamberOfSecrets, OpPrisonerOfAzkaban, OpGobletOfFire, OpOrderOfThePhoenix, OpHalfBloodPrince, OpDeathlyHallows, OpGanttChart},
	Tags: []models.Tag{
		TagFamily, TagFriendship, TagHome, TagLoyalty, TagCourage, TagGoodVsEvil, TagSupernatural,
		TagMercury, TagVenus, TagEarth, TagMars, TagJupiter, TagSaturn, TagNeptune,

		// common tags among all operations
		CommonTagWhoSS, CommonTagWhatSS, CommonTagWhereSS, CommonTagWhenSS, CommonTagWhySS,
		CommonTagWhoCoS, CommonTagWhatCoS, CommonTagWhereCoS, CommonTagWhenCoS, CommonTagWhyCoS,
		CommonTagWhoGoF, CommonTagWhatGoF, CommonTagWhereGoF, CommonTagWhenGoF, CommonTagWhyGoF,
		CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork, TagGanttWalk,
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

		// Every user should be part of OpGobletOfFire (exception: deleted users)
		newUserOpPermission(UserHarry, OpGobletOfFire, policy.OperationRoleAdmin),
		newUserOpPermission(UserRon, OpGobletOfFire, policy.OperationRoleWrite),
		newUserOpPermission(UserGinny, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserHermione, OpGobletOfFire, policy.OperationRoleWrite),
		newUserOpPermission(UserNeville, OpGobletOfFire, policy.OperationRoleWrite),
		newUserOpPermission(UserSeamus, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserDraco, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserSnape, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserHagrid, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserHeadlessNick, OpGobletOfFire, policy.OperationRoleWrite),
		newUserOpPermission(UserCedric, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserFleur, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserViktor, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserAlastor, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserMinerva, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserLucius, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserSirius, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserPeter, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserParvati, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserPadma, OpGobletOfFire, policy.OperationRoleRead),
		newUserOpPermission(UserCho, OpGobletOfFire, policy.OperationRoleRead),

		newUserOpPermission(UserDumbledore, OpSorcerersStone, policy.OperationRoleAdmin),
		newUserOpPermission(UserDumbledore, OpChamberOfSecrets, policy.OperationRoleAdmin),
		newUserOpPermission(UserDumbledore, OpGobletOfFire, policy.OperationRoleAdmin),

		newUserOpPermission(UserDumbledore, OpGanttChart, policy.OperationRoleAdmin),
		newUserOpPermission(UserHarry, OpGanttChart, policy.OperationRoleWrite),
		newUserOpPermission(UserGinny, OpGanttChart, policy.OperationRoleRead),
	},
	Findings: []models.Finding{
		FindingBook2Magic, FindingBook2CGI, FindingBook2SpiderFear,
	},
	Evidences: []models.Evidence{
		EviDursleys, EviMirrorOfErised, EviLevitateSpell, EviRulesForQuidditch,
		EviFlyingCar, EviDobby, EviSpiderAragog, EviMoaningMyrtle, EviWhompingWillow, EviTomRiddlesDiary, EviPetrifiedHermione, EviHeadlessHuntApplication,
		EviTristateTrophy, EviEntryForm, EviWizardDance, EviPolyjuice, EviWarewolf,
		EviGantt01, EviGantt02, EviGantt03, EviGantt04, EviGantt05, EviGantt06, EviGantt07, EviGantt08, EviGantt09, EviGantt10,
		EviGantt11, EviGantt12, EviGantt13, EviGantt14, EviGantt15, EviGantt16, EviGantt17, EviGantt18, EviGantt19, EviGantt20, EviGanttExtra,
	},
	TagEviMap: unionTagEviMap(
		associateTagsToEvidence(EviDursleys, TagFamily, TagHome),
		associateTagsToEvidence(EviFlyingCar, TagEarth, TagSaturn),
		associateTagsToEvidence(EviDobby, TagMars, TagJupiter, TagMercury),
		associateTagsToEvidence(EviPetrifiedHermione, TagMars, CommonTagWhatCoS, CommonTagWhoCoS),

		associateTagsToEvidence(EviTristateTrophy, CommonTagWhoGoF, CommonTagWhereGoF, CommonTagWhyGoF),
		associateTagsToEvidence(EviEntryForm, CommonTagWhatGoF, CommonTagWhereGoF, CommonTagWhenGoF),
		associateTagsToEvidence(EviWizardDance, CommonTagWhereGoF, CommonTagWhenGoF),
		associateTagsToEvidence(EviPolyjuice, CommonTagWhatGoF, CommonTagWhereGoF, CommonTagWhenGoF, CommonTagWhyGoF),
		associateTagsToEvidence(EviWarewolf, CommonTagWhoGoF, CommonTagWhereGoF, CommonTagWhyGoF),

		// tags are in a pattern: the first 10 columns are dedicated to a pineapple, the second to an apple
		// -- pineapple      apple
		//    1234567890     1234567890
		//  1 .###.###..     ......##..
		//  2 #########.     .....##...
		//  3 #..###..#.     ..##.#.##.
		//  4 ..#####...     .####.####
		//  5 .#.#.#.#..     .#.#######
		//  6 .#######..     .#.#######
		//  7 .#.#.#.#..     .#.#######
		//  8 .#######..     .##.######
		//  9 .#.#.#.#..     ..#######.
		// 10 ..#####...     ...##.##..
		//  Note: EviGanttExtra is present to check multiple-tag-usage on same-day
		associateTagsToEvidence(EviGantt01, CommonTagWhatGantt, CommonTagWhereGantt),
		associateTagsToEvidence(EviGantt02, CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork),
		associateTagsToEvidence(EviGantt03, CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhenGantt, TagGanttBroom, TagGanttAparate, TagGanttWalk),
		associateTagsToEvidence(EviGantt04, CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork, TagGanttWalk),
		associateTagsToEvidence(EviGantt05, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhenGantt, TagGanttBroom, TagGanttAparate, TagGanttWalk),
		associateTagsToEvidence(EviGantt06, CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork, TagGanttWalk),
		associateTagsToEvidence(EviGantt07, CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhenGantt, TagGanttBroom, TagGanttAparate, TagGanttWalk),
		associateTagsToEvidence(EviGantt08, CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork),
		associateTagsToEvidence(EviGantt09, CommonTagWhatGantt, CommonTagWhereGantt),
		associateTagsToEvidence(EviGantt10),
		associateTagsToEvidence(EviGantt11),
		associateTagsToEvidence(EviGantt12, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate),
		associateTagsToEvidence(EviGantt13, CommonTagWhereGantt, CommonTagWhenGantt, TagGanttAparate, TagGanttFlooNetwork),
		associateTagsToEvidence(EviGantt14, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttFlooNetwork, TagGanttWalk),
		associateTagsToEvidence(EviGantt15, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork, TagGanttWalk),
		associateTagsToEvidence(EviGantt16, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork),
		associateTagsToEvidence(EviGantt17, CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork, TagGanttWalk),
		associateTagsToEvidence(EviGantt18, CommonTagWhoGantt, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork, TagGanttWalk),
		associateTagsToEvidence(EviGantt19, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork),
		associateTagsToEvidence(EviGantt20, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate),
		associateTagsToEvidence(EviGanttExtra, CommonTagWhoGantt),
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

var newHPUser = newUserGen(1, func(f, l string) string { return strings.ToLower(f + "." + strings.Replace(l, " ", "", -1)) })

// UserDumbledore is reserved to be a super admin.
var UserDumbledore = newHPUser(newUserInput{FirstName: "Albus", LastName: "Dumbledore", Birthday: date(1970, 8, 1), SetLastUpdated: true, IsAdmin: true}) // birthday should be in 1881, but timestamp range is 1970-2038

var UserHarry = newHPUser(newUserInput{FirstName: "Harry", LastName: "Potter", Birthday: date(1980, 7, 31), SetLastUpdated: true})
var UserRon = newHPUser(newUserInput{FirstName: "Ron", LastName: "Weasley", Birthday: date(1980, 3, 1), SetLastUpdated: true})
var UserGinny = newHPUser(newUserInput{FirstName: "Ginny", LastName: "Weasley", Birthday: date(1981, 3, 1), SetLastUpdated: true})
var UserHermione = newHPUser(newUserInput{FirstName: "Hermione", LastName: "Granger", Birthday: date(1979, 9, 19), SetLastUpdated: true})
var UserNeville = newHPUser(newUserInput{FirstName: "Neville", LastName: "Longbottom", Birthday: date(1979, 9, 19), SetLastUpdated: true})
var UserSeamus = newHPUser(newUserInput{FirstName: "Seamus", LastName: "Finnigan", Birthday: date(1980, 9, 1), SetLastUpdated: true})
var UserDraco = newHPUser(newUserInput{FirstName: "Draco", LastName: "Malfoy", Birthday: date(1980, 6, 5), SetLastUpdated: true})
var UserSnape = newHPUser(newUserInput{FirstName: "Serverus", LastName: "Snape", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserCedric = newHPUser(newUserInput{FirstName: "Cedric", LastName: "Digory", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserFleur = newHPUser(newUserInput{FirstName: "Fleur", LastName: "Delacour", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserViktor = newHPUser(newUserInput{FirstName: "Viktor", LastName: "Krum", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserAlastor = newHPUser(newUserInput{FirstName: "Alastor", LastName: "Moody", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserMinerva = newHPUser(newUserInput{FirstName: "Minerva", LastName: "McGonagall", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserLucius = newHPUser(newUserInput{FirstName: "Lucius", LastName: "Malfoy", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserSirius = newHPUser(newUserInput{FirstName: "Sirius", LastName: "Black", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserPeter = newHPUser(newUserInput{FirstName: "Peter", LastName: "Pettigrew", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserParvati = newHPUser(newUserInput{FirstName: "Parvati", LastName: "Patil", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserPadma = newHPUser(newUserInput{FirstName: "Padma", LastName: "Patil", Birthday: date(1980, 1, 1), SetLastUpdated: true})
var UserCho = newHPUser(newUserInput{FirstName: "Cho", LastName: "Chang", Birthday: date(1980, 1, 1), SetLastUpdated: true})

// UserHagrid is reserved to test disabled users
var UserHagrid = newHPUser(newUserInput{FirstName: "Rubeus", LastName: "Hagrid", Birthday: date(1980, 1, 1), SetLastUpdated: true, Disabled: true})

// UserTomRiddle is reserved to test deleted users
var UserTomRiddle = newHPUser(newUserInput{FirstName: "Tom", LastName: "Riddle", Birthday: date(1980, 1, 1), SetLastUpdated: true, Deleted: true})

// UserHeadlessNick is reserved to test api-only access/"headless" users
var UserHeadlessNick = newHPUser(newUserInput{FirstName: "Nicholas", LastName: "de Mimsy-Porpington", Birthday: date(1980, 1, 1), SetLastUpdated: true, Headless: true})

// Reserved users: Luna Lovegood (Create user test)

var newAPIKey = newAPIKeyGen(1)
var APIKeyHarry1 = newAPIKey(UserHarry.ID, "harry-abc", []byte{0x01, 0x02, 0x03})
var APIKeyHarry2 = newAPIKey(UserHarry.ID, "harry-123", []byte{0x11, 0x12, 0x13})
var APIKeyRon1 = newAPIKey(UserRon.ID, "ron-abc", []byte{0x01, 0x02, 0x03})
var APIKeyRon2 = newAPIKey(UserRon.ID, "ron-123", []byte{0x11, 0x12, 0x13})
var APIKeyHeadlessNick1 = newAPIKey(UserHeadlessNick.ID, "DAYPFGHnm1Pqes-l0Fm76_y1", []byte("HqmuWylLznR+tqSotZAOc+w47buSFaKKTJozpXEYkuNBiuRJgw3NeJOuVP6kbQBQmiYTqYAaiIKbcO1BxcH52Q==")) // realKey

var newHPOp = newOperationGen(1)

// OpSorcerersStone is reserved to test permission issues. Also used with OpChamberOfSecrets
var OpSorcerersStone = newHPOp("HPSS", "Harry Potter and The Sorcerer's Stone")

// OpChamberOfSecrets is reserved to test permission issues. Also used with OpSorcerersStone
var OpChamberOfSecrets = newHPOp("HPCoS", "Harry Potter and The Chamber Of Secrets")

// OpPrisonerOfAzkaban is reserved as a no-user (orphaned) operation
var OpPrisonerOfAzkaban = newHPOp("HPPoA", "Harry Potter and The Prisoner Of Azkaban")

// OpGobletOfFire is reserved for a common-operation for all users
var OpGobletOfFire = newHPOp("HPGoF", "Harry Potter and The Goblet Of Fire")

// OpOrderOfThePhoenix is available for use
var OpOrderOfThePhoenix = newHPOp("HPOotP", "Harry Potter and The Order Of The Phoenix")

// OpHalfBloodPrince is available for use
var OpHalfBloodPrince = newHPOp("HPHBP", "Harry Potter and The Half Blood Prince")

// OpDeathlyHallows is a vailable for use
var OpDeathlyHallows = newHPOp("HPDH", "Harry Potter and The Deathly Hallows")

// OpGanttChart is reserved to verify the Overview feature
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

// Common tags are used to test migrating content from one operation to another
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

var CommonTagWhoGoF = newHPTag(OpGobletOfFire.ID, "Who", "lightRed")
var CommonTagWhatGoF = newHPTag(OpGobletOfFire.ID, "What", "lightBlue")
var CommonTagWhereGoF = newHPTag(OpGobletOfFire.ID, "Where", "lightGreen")
var CommonTagWhenGoF = newHPTag(OpGobletOfFire.ID, "When", "lightIndigo")
var CommonTagWhyGoF = newHPTag(OpGobletOfFire.ID, "Why", "lightYellow")

var CommonTagWhoGantt = newHPTag(OpGanttChart.ID, "Who", "lightRed")
var CommonTagWhatGantt = newHPTag(OpGanttChart.ID, "What", "lightBlue")
var CommonTagWhereGantt = newHPTag(OpGanttChart.ID, "Where", "lightGreen")
var CommonTagWhenGantt = newHPTag(OpGanttChart.ID, "When", "lightIndigo")
var CommonTagWhyGantt = newHPTag(OpGanttChart.ID, "Why", "lightYellow")
var TagGanttBroom = newHPTag(OpGanttChart.ID, "Broom", "red")
var TagGanttHippogriff = newHPTag(OpGanttChart.ID, "Hippogriff", "blue")
var TagGanttAparate = newHPTag(OpGanttChart.ID, "Aparate", "green")
var TagGanttFlooNetwork = newHPTag(OpGanttChart.ID, "Floo", "indigo")
var TagGanttWalk = newHPTag(OpGanttChart.ID, "Walk", "yellow")

var newHPEvidence = newEvidenceGen(1)
var EviDursleys = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_dursleys", "Family of self-centered muggles + Harry", "image", 0)
var EviMirrorOfErised = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_mirror", "Mysterious mirror that shows you your deepest desires", "image", 0)
var EviLevitateSpell = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_md_levitate", "Documented Levitation Spell effects", "codeblock", 0)
var EviRulesForQuidditch = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_rs_aoc201501", "Complex rules for a simple game", "codeblock", 0)

var EviFlyingCar = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_car", "A Car that flies", "image", 0)
var EviDobby = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_dobby", "an elf?", "image", 0)
var EviSpiderAragog = newHPEvidence(OpChamberOfSecrets.ID, UserHagrid.ID, "seed_aragog", "Just a big spider", "image", 0)
var EviMoaningMyrtle = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_myrtle", "She's very sad", "image", 0)
var EviWhompingWillow = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_tree", "Don't get too close", "image", 0)
var EviTomRiddlesDiary = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_go_aoc201614", "What's a Horcrux?", "codeblock", 0)
var EviHeadlessHuntApplication = newHPEvidence(OpChamberOfSecrets.ID, UserRon.ID, "seed_py_aoc201717", "This group is very particular", "codeblock", 0)
var EviPetrifiedHermione = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_statue", "Strangely real-looking statue", "image", 0)

var EviTristateTrophy = newHPEvidence(OpGobletOfFire.ID, UserHarry.ID, "seed_trophy", "First Triwizard Champion Trophy", "image", 0)
var EviEntryForm = newHPEvidence(OpGobletOfFire.ID, UserCedric.ID, "seed_entry", "Cedric's entry form for Triwizard competition", "codeblock", 0)
var EviWizardDance = newHPEvidence(OpGobletOfFire.ID, UserCho.ID, "seed_dance", "Advertising for the Triwizard Dance", "image", 0)
var EviPolyjuice = newHPEvidence(OpGobletOfFire.ID, UserAlastor.ID, "seed_juice", "DIY instructions for Polyjuice Potion", "codeblock", 0)
var EviWarewolf = newHPEvidence(OpGobletOfFire.ID, UserViktor.ID, "seed_wolf", "Strangely real-looking statue", "terminal-recording", 0)

var EviGantt01 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_01", "", "none", -19)
var EviGantt02 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_02", "", "none", -18)
var EviGantt03 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_03", "", "none", -17)
var EviGantt04 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_04", "", "none", -16)
var EviGantt05 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_05", "", "none", -15)
var EviGantt06 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_06", "", "none", -14)
var EviGantt07 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_07", "", "none", -13)
var EviGantt08 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_08", "", "none", -12)
var EviGantt09 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_09", "", "none", -11)
var EviGantt10 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_10", "", "none", -10)
var EviGantt11 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_11", "", "none", -9)
var EviGantt12 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_12", "", "none", -8)
var EviGantt13 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_13", "", "none", -7)
var EviGantt14 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_14", "", "none", -6)
var EviGantt15 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_15", "", "none", -5)
var EviGantt16 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_16", "", "none", -4)
var EviGantt17 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_17", "", "none", -3)
var EviGantt18 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_18", "", "none", -2)
var EviGantt19 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_19", "", "none", -1)
var EviGantt20 = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_20", "", "none", 0)
var EviGanttExtra = newHPEvidence(OpGanttChart.ID, UserHarry.ID, "seed_gantt_extra", "", "none", -19)

var newHPQuery = newQueryGen(1)
var QuerySalazarsHier = newHPQuery(OpChamberOfSecrets.ID, "Find Heir", "Magic Query String", "findings")
var QueryWhereIsTheChamberOfSecrets = newHPQuery(OpChamberOfSecrets.ID, "Locate Chamber", "Fancy Query", "evidence")

var newHPFinding = newFindingGen(1)
var noLink = ""
var spiderLink = "https://www.google.com/search?q=spider+predators"
var FindingBook2Magic = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-b2magic", "some-category", "lots o' magic", "Magic plagues Harry's life", nil)
var FindingBook2CGI = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-cgi", "alt-category", "this looks fake", "I'm not entirely sure this is all above board", &noLink)
var FindingBook2SpiderFear = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-spider", "some-category", "how to scare spiders", "Who would have thought?", &spiderLink)
