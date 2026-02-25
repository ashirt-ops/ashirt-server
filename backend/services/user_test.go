package services_test

import (
	"context"
	"testing"

	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/services"
	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"
)

type userValidator func(*testing.T, UserOpPermJoinUser, *dtos.UserOperationRole)

func TestCreateUser(t *testing.T) {
	RunDisposableDBTestWithSeed(t, NoSeedData, func(db *database.Connection, _ TestSeedData) {
		// verify first user is an admin
		i := services.CreateUserInput{
			FirstName: "Luna",
			LastName:  "Lovegood",
			Slug:      "luna.lovegood",
			Email:     "luna.lovegood@hogwarts.edu",
		}

		createUserOutput, err := services.CreateUser(context.Background(), db, i)
		require.NoError(t, err)
		require.Equal(t, createUserOutput.RealSlug, i.Slug)
		luna := getUserProfile(t, db, createUserOutput.UserID)

		require.Equal(t, true, luna.Admin)
		require.Equal(t, luna.FirstName, i.FirstName)
		require.Equal(t, luna.Email, i.Email)
		require.Equal(t, luna.LastName, i.LastName)

		// Verify re-register will fail (due to unique email constraint)
		_, err = services.CreateUser(context.Background(), db, i)
		require.Error(t, err)

		// Verify 2nd user (non-admin, no matching slug)
		i.Email = "luna.lovegood+extra@hogwarts.edu" // change the password to something that won't exist
		createUserOutput, err = services.CreateUser(context.Background(), db, i)
		require.NoError(t, err)
		// Since Luna's already exists, a new slug should be created
		require.NotEqual(t, i.Slug, createUserOutput.RealSlug)
		require.Contains(t, createUserOutput.RealSlug, i.Slug)
		newLuna := getUserProfile(t, db, createUserOutput.UserID)

		require.Equal(t, false, newLuna.Admin)
		require.Equal(t, i.FirstName, newLuna.FirstName)
		require.Equal(t, i.Email, newLuna.Email)
		require.Equal(t, i.LastName, newLuna.LastName)
	})
}

func TestCreateHeadlessUser(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		i := services.CreateUserInput{
			FirstName: "Extra",
			LastName:  "Headless Hunt Member",
			Slug:      "sir.nobody",
			Email:     "sir.nobody@hogwarts.edu",
		}

		// Verify non-admin can not create headless users
		ctx := contextForUser(UserHarry, db)
		_, err := services.CreateHeadlessUser(ctx, db, i)
		require.Error(t, err)

		ctx = contextForUser(UserDumbledore, db)
		result, err := services.CreateHeadlessUser(ctx, db, i)
		require.NoError(t, err)

		foundUser := getUserBySlug(t, db, result.RealSlug)

		require.True(t, foundUser.Headless)
	})
}

func TestDeleteUser(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		targetUser := UserRon
		admin := UserDumbledore

		require.True(t, 0 < countRows(t, db, "api_keys", "user_id=?", targetUser.ID))
		require.True(t, 0 < countRows(t, db, "auth_scheme_data", "user_id=?", targetUser.ID))
		require.True(t, 0 < countRows(t, db, "user_operation_permissions", "user_id=?", targetUser.ID))

		// verify that non-admins cannot delete
		ctx := contextForUser(UserDraco, db)
		err := services.DeleteUser(ctx, db, targetUser.Slug)
		require.Error(t, err)

		// verify user cannot delete themselves
		ctx = contextForUser(admin, db)
		err = services.DeleteUser(ctx, db, admin.Slug)
		require.NotNil(t, err)

		// Verify delete actually works
		err = services.DeleteUser(ctx, db, targetUser.Slug)
		require.Nil(t, err)

		require.True(t, countRows(t, db, "api_keys", "user_id=?", targetUser.ID) == 0)
		require.True(t, countRows(t, db, "auth_scheme_data", "user_id=?", targetUser.ID) == 0)
		require.True(t, countRows(t, db, "user_operation_permissions", "user_id=?", targetUser.ID) == 0)

		var user models.User
		err = db.Get(&user, sq.Select("*").From("users").Where(sq.Eq{"id": targetUser.ID}))
		require.Nil(t, err)
		require.NotNil(t, user.DeletedAt)
	})
}

func TestListEvidenceCreatorsForOperation(t *testing.T) {
	// TODO
}

