// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"

	sq "github.com/Masterminds/squirrel"
)

type CreateUserGroupInput struct {
	Name      string
	UserSlugs []string
}

type ModifyUserGroupInput struct {
	Slug      string
	UserSlugs []string
}

type ListUserGroupsForAdminInput struct {
	UserFilter
	Pagination
	IncludeDeleted bool
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

func CreateUserGroup(ctx context.Context, db *database.Connection, i CreateUserGroupInput) (*dtos.CreateUserGroupOutput, error) {
	// TODO TN add validation?
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unwilling to create a user group", backend.UnauthorizedReadErr(err))
	}

	// TODO TN: how to ensure operations without users are shown?
	for {
		id, err := db.Insert("user_groups", map[string]interface{}{
			"slug": i.Name,
		})
		if err != nil {
			if database.IsAlreadyExistsError(err) {
				return nil, backend.WrapError("Unable to create user group. User group slug already exists.", backend.BadInputErr(err, "A user group with this name already exists; please choose another name"))
			}
		}
		AddUsersToGroup(db, i.UserSlugs, id)
		break
	}

	return nil, nil
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

// TODO TN - where does this get used?
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

func ListUserGroupsForAdmin(ctx context.Context, db *database.Connection, i ListUserGroupsForAdminInput) (*dtos.PaginationWrapper, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unwilling to list user groups", backend.UnauthorizedReadErr(err))
	}

	var slugMap []struct {
		UserSlug  string         `db:"user_slug"`
		GroupSlug string         `db:"group_slug"`
		Deleted   sql.NullString `db:"deleted"`
	}

	sb := sq.Select("user_groups.slug AS group_slug, users.slug AS user_slug, user_groups.deleted_at AS deleted").
		From("group_user_map").
		LeftJoin("user_groups ON group_user_map.group_id = user_groups.id").
		Join("users ON group_user_map.user_id = users.id")

	i.AddWhere(&sb)

	// write test data for this TODO TN
	if !i.IncludeDeleted {
		sb = sb.Where(sq.Eq{"user_groups.deleted_at": nil})
	}

	// err := i.Pagination.Select(ctx, db, &slugMap, sb)

	err := db.Select(&slugMap, sb)

	if err != nil {
		return nil, backend.WrapError("unable to get map of user IDs to group IDs from database", backend.DatabaseErr(err))
	}

	if err != nil {
		return nil, backend.WrapError("Cannot list user groups for admin", backend.DatabaseErr(err))
	}

	type tempGroup struct {
		Slug      string
		UserSlugs []string
		Deleted   bool
	}

	userGroupsDTO := []dtos.UserGroupAdminView{}
	tempGroupMap := dtos.UserGroupAdminView{}

	for j := 0; j < len(slugMap); j++ {
		if j == 0 {
			tempGroupMap = dtos.UserGroupAdminView{
				Slug: slugMap[j].GroupSlug,
				UserSlugs: []string{
					slugMap[j].UserSlug,
				},
				Deleted: &slugMap[j].Deleted != nil,
			}
		} else if slugMap[j].GroupSlug == slugMap[j-1].GroupSlug {
			tempGroupMap.UserSlugs = append(tempGroupMap.UserSlugs, slugMap[j].UserSlug)
			// TODO TN - make this into a part of the main clause
			if j == len(slugMap)-1 {
				userGroupsDTO = append(userGroupsDTO, tempGroupMap)
			}
		} else {
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
			tempGroupMap = dtos.UserGroupAdminView{
				Slug: slugMap[j].GroupSlug,
				UserSlugs: []string{
					slugMap[j].UserSlug,
				},
			}
		}
	}

	p := i.Pagination

	prevLastIndex := (p.Page - 1) * p.PageSize
	remainingItems := (len(userGroupsDTO) - int(prevLastIndex)) % int(p.PageSize)
	currLastIndex := int(math.Min(float64(p.Page*p.PageSize), float64(remainingItems)))
	paginatedResults := userGroupsDTO[prevLastIndex:currLastIndex]

	numPages := len(userGroupsDTO) / int(p.PageSize)
	totalPages := math.Ceil(float64(numPages))
	// TODO TN test that this loads user groups with no users
	paginatedData := &dtos.PaginationWrapper{
		PageNumber: p.Page,
		PageSize:   p.PageSize,
		Content:    paginatedResults,
		TotalCount: p.TotalCount,
		TotalPages: int64(totalPages),
	}

	return paginatedData, nil
}
