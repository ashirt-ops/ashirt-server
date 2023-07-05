// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/alexedwards/scs/v2"
	"github.com/theparanoids/ashirt-server/backend"
	localauth "github.com/theparanoids/ashirt-server/backend/authschemes/localauth/constants"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateUserInput struct {
	FirstName string
	LastName  string
	Slug      string
	Email     string
	Headless  bool
}

type ListEvidenceCreatorsForOperationInput struct {
	OperationSlug string
}

type ListUsersForAdminInput struct {
	UserFilter
	Pagination
	IncludeDeleted bool
}

type ListUsersForOperationInput struct {
	Pagination
	UserFilter
	OperationSlug string
}

type userAndRole struct {
	models.User
	Role policy.OperationRole `db:"role"`
}

type ListUsersInput struct {
	Query          string
	IncludeDeleted bool
}

type UpdateUserProfileInput struct {
	UserSlug  string
	FirstName string
	LastName  string
	Email     string
}

type SetUserFlagsInput struct {
	Slug     string
	Disabled *bool
	Admin    *bool
}

func (cui CreateUserInput) validate() error {
	if cui.Slug == "" {
		return backend.MissingValueErr("User Slug")
	}
	if cui.FirstName == "" {
		return backend.MissingValueErr("First Name")
	}
	if cui.LastName == "" {
		return backend.MissingValueErr("Last Name")
	}
	if cui.Email == "" {
		return backend.MissingValueErr("Email")
	}
	return nil
}

// CreateHeadlessUser is really just CreateUser. The difference here is that _headless_ users will not have
// authentication, and instead rely on user-impersonation and API keys for access.
func CreateHeadlessUser(ctx context.Context, db *database.Connection, i CreateUserInput) (*dtos.CreateUserOutput, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unable to create new headless user", backend.UnauthorizedWriteErr(err))
	}
	i.Headless = true
	return CreateUser(db, i)
}

// CreateUser generates an entry in the users table in the database. No more is done here, but it is expected
// that the caller will, at a minimum, also want to create an entry in the authentication tables, so
// that the user can actually log in.
//
// Note: CreateUserInput.Slug is a _suggestion_, and it may be altered to ensure uniqueness.
//
// Returns a structure containing both the true slug (i.e. what it was mangled to, if it was infact mangled), plus
// the associated user_id value
func CreateUser(db *database.Connection, i CreateUserInput) (*dtos.CreateUserOutput, error) {
	validationErr := i.validate()
	if validationErr != nil {
		return nil, backend.WrapError("Unable to create new user", validationErr)
	}

	var userID int64
	var err error
	slugSuffix := ""
	var attemptedSlug string
	attemptNumber := 1
	for {
		attemptedSlug = i.Slug + slugSuffix
		userID, err = db.Insert("users", map[string]interface{}{
			"slug":       attemptedSlug,
			"first_name": i.FirstName,
			"last_name":  i.LastName,
			"email":      i.Email,
			"headless":   i.Headless,
		})
		if err != nil {
			if database.IsAlreadyExistsError(err) {
				if strings.Contains(err.Error(), "users.unique_email") { // not sure how else to check if this is a duplicate slug vs a duplicate email address
					return nil, backend.WrapError("Unable to insert new user", backend.DatabaseErr(err))
				}

				if attemptNumber > 5 {
					return nil, backend.WrapError("Unable to create new user after many attempts", backend.DatabaseErr(err))
				}

				logging.GetSystemLogger().Log(
					"msg", "Unable to create user with slug; trying alternative",
					"slug", attemptedSlug,
					"attempt", attemptNumber,
					"error", err.Error(),
				)
				attemptNumber++

				// an account with this slug already exists, attempt creating it again with a suffix
				// TODO: There's a possible, but impractical infinite loop here. We need some way to escape this
				slugSuffix = fmt.Sprintf("-%d", rand.Intn(99999))
				continue
			}
			return nil, backend.WrapError("Unable to insert new user", backend.DatabaseErr(err))
		}
		break
	}
	if userID == 1 {
		err := db.Update(sq.Update("users").Set("admin", true).Where(sq.Eq{"id": userID}))
		if err != nil {
			logging.GetSystemLogger().Log("msg", "Unable to make the first user an admin", "error", err.Error())
		}
	}
	return &dtos.CreateUserOutput{
		RealSlug: attemptedSlug,
		UserID:   userID,
	}, nil
}