func TestListUsersForAdmin(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		allUsers := getAllUsers(t, db)
		allDeletedUsers := getAllDeletedUsers(t, db)

		input := services.ListUsersForAdminInput{
			Pagination: services.Pagination{
				Page:     1,
				PageSize: 250,
			},
			IncludeDeleted: false,
		}
		input.Pagination.SetMaxItems(input.Pagination.PageSize) // force constrain not to affect us

		// verify access restricted for non-admins
		ctx := contextForUser(UserDraco, db)
		_, err := services.ListUsersForAdmin(ctx, db, input)
		require.Error(t, err)

		// Verify admins can list users (no deleted users)
		ctx = contextForUser(UserDumbledore, db)
		pagedUsers, err := services.ListUsersForAdmin(ctx, db, input)
		require.NoError(t, err)

		require.Equal(t, input.Pagination.PageSize, pagedUsers.PageSize)
		require.Equal(t, int64(len(allUsers)-len(allDeletedUsers)), pagedUsers.TotalCount)

		usersDto, ok := pagedUsers.Content.([]*dtos.UserAdminView)
		require.True(t, ok)
		dtoIndex := 0
		for i := 0; i < len(allUsers); i++ {
			if allUsers[i].DeletedAt != nil {
				continue
			}
			require.Equal(t, allUsers[i].Slug, usersDto[dtoIndex].Slug)
			require.Equal(t, allUsers[i].FirstName, usersDto[dtoIndex].FirstName)
			require.Equal(t, allUsers[i].LastName, usersDto[dtoIndex].LastName)
			require.Equal(t, allUsers[i].Admin, usersDto[dtoIndex].Admin)
			require.Equal(t, allUsers[i].Disabled, usersDto[dtoIndex].Disabled)
			require.Equal(t, allUsers[i].Email, usersDto[dtoIndex].Email)
			require.Equal(t, allUsers[i].Headless, usersDto[dtoIndex].Headless)
			require.Equal(t, false, usersDto[dtoIndex].Deleted)
			dtoIndex++
		}

		// verify deleted users can be shown
		input.IncludeDeleted = true
		pagedUsers, err = services.ListUsersForAdmin(ctx, db, input)
		require.Nil(t, err)

		usersDto, _ = pagedUsers.Content.([]*dtos.UserAdminView)
		for i := 0; i < len(allUsers); i++ {
			require.Equal(t, allUsers[i].Slug, usersDto[i].Slug)
			require.Equal(t, allUsers[i].FirstName, usersDto[i].FirstName)
			require.Equal(t, allUsers[i].LastName, usersDto[i].LastName)
			require.Equal(t, allUsers[i].Admin, usersDto[i].Admin)
			require.Equal(t, allUsers[i].Disabled, usersDto[i].Disabled)
			require.Equal(t, allUsers[i].Email, usersDto[i].Email)
			require.Equal(t, allUsers[i].Headless, usersDto[i].Headless)
			require.Equal(t, (allUsers[i].DeletedAt != nil), usersDto[i].Deleted)
		}
	})
}

func TestListUsersForOperation(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		ctx := contextForUser(UserRon, db)

		masterOp := OpChamberOfSecrets
		allUserOpRoles := getUsersWithRoleForOperationByOperationID(t, db, masterOp.ID)
		require.NotEqual(t, len(allUserOpRoles), 0, "Some users should be attached to this operation")

		input := services.ListUsersForOperationInput{
			OperationSlug: masterOp.Slug,
		}

		content, err := services.ListUsersForOperation(ctx, db, input)
		require.NoError(t, err)

		require.Equal(t, len(content), len(allUserOpRoles))
		validateUserSets(t, content, allUserOpRoles, validateUser)
	})
}

func TestListUsers(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		testListUsersCase(t, db, "harry potter", true, []models.User{UserHarry})
		testListUsersCase(t, db, "granger", true, []models.User{UserHermione})
		testListUsersCase(t, db, "al", true, []models.User{UserAlastor, UserDumbledore, UserDraco, UserLucius, UserMinerva})
		testListUsersCase(t, db, "dra mal", true, []models.User{UserDraco})
		testListUsersCase(t, db, "", true, []models.User{})
		testListUsersCase(t, db, "  ", true, []models.User{})
		testListUsersCase(t, db, "%", true, []models.User{})
		testListUsersCase(t, db, "*", true, []models.User{})
		testListUsersCase(t, db, "___", true, []models.User{})

		// test for deleted user filtering
		testListUsersCase(t, db, UserTomRiddle.LastName, true, []models.User{UserTomRiddle})
		testListUsersCase(t, db, UserTomRiddle.LastName, false, []models.User{})
	})
}

func testListUsersCase(t *testing.T, db *database.Connection, query string, includeDeleted bool, expectedUsers []models.User) {
	ctx := contextForUser(UserHarry, db)

	users, err := services.ListUsers(ctx, db, services.ListUsersInput{Query: query, IncludeDeleted: includeDeleted})
	require.NoError(t, err)

	require.Equal(t, len(expectedUsers), len(users), "Expected %d users for query '%s' but got %d", len(expectedUsers), query, len(users))

	for i := range expectedUsers {
		require.Equal(t, expectedUsers[i].Slug, users[i].Slug)
		require.Equal(t, expectedUsers[i].FirstName, users[i].FirstName)
		require.Equal(t, expectedUsers[i].LastName, users[i].LastName)
	}
}

