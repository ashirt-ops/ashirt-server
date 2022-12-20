// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"unicode"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateUserGroupInput struct {
	Name      string
	Slug      string
	UserSlugs []string
}

type ModifyUserGroupInput struct {
	Name          string
	Slug          string
	UsersToAdd    []string
	UsersToRemove []string
}

type ListUserGroupsForAdminInput struct {
	UserFilter
	Pagination
	IncludeDeleted bool
}

type ListUserGroupsForOperationInput struct {
	Pagination
	UserGroupFilter
	OperationSlug string
}

type userGroupAndRole struct {
	models.UserGroup
	Role policy.OperationRole `db:"role"`
}

type ListUserGroupsInput struct {
	Query          string
	IncludeDeleted bool
	OperationSlug  string
}

func (cugi ModifyUserGroupInput) validateUserGroupInput() error {
	if cugi.Slug == "" {
		return backend.MissingValueErr("Slug")
	}
	if cugi.Slug == "" {
		return backend.MissingValueErr("Name")
	}
	return nil
}

func (cugi CreateUserGroupInput) validateUserGroupInput() error {
	if cugi.Slug == "" {
		return backend.MissingValueErr("Slug")
	}
	if cugi.Name == "" {
		return backend.MissingValueErr("Name")
	}
	return nil
}

func AddUsersToGroup(db *database.Connection, userSlugs []string, groupID int64) error {
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

		err = db.Delete(sq.Delete("group_user_map").Where(sq.Eq{"user_id": userID, "group_id": groupID}))

		if err != nil {
			return backend.WrapError("Cannot delete user role", backend.DatabaseErr(err))
		}
	}

	return nil
}

// TODO TN ask Joel about return values
// TODO TN look it up
func CreateUserGroup(ctx context.Context, db *database.Connection, i CreateUserGroupInput) (*dtos.UserGroupOutput, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unwilling to create a user group", backend.UnauthorizedReadErr(err))
	}

	cleanSlug := SanitizeSlug(i.Slug)
	if cleanSlug == "" {
		return nil, backend.BadInputErr(errors.New("Unable to create operation. Invalid operation slug"), "Slug must contain english letters or numbers")
	}

	for {
		id, err := db.Insert("user_groups", map[string]interface{}{
			"slug": cleanSlug,
			"name": i.Name,
		})
		if err != nil {
			if database.IsAlreadyExistsError(err) {
				return nil, backend.WrapError("Unable to create user group. User group slug already exists.", backend.BadInputErr(err, "A user group with this slug already exists; please choose another name"))
			}
		}
		err = AddUsersToGroup(db, i.UserSlugs, id)
		if err != nil {
			// TODO TN fix wrapped error
			// rollback creatation of user group TODO TN
			return nil, backend.WrapError("Unable to add users to user group.", backend.BadInputErr(err, "Unable to create add users to user group."))
		}
		break
	}

	return nil, nil
}

func ListUserGroupsForOperation(ctx context.Context, db *database.Connection, i ListUserGroupsForOperationInput) ([]*dtos.UserGroupOperationRole, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err := policyRequireWithAdminBypass(ctx, policy.CanListUserGroupsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list usergroups", backend.UnauthorizedReadErr(err))
	}
	query, err := prepListUserGroupsForOperation(ctx, db, i, operation.ID)
	if err != nil {
		return nil, err
	}

	var userGroups []userGroupAndRole
	err = db.Select(&userGroups, *query)
	if err != nil {
		return nil, backend.WrapError("Cannot list user groups for operation", backend.DatabaseErr(err))
	}
	userGroupsDTO := wrapListUserGroupsForOperationResponse(userGroups)
	return userGroupsDTO, nil
}

func wrapListUserGroupsForOperationResponse(userGroups []userGroupAndRole) []*dtos.UserGroupOperationRole {
	userGroupsDTO := make([]*dtos.UserGroupOperationRole, len(userGroups))
	for idx, userGroup := range userGroups {
		userGroupsDTO[idx] = &dtos.UserGroupOperationRole{
			UserGroup: dtos.UserGroupAdminView{
				Slug: userGroup.Slug,
				Name: userGroup.Name,
			},
			Role: userGroup.Role,
		}
	}
	return userGroupsDTO
}