// DeleteUser provides the ability for a super admin to remove a user from the system. Doing so
// removes access only. Evidence and other contributions remain. Note that users are not able to
// delete their own accounts to prevent accidents. Also note that once a user has been deleted,
// they cannot be restored.
func DeleteUser(ctx context.Context, sessionManager *scs.SessionManager, db *database.Connection, slug string) error {
	if !middleware.IsAdmin(ctx) {
		return backend.WrapError("Unwilling to delete user", backend.UnauthorizedWriteErr(fmt.Errorf("Requesting user is not an admin")))
	}

	userID, err := userSlugToUserID(db, slug)
	if err != nil {
		return backend.WrapError("Unable to delete user", backend.DatabaseErr(err))
	}

	if userID == middleware.UserID(ctx) {
		return backend.BadInputErr(fmt.Errorf("User is trying to delete themself"), "Users cannot delete themselves")
	}

	disabled := true
	// session data is deleted when disabling the user
	disableErr := SetUserFlags(ctx, sessionManager, db, SetUserFlagsInput{
		Slug:     slug,
		Disabled: &disabled,
	})
	if disableErr != nil {
		return backend.WrapError("Could not set user to disabled prior to deletion", disableErr)
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		tx.Delete(sq.Delete("api_keys").Where(sq.Eq{"user_id": userID}))
		tx.Delete(sq.Delete("auth_scheme_data").Where(sq.Eq{"user_id": userID}))
		tx.Delete(sq.Delete("user_operation_permissions").Where(sq.Eq{"user_id": userID}))
		tx.Update(sq.Update("users").Set("deleted_at", time.Now()).Where(sq.Eq{"slug": slug}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete user", backend.DatabaseErr(err))
	}

	return nil
}

// ListEvidenceCreatorsForOperation returns a list of all users that have (ever) created a piece of
// evidence for a given operation slug. Note that this won't return users that _had_ created evidence
// that has since been deleted
func ListEvidenceCreatorsForOperation(ctx context.Context, db *database.Connection, i ListEvidenceCreatorsForOperationInput) ([]*dtos.User, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list evidence for an operation", backend.UnauthorizedReadErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list evidence for an operation", backend.UnauthorizedReadErr(err))
	}

	var users []struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Slug      string `db:"slug"`
	}

	sb := sq.Select("users.slug", "users.first_name", "users.last_name").
		Distinct().
		From("operations").
		LeftJoin("evidence ON operations.id = evidence.operation_id").
		InnerJoin("users ON evidence.operator_id = users.id").
		Where(sq.Eq{"operations.slug": i.OperationSlug}).
		OrderBy("users.first_name ASC")

	err = db.Select(&users, sb)
	if err != nil {
		return nil, backend.WrapError("Cannot list evidence creators for an operation", backend.DatabaseErr(err))
	}

	usersDTO := make([]*dtos.User, len(users))
	for idx, user := range users {
		usersDTO[idx] = &dtos.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Slug:      user.Slug,
		}
	}
	return usersDTO, nil
}

// ListUsersForAdmin retreives standard User (public) details, and aguments with some particular fields
// meant for admin review. For use in admin views only.
func ListUsersForAdmin(ctx context.Context, db *database.Connection, i ListUsersForAdminInput) (*dtos.PaginationWrapper, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unwilling to list users", backend.UnauthorizedReadErr(err))
	}

	var users []struct {
		models.User
		AuthSchemes   *string `db:"auth_schemes"`
		UsesLocalTOTP bool    `db:"has_local_totp"`
	}

	sb := sq.Select("slug", "first_name", "last_name", "email", "admin", "disabled", "headless",
		"deleted_at", "GROUP_CONCAT(auth_scheme) AS auth_schemes").
		Column("SUM(auth_scheme='" + localauth.Code + "' AND totp_secret IS NOT NULL)>0 AS has_local_totp"). // does the user have *local* totp enabled
		From("users").
		LeftJoin("auth_scheme_data ON auth_scheme_data.user_id = users.id").
		GroupBy("users.id")

	i.AddWhere(&sb)

	if !i.IncludeDeleted {
		sb = sb.Where(sq.Eq{"deleted_at": nil})
	}

	err := i.Pagination.Select(ctx, db, &users, sb)
	if err != nil {
		return nil, backend.WrapError("Cannot list users for admin", backend.DatabaseErr(err))
	}

	usersDTO := []*dtos.UserAdminView{}
	for _, user := range users {
		// Group_Concat will return null if there are no authentication schemes listed. The below forces the schemes into a slice to avoid errors.
		authSchemes := []string{}
		if user.AuthSchemes != nil {
			authSchemes = strings.Split(*user.AuthSchemes, ",")
		}

		usersDTO = append(usersDTO, &dtos.UserAdminView{
			User: dtos.User{
				Slug:      user.Slug,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
			Email:         user.Email,
			Admin:         user.Admin,
			Headless:      user.Headless,
			AuthSchemes:   authSchemes,
			Disabled:      user.Disabled,
			UsesLocalTOTP: user.UsesLocalTOTP,
			Deleted:       user.DeletedAt != nil,
		})
	}

	return i.Pagination.WrapData(usersDTO), nil
}

