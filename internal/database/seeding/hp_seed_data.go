package seeding

import (
	"strings"
	"time"

	"github.com/ashirt-ops/ashirt-server/internal/helpers"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/servicetypes/evidencemetadata"
)

var HarryPotterSeedData = Seeder{
	FindingCategories: []models.FindingCategory{
		ProductFindingCategory, NetworkFindingCategory, EnterpriseFindingCategory, VendorFindingCategory, BehavioralFindingCategory, DetectionGapFindingCategory,
		DeletedCategory, SomeFindingCategory, SomeOtherFindingCategory,
	},
	Users: []models.User{UserHarry, UserRon, UserGinny, UserHermione, UserNeville, UserSeamus, UserDraco, UserSnape, UserDumbledore, UserHagrid, UserTomRiddle, UserHeadlessNick,
		UserCedric, UserFleur, UserViktor, UserAlastor, UserMinerva, UserLucius, UserSirius, UserPeter, UserParvati, UserPadma, UserCho,
	},
	UserGroups: []models.UserGroup{
		UserGroupGryffindor, UserGroupHufflepuff, UserGroupRavenclaw, UserGroupSlytherin, UserGroupOtherHouse,
	},
	UserGroupMaps: []models.UserGroupMap{
		AddHarryToGryffindor, AddRonToGryffindor, AddGinnyToGryffindor, AddHermioneToGryffindor,
		AddMalfoyToSlytherin, AddSnapeToSlytherin, AddLuciusToSlytherin,
		AddCedricToHufflepuff, AddFleurToHufflepuff,
		AddChoToRavenclaw, AddViktorToRavenclaw,
	},
	Operations: []models.Operation{OpSorcerersStone, OpChamberOfSecrets, OpPrisonerOfAzkaban, OpGobletOfFire, OpOrderOfThePhoenix, OpHalfBloodPrince, OpDeathlyHallows},
	Tags: []models.Tag{
		TagFamily, TagFriendship, TagHome, TagLoyalty, TagCourage, TagGoodVsEvil, TagSupernatural,
		TagMercury, TagVenus, TagEarth, TagMars, TagJupiter, TagSaturn, TagNeptune,

		// common tags among all operations
		CommonTagWhoSS, CommonTagWhatSS, CommonTagWhereSS, CommonTagWhenSS, CommonTagWhySS,
		CommonTagWhoCoS, CommonTagWhatCoS, CommonTagWhereCoS, CommonTagWhenCoS, CommonTagWhyCoS,
		CommonTagWhoGoF, CommonTagWhatGoF, CommonTagWhereGoF, CommonTagWhenGoF, CommonTagWhyGoF,
	},
	DefaultTags: []models.DefaultTag{
		DefaultTagWho, DefaultTagWhat, DefaultTagWhere, DefaultTagWhen, DefaultTagWhy,
	},
	APIKeys: []models.APIKey{
		APIKeyHarry1, APIKeyHarry2,
		APIKeyRon1, APIKeyRon2,
		APIKeyNick,
	},
	UserOpPrefMap: []models.UserOperationPreferences{
		newUserOperationPreferences(UserRon, OpChamberOfSecrets, true),
		newUserOperationPreferences(UserDumbledore, OpGobletOfFire, true),
		// This user doesn't have permission to view this operation. This helps check that it's okay
		// to not clean this data up. (and other checks should verify the users don't need to have
		// not-a-favorite set up.)
		newUserOperationPreferences(UserDraco, OpChamberOfSecrets, true),
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
	},
	UserGroupOpMap: []models.UserGroupOperationPermission{
		newUserGroupOpPermission(UserGroupGryffindor, OpSorcerersStone, policy.OperationRoleRead),
		newUserGroupOpPermission(UserGroupHufflepuff, OpSorcerersStone, policy.OperationRoleWrite),
		newUserGroupOpPermission(UserGroupSlytherin, OpSorcerersStone, policy.OperationRoleAdmin),
	},
	Findings: []models.Finding{
		FindingBook2Magic, FindingBook2CGI, FindingBook2SpiderFear, FindingBook2Robes,
	},
	Evidences: []models.Evidence{
		EviDursleys, EviMirrorOfErised, EviLevitateSpell, EviRulesForQuidditch,
		EviFlyingCar, EviDobby, EviSpiderAragog, EviMoaningMyrtle, EviWhompingWillow, EviTomRiddlesDiary, EviPetrifiedHermione, EviLazyHar, EviHeadlessHuntApplication,
		EviTristateTrophy, EviEntryForm, EviWizardDance, EviPolyjuice, EviWarewolf,
	},
	EvidenceMetadatas: []models.EvidenceMetadata{
		EviMetaDursleys, EviMetaMirrorOfErised, EviMetaLevitateSpell, EviMetaRulesForQuidditch, EviMetaRulesForQuidditchTwo,
		EviMetaFlyingCar, EviMetaDobby, EviMetaSpiderAragog, EviMetaMoaningMyrtle, EviMetaWhompingWillow, EviMetaTomRiddlesDiary, EviMetaTomRiddlesDiaryTwo, EviMetaPetrifiedHermione, EviMetaHeadlessHuntApplication,
		EviMetaTristateTrophy, EviMetaEntryForm, EviMetaWizardDance, EviMetaPolyjuice, EviMetaEntryFormTwo, EviMetaWarewolf, EviMetaWarewolfOther,
	},
	TagEviMap: unionTagEviMap(
		associateTagsToEvidence(EviDursleys, TagFamily, TagHome),
		associateTagsToEvidence(EviFlyingCar, TagEarth, TagSaturn),
		associateTagsToEvidence(EviDobby, TagMars, TagJupiter, TagMercury),
		associateTagsToEvidence(EviPetrifiedHermione, TagMars, CommonTagWhatCoS, CommonTagWhoCoS),
		associateTagsToEvidence(EviLazyHar, CommonTagWhatCoS),

		associateTagsToEvidence(EviTristateTrophy, CommonTagWhoGoF, CommonTagWhereGoF, CommonTagWhyGoF),
		associateTagsToEvidence(EviEntryForm, CommonTagWhatGoF, CommonTagWhereGoF, CommonTagWhenGoF),
		associateTagsToEvidence(EviWizardDance, CommonTagWhereGoF, CommonTagWhenGoF),
		associateTagsToEvidence(EviPolyjuice, CommonTagWhatGoF, CommonTagWhereGoF, CommonTagWhenGoF, CommonTagWhyGoF),
		associateTagsToEvidence(EviWarewolf, CommonTagWhoGoF, CommonTagWhereGoF, CommonTagWhyGoF),
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
		DemoServiceWorker,
	},
	GlobalVars: []models.GlobalVar{
		VarExpelliarmus, VarAlohomora, VarAscendio, VarImperio, VarLumos, VarObliviate,
	},
	OperationVars: []models.OperationVar{
		OpVarImmobulus, OpVarObscuro, OpVarProtego, OpVarReparo, OpVarStupefy, OpVarWingardiumLeviosa,
	},
	VarOperationMap: unionVarOperationMap(
		associateVarsToOperation(OpSorcerersStone, OpVarImmobulus, OpVarObscuro),
		associateVarsToOperation(OpChamberOfSecrets, OpVarProtego, OpVarReparo),
		associateVarsToOperation(OpGobletOfFire, OpVarStupefy, OpVarWingardiumLeviosa),
	),
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

var newUserGroup = newUserGroupGen(1)

var UserGroupGryffindor = newUserGroup("Gryffindor", false)
var UserGroupHufflepuff = newUserGroup("Hufflepuff", false)
var UserGroupRavenclaw = newUserGroup("Ravenclaw", false)
var UserGroupSlytherin = newUserGroup("Slytherin", false)

// UserGroupOtherHouse is reserved to test deleted user groups
var UserGroupOtherHouse = newUserGroup("Other House", true)

var AddHarryToGryffindor = newUserGroupMapping(UserHarry, UserGroupGryffindor)
var AddRonToGryffindor = newUserGroupMapping(UserRon, UserGroupGryffindor)
var AddGinnyToGryffindor = newUserGroupMapping(UserGinny, UserGroupGryffindor)
var AddHermioneToGryffindor = newUserGroupMapping(UserHermione, UserGroupGryffindor)

var AddMalfoyToSlytherin = newUserGroupMapping(UserDraco, UserGroupSlytherin)
var AddLuciusToSlytherin = newUserGroupMapping(UserLucius, UserGroupSlytherin)
var AddSnapeToSlytherin = newUserGroupMapping(UserSnape, UserGroupSlytherin)

var AddCedricToHufflepuff = newUserGroupMapping(UserCedric, UserGroupHufflepuff)
var AddFleurToHufflepuff = newUserGroupMapping(UserFleur, UserGroupHufflepuff)

var AddViktorToRavenclaw = newUserGroupMapping(UserViktor, UserGroupRavenclaw)
var AddChoToRavenclaw = newUserGroupMapping(UserCho, UserGroupRavenclaw)

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

var timeNow = time.Now()

var newHPEvidence = newEvidenceGen(1)
var EviDursleys = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_dursleys", "Family of self-centered muggles + Harry", "image", 0, &timeNow)
var EviMirrorOfErised = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_mirror", "Mysterious mirror that shows you your deepest desires", "image", 0, nil)
var EviLevitateSpell = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_md_levitate", "Documented Levitation Spell effects", "codeblock", 0, &timeNow)
var EviRulesForQuidditch = newHPEvidence(OpSorcerersStone.ID, UserHarry.ID, "seed_rs_aoc201501", "Complex rules for a simple game", "codeblock", 0, &timeNow)

var newHPEviMetadata = newEvidenceMetadataGen(1)
var EviMetaDursleys = newHPEviMetadata(EviDursleys.ID, "color-averager", "rgb(65, 65, 65)\n#414141\nhsl(0, 0%, 25%)", nil, helpers.PTrue(), 0)
var EviMetaMirrorOfErised = newHPEviMetadata(EviMirrorOfErised.ID, "color-averager", "rgb(111, 77, 14)\n#6f4d0e\nhsl(39, 78%, 25%)", nil, helpers.PTrue(), 0)
var EviMetaLevitateSpell = newHPEviMetadata(EviLevitateSpell.ID, "wc -l", "12 seed_md_levitate", nil, helpers.PTrue(), 0)
var EviMetaRulesForQuidditch = newHPEviMetadata(EviRulesForQuidditch.ID, "run-result", "Last floor reached: 138 \n Step on which basement is reached (first time): 1771", nil, helpers.PTrue(), 0)
var EviMetaRulesForQuidditchTwo = newHPEviMetadata(EviRulesForQuidditch.ID, "wc-l", "33 main.rs", nil, helpers.PTrue(), 0)

var EviFlyingCar = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_car", "A Car that flies", "image", 0, &timeNow)
var EviDobby = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_dobby", "an elf?", "image", 0, nil)
var EviSpiderAragog = newHPEvidence(OpChamberOfSecrets.ID, UserHagrid.ID, "seed_aragog", "Just a big spider", "image", 0, nil)
var EviMoaningMyrtle = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_myrtle", "She's very sad", "image", 0, nil)
var EviWhompingWillow = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_tree", "Don't get too close", "image", 0, &timeNow)
var EviTomRiddlesDiary = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_go_aoc201614", "What's a Horcrux?", "codeblock", 0, &timeNow)
var EviHeadlessHuntApplication = newHPEvidence(OpChamberOfSecrets.ID, UserRon.ID, "seed_py_aoc201717", "This group is very particular", "codeblock", 0, &timeNow)
var EviPetrifiedHermione = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_statue", "Strangely real-looking statue", "image", 0, &timeNow)
var EviLazyHar = newHPEvidence(OpChamberOfSecrets.ID, UserHarry.ID, "seed_har", "Joel couldn't be bothered to add a themed asset.", "http-request-cycle", 0, &timeNow)

var EviMetaFlyingCar = newHPEviMetadata(EviFlyingCar.ID, "color-averager", "rgb(106, 109, 84)\n#6a6d54\nhsl(67, 13%, 38%)", evidencemetadata.StatusCompleted.Ptr(), helpers.PTrue(), 0)
var EviMetaDobby = newHPEviMetadata(EviDobby.ID, "color-averager", "rgb(74, 51, 32)\n#4a3320\nhsl(27, 40%, 21%)", nil, helpers.PTrue(), 0)
var EviMetaSpiderAragog = newHPEviMetadata(EviSpiderAragog.ID, "color-averager", "rgb(189, 156, 146)\n#bd9c92\nhsl(14, 25%, 66%)", nil, helpers.PTrue(), 0)
var EviMetaMoaningMyrtle = newHPEviMetadata(EviMoaningMyrtle.ID, "color-averager", "rgb(118, 103, 102)\n#766766\nhsl(4, 7%, 43%)", nil, helpers.PTrue(), 0)
var EviMetaWhompingWillow = newHPEviMetadata(EviWhompingWillow.ID, "color-averager", "rgb(115, 109, 81)\n#736d51\nhsl(49, 17%, 38%)", nil, helpers.PTrue(), 0)
var EviMetaTomRiddlesDiary = newHPEviMetadata(EviTomRiddlesDiary.ID, "run-result", "All keys found by index:  19968", nil, helpers.PTrue(), 0)
var EviMetaTomRiddlesDiaryTwo = newHPEviMetadata(EviTomRiddlesDiary.ID, "wc -l", "98 main.go", nil, helpers.PTrue(), 0)
var EviMetaHeadlessHuntApplication = newHPEviMetadata(EviHeadlessHuntApplication.ID, "run-result", "41797835\nelapsed time (seconds): 3.772843360900879", nil, helpers.PTrue(), 0)
var EviMetaPetrifiedHermione = newHPEviMetadata(EviPetrifiedHermione.ID, "color-averager", "rgb(162, 104, 101)\n#a26865\nhsl(3, 25%, 52%)", nil, helpers.PTrue(), 0)

var EviTristateTrophy = newHPEvidence(OpGobletOfFire.ID, UserHarry.ID, "seed_trophy", "First Triwizard Champion Trophy", "image", 0, nil)
var EviEntryForm = newHPEvidence(OpGobletOfFire.ID, UserCedric.ID, "seed_entry", "Cedric's entry form for Triwizard competition", "codeblock", 0, nil)
var EviWizardDance = newHPEvidence(OpGobletOfFire.ID, UserCho.ID, "seed_dance", "Advertising for the Triwizard Dance", "image", 0, &timeNow)
var EviPolyjuice = newHPEvidence(OpGobletOfFire.ID, UserAlastor.ID, "seed_juice", "DIY instructions for Polyjuice Potion", "codeblock", 0, &timeNow)
var EviWarewolf = newHPEvidence(OpGobletOfFire.ID, UserViktor.ID, "seed_wolf", "Strangely real-looking statue", "terminal-recording", 0, &timeNow)

var EviMetaTristateTrophy = newHPEviMetadata(EviTristateTrophy.ID, "color-averager", "rgb(182, 184, 183)\n#b6b8b7\nhsl(150, 1%, 72%)", nil, helpers.PTrue(), 0)
var EviMetaEntryForm = newHPEviMetadata(EviEntryForm.ID, "run-result", "(No output)", nil, helpers.PTrue(), 0)
var EviMetaEntryFormTwo = newHPEviMetadata(EviEntryForm.ID, "wc -l", "13 seed_entry", nil, helpers.PTrue(), 0)
var EviMetaWizardDance = newHPEviMetadata(EviWizardDance.ID, "color-averager", "rgb(22, 19, 2nil, 0)\n#161314\nhsl(340, 7%, 8%)", nil, helpers.PTrue(), 0)
var EviMetaPolyjuice = newHPEviMetadata(EviPolyjuice.ID, "wc -l", "13 seed_juice", nil, helpers.PTrue(), 0)
var EviMetaWarewolf = newHPEviMetadata(EviWarewolf.ID, "duration", "348.163058815", nil, helpers.PTrue(), 0)
var EviMetaWarewolfOther = newHPEviMetadata(EviWarewolf.ID, "color-averager", "Yo, I can't process this", nil, helpers.PFalse(), 0)

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
var DemoServiceWorker = newHPServiceWorker("Demo", `{ "type": "web",  "version": 1, "url": "http://demo:8080/process" }`)

var newGlobalVar = newGlobalVarGen(1)
var VarExpelliarmus = newGlobalVar("EXPELLIARMUS", "disarm an opponent")
var VarAlohomora = newGlobalVar("ALOHOMORA", "unlock doors")
var VarAscendio = newGlobalVar("ASCENDIO", "lifts the caster high into the air")
var VarImperio = newGlobalVar("IMPERIO", "control another person")
var VarLumos = newGlobalVar("LUMOS", "creates a narrow beam of light")
var VarObliviate = newGlobalVar("PETRIFICUS_TOTALUS", "paralyzes someone")

var newOperationVar = newOperationVarGen(1)
var OpVarImmobulus = newOperationVar("immobulus", "IMMOBULUS", "freezes objects")
var OpVarObscuro = newOperationVar("obscuro", "OBSCURO", "blindfolds the victim")
var OpVarProtego = newOperationVar("protego", "PROTEGO", "shield charm")
var OpVarReparo = newOperationVar("reparo", "REPARO", "repairs broken objects")
var OpVarStupefy = newOperationVar("stupefy", "STUPEFY", "knocks out opponent")
var OpVarWingardiumLeviosa = newOperationVar("wingardium_leviosa", "WINGARDIUM_LEVIOSA", "levitates objects")
