// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package seeding

import (
	"strings"

	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
)

var HarryPotterSeedData = Seeder{
	FindingCategories: []models.FindingCategory{
		ProductFindingCategory, NetworkFindingCategory, EnterpriseFindingCategory, VendorFindingCategory, BehavioralFindingCategory, DetectionGapFindingCategory,
		DeletedCategory, SomeFindingCategory, SomeOtherFindingCategory,
	},
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
		CommonTagWhoGantt, CommonTagWhatGantt, CommonTagWhereGantt, CommonTagWhenGantt, CommonTagWhyGantt, TagGanttBroom, TagGanttHippogriff, TagGanttAparate, TagGanttFlooNetwork, TagGanttWalk, TagHowGantt,
	},
	DefaultTags: []models.DefaultTag{
		DefaultTagWho, DefaultTagWhat, DefaultTagWhere, DefaultTagWhen, DefaultTagWhy,
	},
	APIKeys: []models.APIKey{
		APIKeyHarry1, APIKeyHarry2,
		APIKeyRon1, APIKeyRon2,
		APIKeyNick,
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
		FindingBook2Magic, FindingBook2CGI, FindingBook2SpiderFear, FindingBook2Robes,
	},
	Evidences: []models.Evidence{
		EviDursleys, EviMirrorOfErised, EviLevitateSpell, EviRulesForQuidditch,
		EviFlyingCar, EviDobby, EviSpiderAragog, EviMoaningMyrtle, EviWhompingWillow, EviTomRiddlesDiary, EviPetrifiedHermione, EviHeadlessHuntApplication,
		EviTristateTrophy, EviEntryForm, EviWizardDance, EviPolyjuice, EviWarewolf,
		EviGantt01, EviGantt02, EviGantt03, EviGantt04, EviGantt05, EviGantt06, EviGantt07, EviGantt08, EviGantt09, EviGantt10,
		EviGantt11, EviGantt12, EviGantt13, EviGantt14, EviGantt15, EviGantt16, EviGantt17, EviGantt18, EviGantt19, EviGantt20, EviGanttExtra,

		EviGanttLong00, EviGanttLong01, EviGanttLong02, EviGanttLong03, EviGanttLong04, EviGanttLong05, EviGanttLong06, EviGanttLong07, EviGanttLong08, EviGanttLong09,
		EviGanttLong10, EviGanttLong11, EviGanttLong12, EviGanttLong13, EviGanttLong14, EviGanttLong15, EviGanttLong16, EviGanttLong17, EviGanttLong18, EviGanttLong19,
		EviGanttLong20, EviGanttLong21, EviGanttLong22, EviGanttLong23, EviGanttLong24, EviGanttLong25, EviGanttLong26, EviGanttLong27, EviGanttLong28, EviGanttLong29,
		EviGanttLong30, EviGanttLong31, EviGanttLong32, EviGanttLong33, EviGanttLong34, EviGanttLong35, EviGanttLong36, EviGanttLong37, EviGanttLong38, EviGanttLong39,
		EviGanttLong40, EviGanttLong41, EviGanttLong42, EviGanttLong43, EviGanttLong44, EviGanttLong45, EviGanttLong46, EviGanttLong47, EviGanttLong48, EviGanttLong49,
		EviGanttLong50, EviGanttLong51, EviGanttLong52, EviGanttLong53, EviGanttLong54, EviGanttLong55, EviGanttLong56, EviGanttLong57, EviGanttLong58, EviGanttLong59,
	},
	EvidenceMetadatas: []models.EvidenceMetadata{
		EviMetaDursleys, EviMetaMirrorOfErised, EviMetaLevitateSpell, EviMetaRulesForQuidditch, EviMetaRulesForQuidditchTwo,
		EviMetaFlyingCar, EviMetaDobby, EviMetaSpiderAragog, EviMetaMoaningMyrtle, EviMetaWhompingWillow, EviMetaTomRiddlesDiary, EviMetaTomRiddlesDiaryTwo, EviMetaPetrifiedHermione, EviMetaHeadlessHuntApplication,
		EviMetaTristateTrophy, EviMetaEntryForm, EviMetaWizardDance, EviMetaPolyjuice, EviMetaEntryFormTwo, EviMetaWarewolf,
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
	ServiceWorkers: []models.ServiceWorker{
		demoServiceWorker,
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
var APIKeyNick = newAPIKey(UserHeadlessNick.ID, "gR6nVtaQmp2SvzIqLUWdedDk", []byte{
	// Corresponds to secret key: WvtvxFaJS0mPs82nCzqamI+bOGXpq7EIQhg4UD8nxS5448XG9N0gNAceJGBLPdCA3kAzC4MdUSHnKCJ/lZD++A==
	0x5A, 0xFB, 0x6F, 0xC4, 0x56, 0x89, 0x4B, 0x49, 0x8F, 0xB3, 0xCD, 0xA7, 0x0B, 0x3A, 0x9A, 0x98,
	0x8F, 0x9B, 0x38, 0x65, 0xE9, 0xAB, 0xB1, 0x08, 0x42, 0x18, 0x38, 0x50, 0x3F, 0x27, 0xC5, 0x2E,
	0x78, 0xE3, 0xC5, 0xC6, 0xF4, 0xDD, 0x20, 0x34, 0x07, 0x1E, 0x24, 0x60, 0x4B, 0x3D, 0xD0, 0x80,
	0xDE, 0x40, 0x33, 0x0B, 0x83, 0x1D, 0x51, 0x21, 0xE7, 0x28, 0x22, 0x7F, 0x95, 0x90, 0xFE, 0xF8,
})

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

var newDefaultHPTag = newDefaultTagGen(1)
var DefaultTagWho = newDefaultHPTag("Who", "lightRed")
var DefaultTagWhat = newDefaultHPTag("What", "lightBlue")
var DefaultTagWhere = newDefaultHPTag("Where", "lightGreen")
var DefaultTagWhen = newDefaultHPTag("When", "lightIndigo")
var DefaultTagWhy = newDefaultHPTag("Why", "lightYellow")

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
var TagHowGantt = newHPTag(OpGanttChart.ID, "How", "lightTeal")

var newHPEvidence = newEvidenceGen(1)
var EviDursleys = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_dursleys", "Family of self-centered muggles + Harry", "image", 0)
var EviMirrorOfErised = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_mirror", "Mysterious mirror that shows you your deepest desires", "image", 0)
var EviLevitateSpell = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_md_levitate", "Documented Levitation Spell effects", "codeblock", 0)
var EviRulesForQuidditch = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_rs_aoc201501", "Complex rules for a simple game", "codeblock", 0)

var newHPEviMetadata = newEvidenceMetadataGen(1)
var EviMetaDursleys = newHPEviMetadata(EviDursleys.ID, "color-averager", "rgb(65, 65, 65)\n#414141\nhsl(0, 0%, 25%)", 0)
var EviMetaMirrorOfErised = newHPEviMetadata(EviMirrorOfErised.ID, "color-averager", "rgb(111, 77, 14)\n#6f4d0e\nhsl(39, 78%, 25%)", 0)
var EviMetaLevitateSpell = newHPEviMetadata(EviLevitateSpell.ID, "wc -l", "12 seed_md_levitate", 0)
var EviMetaRulesForQuidditch = newHPEviMetadata(EviRulesForQuidditch.ID, "run-result", "Last floor reached: 138 \n Step on which basement is reached (first time): 1771", 0)
var EviMetaRulesForQuidditchTwo = newHPEviMetadata(EviRulesForQuidditch.ID, "wc-l", "33 main.rs", 0)

var EviFlyingCar = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_car", "A Car that flies", "image", 0)
var EviDobby = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_dobby", "an elf?", "image", 0)
var EviSpiderAragog = newHPEvidence(OpChamberOfSecrets.ID, UserHagrid.ID, "seed_aragog", "Just a big spider", "image", 0)
var EviMoaningMyrtle = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_myrtle", "She's very sad", "image", 0)
var EviWhompingWillow = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_tree", "Don't get too close", "image", 0)
var EviTomRiddlesDiary = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_go_aoc201614", "What's a Horcrux?", "codeblock", 0)
var EviHeadlessHuntApplication = newHPEvidence(OpChamberOfSecrets.ID, UserRon.ID, "seed_py_aoc201717", "This group is very particular", "codeblock", 0)
var EviPetrifiedHermione = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_statue", "Strangely real-looking statue", "image", 0)

var EviMetaFlyingCar = newHPEviMetadata(EviFlyingCar.ID, "color-averager", "rgb(106, 109, 84)\n#6a6d54\nhsl(67, 13%, 38%)", 0)
var EviMetaDobby = newHPEviMetadata(EviDobby.ID, "color-averager", "rgb(74, 51, 32)\n#4a3320\nhsl(27, 40%, 21%)", 0)
var EviMetaSpiderAragog = newHPEviMetadata(EviSpiderAragog.ID, "color-averager", "rgb(189, 156, 146)\n#bd9c92\nhsl(14, 25%, 66%)", 0)
var EviMetaMoaningMyrtle = newHPEviMetadata(EviMoaningMyrtle.ID, "color-averager", "rgb(118, 103, 102)\n#766766\nhsl(4, 7%, 43%)", 0)
var EviMetaWhompingWillow = newHPEviMetadata(EviWhompingWillow.ID, "color-averager", "rgb(115, 109, 81)\n#736d51\nhsl(49, 17%, 38%)", 0)
var EviMetaTomRiddlesDiary = newHPEviMetadata(EviTomRiddlesDiary.ID, "run-result", "All keys found by index:  19968", 0)
var EviMetaTomRiddlesDiaryTwo = newHPEviMetadata(EviTomRiddlesDiary.ID, "wc -l", "98 main.go", 0)
var EviMetaHeadlessHuntApplication = newHPEviMetadata(EviHeadlessHuntApplication.ID, "run-result", "41797835\nelapsed time (seconds): 3.772843360900879", 0)
var EviMetaPetrifiedHermione = newHPEviMetadata(EviPetrifiedHermione.ID, "color-averager", "rgb(162, 104, 101)\n#a26865\nhsl(3, 25%, 52%)", 0)

var EviTristateTrophy = newHPEvidence(OpGobletOfFire.ID, UserHarry.ID, "seed_trophy", "First Triwizard Champion Trophy", "image", 0)
var EviEntryForm = newHPEvidence(OpGobletOfFire.ID, UserCedric.ID, "seed_entry", "Cedric's entry form for Triwizard competition", "codeblock", 0)
var EviWizardDance = newHPEvidence(OpGobletOfFire.ID, UserCho.ID, "seed_dance", "Advertising for the Triwizard Dance", "image", 0)
var EviPolyjuice = newHPEvidence(OpGobletOfFire.ID, UserAlastor.ID, "seed_juice", "DIY instructions for Polyjuice Potion", "codeblock", 0)
var EviWarewolf = newHPEvidence(OpGobletOfFire.ID, UserViktor.ID, "seed_wolf", "Strangely real-looking statue", "terminal-recording", 0)

var EviMetaTristateTrophy = newHPEviMetadata(EviTristateTrophy.ID, "color-averager", "rgb(182, 184, 183)\n#b6b8b7\nhsl(150, 1%, 72%)", 0)
var EviMetaEntryForm = newHPEviMetadata(EviEntryForm.ID, "run-result", "(No output)", 0)
var EviMetaEntryFormTwo = newHPEviMetadata(EviEntryForm.ID, "wc -l", "13 seed_entry", 0)
var EviMetaWizardDance = newHPEviMetadata(EviWizardDance.ID, "color-averager", "rgb(22, 19, 20)\n#161314\nhsl(340, 7%, 8%)", 0)
var EviMetaPolyjuice = newHPEviMetadata(EviPolyjuice.ID, "wc -l", "13 seed_juice", 0)
var EviMetaWarewolf = newHPEviMetadata(EviWarewolf.ID, "duration", "348.163058815", 0)

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

var newHpFindingCategory = newFindingCategoryGen(1)

var ProductFindingCategory = newHpFindingCategory("Product", false)
var NetworkFindingCategory = newHpFindingCategory("Network", false)
var EnterpriseFindingCategory = newHpFindingCategory("Enterprise", false)
var VendorFindingCategory = newHpFindingCategory("Vendor", false)
var BehavioralFindingCategory = newHpFindingCategory("Behavioral", false)
var DetectionGapFindingCategory = newHpFindingCategory("Detection Gap", false)
var SomeFindingCategory = newHpFindingCategory("some-category", false)
var SomeOtherFindingCategory = newHpFindingCategory("alt-category", false)
var DeletedCategory = newHpFindingCategory("I was deleted", true)

var newHPFinding = newFindingGen(1)
var noLink = ""
var spiderLink = "https://www.google.com/search?q=spider+predators"
var FindingBook2Magic = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-b2magic", &SomeFindingCategory.ID, "lots o' magic", "Magic plagues Harry's life", nil)
var FindingBook2CGI = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-cgi", &SomeOtherFindingCategory.ID, "this looks fake", "I'm not entirely sure this is all above board", &noLink)
var FindingBook2SpiderFear = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-spider", &SomeFindingCategory.ID, "how to scare spiders", "Who would have thought?", &spiderLink)
var FindingBook2Robes = newHPFinding(OpChamberOfSecrets.ID, "find-uuid-robes", nil, "Robes for all seasons", "Turns out there's only one kind of robe.", &spiderLink)

var newHPServiceWorker = newServiceWorkerGen(1)
var demoServiceWorker = newHPServiceWorker("Demo", `{ "type": "web",  "version": 1, "url": "http://demo:3001/process" }`)
