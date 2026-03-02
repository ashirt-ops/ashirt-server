package services_test

import (
	"testing"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/services"
	"github.com/stretchr/testify/require"

	sq "github.com/Masterminds/squirrel"
)

func TestSetUserOperationRole(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserRon, db)

		masterOp := OpChamberOfSecrets
		targetUser := UserHarry
		targetRole := policy.OperationRoleRead
		input := services.SetUserOperationRoleInput{
			OperationSlug: masterOp.Slug,
			UserSlug:      targetUser.Slug,
			Role:          targetRole,
		}

		initialRole := seed.UserRoleForOp(targetUser, masterOp)
		require.NotContains(t, []policy.OperationRole{targetRole, ""}, initialRole, "Test user should both have a role, but not have the role we want to use")

		err := services.SetUserOperationRole(ctx, db, input)
		require.NoError(t, err)

		getDBRole := func() (string, error) {
			var newRole string
			err := db.Get(&newRole, sq.Select("role").
				From("user_operation_permissions").
				Where(sq.Eq{"operation_id": masterOp.ID, "user_id": targetUser.ID}))
			return newRole, err
		}
		newRole, err := getDBRole()
		require.NoError(t, err)
		require.Equal(t, string(targetRole), newRole)

		input = services.SetUserOperationRoleInput{
			OperationSlug: masterOp.Slug,
			UserSlug:      targetUser.Slug,
			Role:          "",
		}

		err = services.SetUserOperationRole(ctx, db, input)
		require.NoError(t, err)

		_, err = getDBRole()
		require.True(t, database.IsEmptyResultSetError(err))

		targetRole = policy.OperationRoleAdmin
		input = services.SetUserOperationRoleInput{
			OperationSlug: masterOp.Slug,
			UserSlug:      targetUser.Slug,
			Role:          targetRole,
		}
		err = services.SetUserOperationRole(ctx, db, input)
		require.NoError(t, err)

		newRole, err = getDBRole()
		require.NoError(t, err)
		require.Equal(t, string(targetRole), newRole)
	})
}

// write a test for SetUserGroupOperationRole
func TestSetUserGroupOperationRole(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		ctx := contextForUser(UserDumbledore, db)

		masterOp := OpSorcerersStone
		targetUserGroup := UserGroupSlytherin
		targetRole := policy.OperationRoleRead
		input := services.SetUserGroupOperationRoleInput{
			OperationSlug: masterOp.Slug,
			UserGroupSlug: targetUserGroup.Slug,
			Role:          targetRole,
		}

		initialRole := seed.UserGroupRoleForOp(targetUserGroup, masterOp)
		require.NotContains(t, []policy.OperationRole{targetRole, ""}, initialRole, "Test user group should have a role, but not have the role we want to use")

		err := services.SetUserGroupOperationRole(ctx, db, input)
		require.NoError(t, err)

		getDBRole := func() (string, error) {
			var newRole string
			err := db.Get(&newRole, sq.Select("role").
				From("user_group_operation_permissions").
				Where(sq.Eq{"operation_id": masterOp.ID, "group_id": targetUserGroup.ID}))
			return newRole, err
		}
		newRole, err := getDBRole()
		require.NoError(t, err)
		require.Equal(t, string(targetRole), newRole)

		input = services.SetUserGroupOperationRoleInput{
			OperationSlug: masterOp.Slug,
			UserGroupSlug: targetUserGroup.Slug,
			Role:          "",
		}

		err = services.SetUserGroupOperationRole(ctx, db, input)
		require.NoError(t, err)

		_, err = getDBRole()
		require.True(t, database.IsEmptyResultSetError(err))

		targetRole = policy.OperationRoleAdmin
		input = services.SetUserGroupOperationRoleInput{
			OperationSlug: masterOp.Slug,
			UserGroupSlug: targetUserGroup.Slug,
			Role:          targetRole,
		}
		err = services.SetUserGroupOperationRole(ctx, db, input)
		require.NoError(t, err)

		newRole, err = getDBRole()
		require.NoError(t, err)
		require.Equal(t, string(targetRole), newRole)
	})
}