func validateUser(t *testing.T, expected UserOpPermJoinUser, actual *dtos.UserOperationRole) {
	require.Equal(t, expected.Slug, actual.User.Slug)
	require.Equal(t, expected.FirstName, actual.User.FirstName)
	require.Equal(t, expected.LastName, actual.User.LastName)
	require.Equal(t, expected.Role, actual.Role)
}

func validateUserSets(t *testing.T, dtoSet []*dtos.UserOperationRole, dbSet []UserOpPermJoinUser, validate userValidator) {
	var expected *UserOpPermJoinUser = nil

	for _, dtoItem := range dtoSet {
		expected = nil
		for _, dbItem := range dbSet {
			if dbItem.Slug == dtoItem.User.Slug {
				expected = &dbItem
				break
			}
		}
		require.NotNil(t, expected, "Result should have matching value")
		validate(t, *expected, dtoItem)
	}
}

func TestReadUser(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserRon
		targetUser := UserHarry
		adminUser := UserDumbledore
		ctx := contextForUser(normalUser, db)

		supportedAuthSchemes := []dtos.SupportedAuthScheme{
			{SchemeName: "Local", SchemeCode: "local"},
		}

		// verify read-self
		retrievedUser, err := services.ReadUser(ctx, db, "", &supportedAuthSchemes)
		require.NoError(t, err)
		verifyRetrievedUser(t, normalUser, retrievedUser, supportedAuthSchemes)

		// verify read-self alternative (userslug provided)
		retrievedUser, err = services.ReadUser(ctx, db, normalUser.Slug, &supportedAuthSchemes)
		require.NoError(t, err)
		verifyRetrievedUser(t, normalUser, retrievedUser, supportedAuthSchemes)

		// verify read-other (non-admin : should fail)
		_, err = services.ReadUser(ctx, db, targetUser.Slug, &supportedAuthSchemes)
		require.Error(t, err)

		// verify read-other (as admin)
		ctx = contextForUser(adminUser, db)
		retrievedUser, err = services.ReadUser(ctx, db, targetUser.Slug, &supportedAuthSchemes)
		require.NoError(t, err)
		verifyRetrievedUser(t, targetUser, retrievedUser, supportedAuthSchemes)

		// verify old/removed auth schemes are filtered out
		ctx = contextForUser(normalUser, db)
		supportedAuthSchemes = []dtos.SupportedAuthScheme{
			{SchemeName: "Petronus", SchemeCode: "petroni"},
		}
		retrievedUser, err = services.ReadUser(ctx, db, "", &supportedAuthSchemes)
		require.NoError(t, err)
		verifyRetrievedUser(t, normalUser, retrievedUser, []dtos.SupportedAuthScheme{})
	})
}

func verifyRetrievedUser(t *testing.T, expectedUser models.User, retrievedUser *dtos.UserOwnView, expectedAuths []dtos.SupportedAuthScheme) {
	require.Equal(t, expectedUser.Slug, retrievedUser.Slug)
	require.Equal(t, expectedUser.FirstName, retrievedUser.FirstName)
	require.Equal(t, expectedUser.LastName, retrievedUser.LastName)
	require.Equal(t, expectedUser.Email, retrievedUser.Email)
	for _, expectedAuth := range expectedAuths {
		found := false

		for _, returnedAuth := range retrievedUser.Authentication {
			if expectedAuth.SchemeCode == returnedAuth.AuthSchemeCode {
				found = true
				break
			}
		}
		require.True(t, found)
	}
}

func TestUpdateUserProfile(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		normalUser := UserRon
		targetUser := UserHarry
		adminUser := UserDumbledore
		ctx := contextForUser(normalUser, db)

		// verify read-self
		verifyUserProfileUpdate(t, false, ctx, db, normalUser.ID, services.UpdateUserProfileInput{
			FirstName: "Stan",
			LastName:  "Shunpike",
			Email:     "sshunpike@hogwarts.edu",
		})

		// verify read-self (alternate)
		verifyUserProfileUpdate(t, false, ctx, db, normalUser.ID, services.UpdateUserProfileInput{
			UserSlug:  normalUser.Slug,
			FirstName: "Stan2",
			LastName:  "Shunpike2",
			Email:     "sshunpike2@hogwarts.edu",
		})

		// verify read-other (non-admin)
		verifyUserProfileUpdate(t, true, ctx, db, targetUser.ID, services.UpdateUserProfileInput{
			UserSlug:  targetUser.Slug,
			FirstName: "Stan3",
			LastName:  "Shunpike3",
			Email:     "sshunpike3@hogwarts.edu",
		})

		// verify read-other (admin)
		ctx = contextForUser(adminUser, db)
		verifyUserProfileUpdate(t, false, ctx, db, targetUser.ID, services.UpdateUserProfileInput{
			UserSlug:  targetUser.Slug,
			FirstName: "Stan4",
			LastName:  "Shunpike4",
			Email:     "sshunpike4@hogwarts.edu",
		})
	})
}