func ListUsersForOperation(ctx context.Context, db *database.Connection, i ListUsersForOperationInput) ([]*dtos.UserOperationRole, error) {
	query, err := prepListUsersForOperation(ctx, db, i)
	if err != nil {
		return nil, err
	}

	var users []userAndRole
	err = db.Select(&users, *query)
	if err != nil {
		return nil, backend.WrapError("Cannot list users for operation", backend.DatabaseErr(err))
	}
	usersDTO := wrapListUsersForOperationResponse(users)
	return usersDTO, nil
}

func ListUsers(ctx context.Context, db *database.Connection, i ListUsersInput) ([]*dtos.User, error) {
	if strings.ContainsAny(i.Query, "%_") || strings.TrimFunc(i.Query, unicode.IsSpace) == "" {
		return []*dtos.User{}, nil
	}

	var users []models.User
	query := sq.Select("slug", "first_name", "last_name").
		From("users").
		Where(sq.Like{"concat(first_name, ' ', last_name)": "%" + strings.ReplaceAll(i.Query, " ", "%") + "%"}).
		OrderBy("first_name").
		Limit(10)
	if !i.IncludeDeleted {
		query = query.Where(sq.Eq{"deleted_at": nil})
	}
	err := db.Select(&users, query)
	if err != nil {
		return nil, backend.WrapError("Cannot list users", backend.DatabaseErr(err))
	}

	usersDTO := []*dtos.User{}
	for _, user := range users {
		if middleware.Policy(ctx).Check(policy.CanReadUser{UserID: user.ID}) {
			usersDTO = append(usersDTO, &dtos.User{
				Slug:      user.Slug,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			})
		}
	}
	return usersDTO, nil
}

func prepListUsersForOperation(ctx context.Context, db *database.Connection, i ListUsersForOperationInput) (*sq.SelectBuilder, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to list users for operation", backend.UnauthorizedReadErr(err))
	}

	if err := policyRequireWithAdminBypass(ctx, policy.CanListUsersOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list users for operation", backend.UnauthorizedReadErr(err))
	}

	query := sq.Select("slug", "first_name", "last_name", "role").
		From("user_operation_permissions").
		LeftJoin("users ON user_operation_permissions.user_id = users.id").
		Where(sq.Eq{"operation_id": operation.ID, "users.deleted_at": nil}).
		OrderBy("user_operation_permissions.created_at ASC")
		// OrderBy("first_name ASC", "last_name ASC", "user_operation_permissions.created_at ASC")

	i.UserFilter.AddWhere(&query)

	return &query, nil
}

