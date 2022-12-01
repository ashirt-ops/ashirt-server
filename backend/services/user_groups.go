// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/models"

	sq "github.com/Masterminds/squirrel"
)

type ModifyUserGroupInput struct {
	Slug      string
	UserSlugs []string
}

func (cugi ModifyUserGroupInput) validateUserGroupInput() error {
	if cugi.Slug == "" {
		return backend.MissingValueErr("Slug")
	}
	if len(cugi.UserSlugs) < 1 {
		return backend.MissingValueErr("User Slugs")
	}
	return nil
}

// TODO TN: how does a group get set up with an operation?
func AddUsersToGroup(db *database.Connection, userSlugs []string, groupID int64) error {
	fmt.Println("Adding users to group", userSlugs)
	for _, userSlug := range userSlugs {
		userID, err := userSlugToUserID(db, userSlug)
		if err != nil {
			return backend.WrapError("Unable to get user id from slug", backend.BadInputErr(err, fmt.Sprintf(`No user with slug %s was found`, userSlug)))
		}

		var userGroupMap models.UserGroupMap
		err = db.Get(&userGroupMap, sq.Select("*").
			From("group_user_map").
			Where(sq.Eq{
				"user_id":  userID,
				"group_id": groupID,
			}))
		if err != nil {
			_, err = db.Insert("group_user_map", map[string]interface{}{
				"user_id":  userID,
				"group_id": groupID,
			})
			if err != nil {
				return backend.WrapError("Unable to connect user to group", backend.DatabaseErr(err))
			}
		}
	}

	return nil
}

func RemoveUsersFromGroup(db *database.Connection, userSlugs []string, groupID int64) error {
	for _, userSlug := range userSlugs {

		userID, err := userSlugToUserID(db, userSlug)
		if err != nil {
			return backend.WrapError("Unable to get user id from slug", backend.BadInputErr(err, fmt.Sprintf(`No user with slug %s was found`, userSlug)))
		}

		var userGroupMap models.UserGroupMap
		err = db.Get(&userGroupMap, sq.Select("*").
			From("group_user_map").
			Where(sq.Eq{
				"user_id":  userID,
				"group_id": groupID,
			}))
		if err == nil {
			err := db.Delete(sq.Delete("group_user_map").Where(sq.Eq{"user_id": userID, "group_id": groupID}))

			if err != nil {
				return backend.WrapError("Cannot delete user role", backend.DatabaseErr(err))
			}
			return nil
		}
	}

	return nil
}

func CreateUserGroup(db *database.Connection, i ModifyUserGroupInput) (*dtos.CreateUserGroupOutput, error) {
	validationErr := i.validateUserGroupInput()
	if validationErr != nil {
		return nil, backend.WrapError("Unable to create new user group", validationErr)
	}

	var userGroupID int64
	var err error
	slugSuffix := ""
	var attemptedSlug string
	attemptNumber := 1
	for {
		attemptedSlug = i.Slug + slugSuffix
		userGroupID, err = db.Insert("user_groups", map[string]interface{}{
			"slug": attemptedSlug,
		})
		if err != nil {
			if database.IsAlreadyExistsError(err) {
				if attemptNumber > 5 {
					return nil, backend.WrapError("Unable to create new user group after many attempts", backend.DatabaseErr(err))
				}

				logging.GetSystemLogger().Log(
					"msg", "Unable to create user group with slug; trying alternative",
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
			return nil, backend.WrapError("Unable to insert new user group", backend.DatabaseErr(err))
		}
		break
	}

	AddUsersToGroup(db, i.UserSlugs, userGroupID)
	return &dtos.CreateUserGroupOutput{
		RealSlug:    attemptedSlug,
		UserGroupID: userGroupID,
	}, nil
}

func DeleteUserGroup(db *database.Connection, slug string) error {
	userGroupID, err := userGroupSlugToUserGroupID(db, slug)
	if err != nil {
		return backend.WrapError("User group does not exist and therefore cannot be deleted", backend.DatabaseErr(err))
	}

	err = db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.Delete(sq.Delete("group_user_map").Where(sq.Eq{"group_id": userGroupID}))
		tx.Update(sq.Update("user_groups").Set("deleted_at", time.Now()).Where(sq.Eq{"slug": slug}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete user group", backend.DatabaseErr(err))
	}

	return nil
}

func GetUserIDsFromGroup(db *database.Connection, groupID int64) ([]int64, error) {
	var userGroupMap []int64
	err := db.Select(&userGroupMap, sq.Select("user_id").
		From("group_user_map").
		Where(sq.Eq{
			"group_id": groupID,
		}))
	if err != nil {
		s := fmt.Sprintf("Cannot get user group map for group %d", groupID)
		return userGroupMap, backend.WrapError(s, backend.DatabaseErr(err))
	}
	return userGroupMap, nil
}

// TODO TN - remove this?
// func GetUserIDsFromGroup(db *database.Connection, groupID int64) ([]models.UserGroupMap, error) {
// 	var userGroupMap []models.UserGroupMap
// 	// TODO TN should I return all here, or just the user IDs?
// 	err := db.Select(&userGroupMap, sq.Select("*").
// 		From("group_user_map").
// 		Where(sq.Eq{
// 			"group_id": groupID,
// 		}))
// 	if err != nil {
// 		s := fmt.Sprintf("Cannot get user group map for group %d", groupID)
// 		return userGroupMap, backend.WrapError(s, backend.DatabaseErr(err))
// 	}
// 	return userGroupMap, nil
// }