func verifyUserProfileUpdate(t *testing.T, expectError bool, ctx context.Context, db *database.Connection, userID int64, updatedData services.UpdateUserProfileInput) {
	err := services.UpdateUserProfile(ctx, db, updatedData)
	if expectError {
		require.NotNil(t, err)
		return
	}

	require.NoError(t, err)

	newProfile := getUserProfile(t, db, userID)
	require.NoError(t, err)
	require.Equal(t, updatedData.FirstName, newProfile.FirstName)
	require.Equal(t, updatedData.LastName, newProfile.LastName)
	require.Equal(t, updatedData.Email, newProfile.Email)
}

func TestDeleteSessionsForUserSlug(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		targetedUser := UserDraco
		alsoPresentUser := UserHarry

		// populate some sessions
		sessionsToAdd := []models.Session{
			{UserID: targetedUser.ID, SessionData: []byte("a")},
			{UserID: alsoPresentUser.ID, SessionData: []byte("b")},
			{UserID: targetedUser.ID, SessionData: []byte("c")},
			{UserID: alsoPresentUser.ID, SessionData: []byte("d")},
		}
		err := db.BatchInsert("sessions", len(sessionsToAdd), func(i int) map[string]interface{} {
			return map[string]interface{}{
				"user_id":      sessionsToAdd[i].UserID,
				"session_data": sessionsToAdd[i].SessionData,
			}
		})

		require.NoError(t, err)

		// verify sessions exist
		var targetedUserSessions []models.Session
		err = db.Select(&targetedUserSessions, sq.Select("*").From("sessions").Where(sq.Eq{"user_id": targetedUser.ID}))
		require.NoError(t, err)
		require.True(t, len(targetedUserSessions) > 0)

		// verify non-admin cannot delete session data
		ctx := contextForUser(UserHarry, db)
		err = services.DeleteSessionsForUserSlug(ctx, db, targetedUser.Slug)
		require.Error(t, err)

		// verify admin can delete session data
		ctx = contextForUser(UserDumbledore, db)
		err = services.DeleteSessionsForUserSlug(ctx, db, targetedUser.Slug)
		require.NoError(t, err)

		targetedUserSessions = []models.Session{}
		err = db.Select(&targetedUserSessions, sq.Select("*").From("sessions").Where(sq.Eq{"user_id": targetedUser.ID}))
		require.NoError(t, err)
		require.True(t, len(targetedUserSessions) == 0)
	})
}

func TestSetUserFlags(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		targetUser := UserHarry
		adminUser := UserDumbledore
		admin := true
		disabled := true
		input := services.SetUserFlagsInput{
			Slug:     targetUser.Slug,
			Admin:    &admin,
			Disabled: &disabled,
		}

		// verify access restricted for non-admins
		ctx := contextForUser(UserDraco, db)
		err := services.SetUserFlags(ctx, db, input)
		require.Error(t, err)

		// As an admin
		ctx = contextForUser(adminUser, db)

		// verify users can't disable themselves
		sameUserInput := services.SetUserFlagsInput{
			Slug:     adminUser.Slug,
			Admin:    &admin,    // true at this point (no change)
			Disabled: &disabled, // true at this point
		}
		err = services.SetUserFlags(ctx, db, sameUserInput)
		require.Error(t, err)

		// verify users can't demote themselves
		disabled = false
		admin = false
		err = services.SetUserFlags(ctx, db, sameUserInput)
		require.Error(t, err)

		// reset for next tests
		disabled = true
		admin = true

		// try setting and then unsetting admin/disabled
		for i := 0; i < 2; i++ {
			err = services.SetUserFlags(ctx, db, input)
			require.NoError(t, err)

			dbProfile := getUserProfile(t, db, targetUser.ID)

			require.Equal(t, admin, dbProfile.Admin)
			require.Equal(t, disabled, dbProfile.Disabled)

			// second test: Make sure setting to false also works
			admin = !admin
			disabled = !disabled
		}

		// verify headless users cannot be admins
		admin = true
		err = services.SetUserFlags(ctx, db, services.SetUserFlagsInput{
			Slug:  UserHeadlessNick.Slug,
			Admin: &admin,
		})
		require.Error(t, err)
	})
}