// write a function that modifies a user group
func ModifyUserGroup(ctx context.Context, db *database.Connection, i ModifyUserGroupInput) (*dtos.UserGroupOutput, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unwilling to modify a user group", backend.UnauthorizedReadErr(err))
	}

	if err := i.validateUserGroupInput(); err != nil {
		return nil, backend.WrapError("Unable to modify user group", backend.BadInputErr(err, "Unable to modify user group due to bad input"))
	}

	userGroup, err := lookupUserGroup(db, i.Slug)
	if err != nil {
		return nil, backend.WrapError("Unable to modify user group", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(context.Background(), func(tx *database.Transactable) {
		if i.Name != "" {
			tx.Update(sq.Update("user_groups").Set("name", i.Name).Where(sq.Eq{"id": userGroup.ID}))
		}
		// TODO TN figure out how to make these transactions work
		RemoveUsersFromGroup(db, i.UsersToRemove, userGroup.ID)
		AddUsersToGroup(db, i.UsersToAdd, userGroup.ID)
	})
	if err != nil {
		return nil, backend.WrapError("Unable to modify user group", backend.DatabaseErr(err))
	}

	return nil, nil

}

// TODO TN are these return values similar to other functions?
func DeleteUserGroup(ctx context.Context, db *database.Connection, slug string) error {
	if err := isAdmin(ctx); err != nil {
		return backend.WrapError("Unwilling to delete a user group", backend.UnauthorizedReadErr(err))
	}
	userGroup, err := lookupUserGroup(db, slug)
	if err != nil {
		return backend.WrapError("Unable to delete user group", backend.UnauthorizedWriteErr(err))
	}

	// if err := policyRequireWithAdminBypass(ctx, policy.CanDeleteOperation{UsergroupID: userGroup.ID}); err != nil {
	// 	return backend.WrapError("Unwilling to delete user group", backend.UnauthorizedWriteErr(err))
	// }
	// TODO TN ADd this in later

	err = db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.Delete(sq.Delete("user_group_operation_permissions").Where(sq.Eq{"group_id": userGroup.ID}))
		tx.Update(sq.Update("user_groups").Set("deleted_at", time.Now()).Where(sq.Eq{"slug": slug}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete user group", backend.DatabaseErr(err))
	}

	return nil
}

var slugMap []struct {
	UserSlug  sql.NullString `db:"user_slug"`
	GroupSlug string         `db:"group_slug"`
	GroupName string         `db:"group_name"`
	Deleted   sql.NullString `db:"deleted"`
}

type tempGroup struct {
	Slug      string
	UserSlugs []string
	Deleted   bool
}

// TODO TN how to to more thoroughly test this?
func ListUserGroupsForAdmin(ctx context.Context, db *database.Connection, i ListUserGroupsForAdminInput) (*dtos.PaginationWrapper, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unwilling to list user groups", backend.UnauthorizedReadErr(err))
	}

	sb := sq.Select("user_groups.slug AS group_slug, user_groups.name AS group_name, users.slug AS user_slug, user_groups.deleted_at AS deleted").
		From("group_user_map").
		LeftJoin("user_groups ON group_user_map.group_id = user_groups.id").
		Join("users ON group_user_map.user_id = users.id")

	i.AddWhere(&sb)

	if !i.IncludeDeleted {
		sb = sb.Where(sq.Eq{"user_groups.deleted_at": nil})
	}

	sb2 := sq.Select("user_groups.slug AS group_slug, user_groups.name AS group_name, NULL as user_slug, user_groups.deleted_at AS deleted").
		From("user_groups")

	if !i.IncludeDeleted {
		sb2 = sb2.Where(sq.Eq{"deleted_at": nil})
	}

	sb2 = sb2.OrderBy("group_name")

	sql, args, _ := sb2.ToSql()
	unionSelect := sb.Suffix("UNION "+sql, args...)

	err := db.Select(&slugMap, unionSelect)

	if err != nil {
		return nil, backend.WrapError("unable to get map of user IDs to group IDs from database", backend.DatabaseErr(err))
	}

	userGroupsDTO := []dtos.UserGroupAdminView{}
	tempGroupMap := dtos.UserGroupAdminView{}

	if len(slugMap) == 0 {
		return &dtos.PaginationWrapper{
			PageNumber: 1,
			PageSize:   0,
			TotalCount: int64(0),
			TotalPages: int64(1),
		}, nil
	}

	for j := 0; j < len(slugMap); j++ {
		firstItem := j == 0
		isLastItem := j == len(slugMap)-1
		otherItem := j > 0 && j < len(slugMap)-1
		hasUserSlug := slugMap[j].UserSlug.Valid
		noUserSlug := !hasUserSlug
		sameGroupAsPrev := false
		if j > 0 {
			sameGroupAsPrev = slugMap[j].GroupName == slugMap[j-1].GroupName
		}
		diffGroup := !sameGroupAsPrev

		if firstItem && hasUserSlug {
			tempGroupMap = dtos.UserGroupAdminView{
				Slug: slugMap[j].GroupSlug,
				Name: slugMap[j].GroupName,
				UserSlugs: []string{
					slugMap[j].UserSlug.String,
				},
				Deleted: slugMap[j].Deleted.Valid,
			}
		} else if firstItem && noUserSlug {
			tempGroupMap = dtos.UserGroupAdminView{
				Slug:    slugMap[j].GroupSlug,
				Name:    slugMap[j].GroupName,
				Deleted: slugMap[j].Deleted.Valid,
			}
		} else if otherItem && sameGroupAsPrev && hasUserSlug {
			tempGroupMap.UserSlugs = append(tempGroupMap.UserSlugs, slugMap[j].UserSlug.String)
		} else if otherItem && diffGroup && hasUserSlug {
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
			tempGroupMap = dtos.UserGroupAdminView{
				Slug: slugMap[j].GroupSlug,
				Name: slugMap[j].GroupName,
				UserSlugs: []string{
					slugMap[j].UserSlug.String,
				},
				Deleted: slugMap[j].Deleted.Valid,
			}
		} else if otherItem && diffGroup && noUserSlug {
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
			tempGroupMap = dtos.UserGroupAdminView{
				Slug:    slugMap[j].GroupSlug,
				Name:    slugMap[j].GroupName,
				Deleted: slugMap[j].Deleted.Valid,
			}
		} else if isLastItem && sameGroupAsPrev && hasUserSlug {
			tempGroupMap.UserSlugs = append(tempGroupMap.UserSlugs, slugMap[j].UserSlug.String)
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
		} else if isLastItem && sameGroupAsPrev && noUserSlug {
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
		} else if isLastItem && diffGroup && hasUserSlug {
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
			tempGroupMap = dtos.UserGroupAdminView{
				Slug: slugMap[j].GroupSlug,
				Name: slugMap[j].GroupName,
				UserSlugs: []string{
					slugMap[j].UserSlug.String,
				},
				Deleted: slugMap[j].Deleted.Valid,
			}
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
		} else if isLastItem && diffGroup && noUserSlug {
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
			tempGroupMap = dtos.UserGroupAdminView{
				Slug:    slugMap[j].GroupSlug,
				Name:    slugMap[j].GroupName,
				Deleted: slugMap[j].Deleted.Valid,
			}
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
		}
	}

	p := i.Pagination

	prevLastIndex := (p.Page - 1) * p.PageSize
	groupLength := len(userGroupsDTO)
	totalPages := math.Ceil(float64(groupLength) / float64(p.PageSize))
	remainingItemsCount := (groupLength - int(prevLastIndex)) % int(p.PageSize)

	currLastIndex := int(p.Page * p.PageSize)
	pageSize := p.PageSize
	if p.Page == int64(totalPages) {
		currLastIndex = int(prevLastIndex) + remainingItemsCount
		pageSize = int64(remainingItemsCount)
	}

	paginatedResults := userGroupsDTO[prevLastIndex:currLastIndex]
	paginatedData := &dtos.PaginationWrapper{
		PageNumber: p.Page,
		PageSize:   pageSize,
		Content:    paginatedResults,
		TotalCount: int64(groupLength),
		TotalPages: int64(totalPages),
	}

	return paginatedData, nil
}

// TODO TN add tests for new functions
// TODO TN hide groups label from non-admins

func ListUserGroups(ctx context.Context, db *database.Connection, i ListUserGroupsInput) ([]*dtos.UserGroupAdminView, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err := policyRequireWithAdminBypass(ctx, policy.CanListUserGroupsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list usergroups", backend.UnauthorizedReadErr(err))
	}

	if strings.ContainsAny(i.Query, "%_") || strings.TrimFunc(i.Query, unicode.IsSpace) == "" {
		return []*dtos.UserGroupAdminView{}, nil
	}

	var userGroups []models.UserGroup
	query := sq.Select("slug", "name").
		From("user_groups").
		Where(sq.Like{"name": "%" + strings.ReplaceAll(i.Query, " ", "%") + "%"}).
		OrderBy("name").
		Limit(10)
	if !i.IncludeDeleted {
		query = query.Where(sq.Eq{"deleted_at": nil})
	}
	err = db.Select(&userGroups, query)
	if err != nil {
		return nil, backend.WrapError("Cannot list user groups", backend.DatabaseErr(err))
	}

	userGroupsDTO := []*dtos.UserGroupAdminView{}
	for _, userGroup := range userGroups {
		if middleware.Policy(ctx).Check(policy.CanReadUser{UserID: userGroup.ID}) {
			userGroupsDTO = append(userGroupsDTO, &dtos.UserGroupAdminView{
				Slug: userGroup.Slug,
				Name: userGroup.Name,
			})
		}
	}
	return userGroupsDTO, nil
	// TODO TN should I call user gruops - groups? Doesn't work in DB, but could work elsewhere
	// TODO TN - editing a group and changing the slug and the users, it doesn't work
}

func prepListUserGroupsForOperation(ctx context.Context, db *database.Connection, i ListUserGroupsForOperationInput, operationID int64) (*sq.SelectBuilder, error) {
	query := sq.Select("slug", "name", "role").
		From("user_group_operation_permissions").
		LeftJoin("user_groups ON user_group_operation_permissions.group_id = user_groups.id").
		Where(sq.Eq{"operation_id": operationID, "user_groups.deleted_at": nil}).
		OrderBy("user_group_operation_permissions.created_at ASC")

	i.UserGroupFilter.AddWhere(&query)

	return &query, nil
}