func wrapListUsersForOperationResponse(users []userAndRole) []*dtos.UserOperationRole {
	usersDTO := make([]*dtos.UserOperationRole, len(users))
	for idx, user := range users {
		usersDTO[idx] = &dtos.UserOperationRole{
			User: dtos.User{
				Slug:      user.Slug,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
			Role: user.Role,
		}
	}
	return usersDTO
}

// ReadUser retrieves a detailed view of a user. This is separate from the data retriving by listing
// users, or reading another user's profile (when not an admin)
func ReadUser(ctx context.Context, db *database.Connection, userSlug string, supportedAuthSchemes *[]dtos.SupportedAuthScheme) (*dtos.UserOwnView, error) {
	userID, err := SelfOrSlugToUserID(ctx, db, userSlug)
	if err != nil {
		return nil, backend.WrapError("Unable to read user", backend.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanReadDetailedUser{UserID: userID}); err != nil {
		return nil, backend.WrapError("Unwilling to read user", backend.UnauthorizedReadErr(err))
	}

	supportedAuthCodes := make([]string, len(*supportedAuthSchemes))
	for i, scheme := range *supportedAuthSchemes {
		supportedAuthCodes[i] = scheme.SchemeCode
	}

	var user models.User
	var authSchemes []models.AuthSchemeData
	err = db.WithTx(ctx, func(tx *database.Transactable) {
		db.Get(&user, sq.Select("first_name", "last_name", "slug", "email", "admin", "headless").
			From("users").
			Where(sq.Eq{"id": userID}))

		db.Select(&authSchemes, sq.Select("username", "auth_scheme", "auth_type", "last_login").
			From("auth_scheme_data").
			Where(sq.Eq{
				"user_id":     userID,
				"auth_scheme": supportedAuthCodes,
			}))
	})
	if err != nil {
		return nil, backend.WrapError("Cannot read user", backend.DatabaseErr(err))
	}

	auths := make([]dtos.AuthenticationInfo, len(authSchemes))
	for i, v := range authSchemes {
		index := getMatchingSchemeIndex(supportedAuthSchemes, v.AuthScheme)

		auths[i] = dtos.AuthenticationInfo{
			Username:       v.Username,
			AuthSchemeCode: v.AuthScheme,
			AuthSchemeType: v.AuthType,
			AuthLogin:      v.LastLogin,
			AuthDetails:    nil,
		}
		if index > -1 {
			auths[i].AuthDetails = &(*supportedAuthSchemes)[index]
		}
	}

	return &dtos.UserOwnView{
		User: dtos.User{
			Slug:      user.Slug,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		Email:          user.Email,
		Admin:          user.Admin,
		Headless:       user.Headless,
		Authentication: auths,
	}, nil
}

func UpdateUserProfile(ctx context.Context, db *database.Connection, i UpdateUserProfileInput) error {
	var userID int64
	var err error

	if userID, err = SelfOrSlugToUserID(ctx, db, i.UserSlug); err != nil {
		return backend.WrapError("Unable to update user profile", backend.DatabaseErr(err))
	}

	if err := policy.Require(middleware.Policy(ctx), policy.CanModifyUser{UserID: userID}); err != nil {
		return backend.WrapError("Unwilling to update user profile", backend.UnauthorizedWriteErr(err))
	}

	err = db.Update(sq.Update("users").
		SetMap(map[string]interface{}{
			"first_name": i.FirstName,
			"last_name":  i.LastName,
			"email":      i.Email,
		}).
		Where(sq.Eq{"id": userID}))

	if err != nil {
		return backend.WrapError("Cannot update user profile", backend.DatabaseErr(err))
	}
	return nil
}

// SetUserFlags updates flags for the indicated user, namely: admin and disabled.
// Then removes all sessions for that user (logging them out)
//
// NOTE: The flag is to _disable_ the user, which prevents access. To enable a user, set Disabled=false
func SetUserFlags(ctx context.Context, sessionManager *scs.SessionManager, db *database.Connection, i SetUserFlagsInput) error {
	if !middleware.IsAdmin(ctx) {
		return backend.WrapError("Unwilling to set user flag", backend.UnauthorizedReadErr(fmt.Errorf("Requesting user is not an admin")))
	}

	targetUser, err := db.RetrieveUserWithAuthDataBySlug(i.Slug)
	if err != nil {
		return backend.WrapError("Cannot set user flags", backend.DatabaseErr(err))
	}
	err = validateAdminCanModifyFlag(ctx, targetUser, i)
	if err != nil {
		return backend.BadInputErr(err, err.Error())
	}

	valuesToUpdate := map[string]interface{}{}

	if i.Disabled != nil {
		valuesToUpdate["disabled"] = *i.Disabled
	}
	if i.Admin != nil {
		valuesToUpdate["admin"] = *i.Admin
	}

	if len(valuesToUpdate) > 0 && sessionManager != nil {
		sessionManager.Remove(ctx, config.SessionStoreKey())
	}
	return nil
}

// validateAdminCanModifyFlag does some checks to validate the logic/sanity of the request.
// Checks roughly include logic verifying that a user isn't elevating/demoting their own status,
// users aren't being given status that doesn't make sense (specifically: headless users cannot be admins)
func validateAdminCanModifyFlag(ctx context.Context, targetUser models.UserWithAuthData, flagsToUpdate SetUserFlagsInput) error {
	targetUserIsSelf := targetUser.ID == middleware.UserID(ctx)
	targetUserIsHeadless := targetUser.Headless

	if flagsToUpdate.Admin != nil {
		// Note on valueUpdated: the frontend will supply all values, so an extra check here is done
		// to verify the intent to change the value rather than enforcing that no value is sent, without
		// requiring that the frontend explicitly omit values
		valueUpdated := targetUser.Admin != *flagsToUpdate.Admin
		if targetUserIsSelf && valueUpdated {
			return errors.New("Admins cannot alter their own admin status")
		}
		if targetUserIsHeadless && *flagsToUpdate.Admin {
			return errors.New("Headless users cannot be granted admin status")
		}
	}

	if flagsToUpdate.Disabled != nil {
		valueUpdated := targetUser.Disabled != *flagsToUpdate.Disabled
		if targetUserIsSelf && valueUpdated {
			return errors.New("Admins cannot disable themselves")
		}
	}

	return nil
}
