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
	UserGroupFilter
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

func (i ModifyUserGroupInput) validateUserGroupInput() error {
	if i.Slug == "" {
		return backend.MissingValueErr("Slug")
	}
	return nil
}

func (i CreateUserGroupInput) validateUserGroupInput() error {
	if i.Slug == "" {
		return backend.MissingValueErr("Slug")
	}
	if i.Name == "" {
		return backend.MissingValueErr("Name")
	}
	return nil
}

func AddUsersToGroup(tx *database.Transactable, userSlugs []string, groupID int64) error {
	if len(userSlugs) > 0 {
		questionMarks := "("

		interfaceSlice := make([]interface{}, len(userSlugs)+1)
		for i, v := range userSlugs {
			questionMarks += "?, "
			interfaceSlice[i] = v
		}

		questionMarks = strings.TrimSuffix(questionMarks, ", ")
		questionMarks += ")"

		sqlStatement := fmt.Sprintf(`INSERT IGNORE INTO group_user_map(user_id, group_id)
					SELECT users.id, user_groups.id
					FROM users, user_groups
					WHERE users.slug in %s and user_groups.id = ?;`, questionMarks)

		interfaceSlice[len(userSlugs)] = groupID
		err := tx.Exec(sq.Expr(sqlStatement, interfaceSlice...))

		if err != nil {
			return backend.WrapError("Unable to add users to group", backend.DatabaseErr(err))
		}
		return nil
	}

	return nil
}

func CreateUserGroup(ctx context.Context, db *database.Connection, i CreateUserGroupInput) (*dtos.UserGroup, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unwilling to create a user group", backend.UnauthorizedReadErr(err))
	}

	cleanSlug := SanitizeSlug(i.Slug)
	if cleanSlug == "" {
		return nil, backend.BadInputErr(errors.New("Unable to create operation. Invalid operation slug"), "Slug must contain english letters or numbers")
	}

	err := db.WithTx(context.Background(), func(tx *database.Transactable) {
		id, _ := tx.Insert("user_groups", map[string]interface{}{
			"slug": cleanSlug,
			"name": i.Name,
		})
		AddUsersToGroup(tx, i.UserSlugs, id)
	})

	if err != nil {
		return nil, backend.WrapError("Error creating user group", backend.BadInputErr(err, "A user group with this slug already exists; please choose another name"))
	}
	return &dtos.UserGroup{
		Slug: i.Slug,
		Name: i.Name,
	}, nil
}

func ModifyUserGroup(ctx context.Context, db *database.Connection, i ModifyUserGroupInput) (*dtos.UserGroup, error) {
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
		if len(i.UsersToRemove) > 0 {
			for _, userSlug := range i.UsersToRemove {
				tx.Exec(sq.Expr("DELETE gm FROM group_user_map gm JOIN users u on gm.user_id = u.id WHERE u.slug=?;", userSlug))
			}
		}
		AddUsersToGroup(tx, i.UsersToAdd, userGroup.ID)
	})
	if err != nil {
		return nil, backend.WrapError("Error creating user group", backend.BadInputErr(err, "A user group with this name already exists; please choose another name"))
	}

	return &dtos.UserGroup{
		Slug: i.Slug,
		Name: i.Name,
	}, nil
}

func DeleteUserGroup(ctx context.Context, db *database.Connection, slug string) error {
	if err := isAdmin(ctx); err != nil {
		return backend.WrapError("Unwilling to delete a user group", backend.UnauthorizedReadErr(err))
	}
	userGroup, err := lookupUserGroup(db, slug)
	if err != nil {
		return backend.WrapError("Unable to delete user group", backend.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.Delete(sq.Delete("user_group_operation_permissions").Where(sq.Eq{"group_id": userGroup.ID}))
		tx.Update(sq.Update("user_groups").Set("deleted_at", time.Now()).Where(sq.Eq{"slug": slug}))
	})
	if err != nil {
		return backend.WrapError("Cannot delete user group", backend.DatabaseErr(err))
	}

	return nil
}

type SlugMap []struct {
	UserSlug  sql.NullString `db:"user_slug"`
	GroupSlug string         `db:"group_slug"`
	GroupName string         `db:"group_name"`
	Deleted   sql.NullString `db:"deleted"`
}

// Lists all usergroups for an admin, with pagination
func ListUserGroupsForAdmin(ctx context.Context, db *database.Connection, i ListUserGroupsForAdminInput) (*dtos.PaginationWrapper, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, backend.WrapError("Unwilling to list user groups", backend.UnauthorizedReadErr(err))
	}

	slugMap, _ := GetSlugMap(db, i)

	paginatedSortedUser, err := SortUsersInToGroups(slugMap, i.Pagination)

	if err != nil {
		return nil, backend.WrapError("Unable to list user groups", backend.DatabaseErr(err))
	}

	return paginatedSortedUser, nil
}

func GetSlugMap(db *database.Connection, i ListUserGroupsForAdminInput) (SlugMap, error) {
	sb := sq.Select("user_groups.slug AS group_slug, user_groups.name AS group_name, users.slug AS user_slug, user_groups.deleted_at AS deleted").
		From("group_user_map").
		Join("users ON group_user_map.user_id = users.id").
		RightJoin("user_groups ON group_user_map.group_id = user_groups.id")

	i.AddWhere(&sb)

	if !i.IncludeDeleted {
		sb = sb.Where(sq.Eq{"user_groups.deleted_at": nil})
	}

	sb = sb.OrderBy("group_name")

	var slugMap SlugMap

	err := db.Select(&slugMap, sb)

	if err != nil {
		return nil, backend.WrapError("unable to get map of user IDs to group IDs from database", backend.DatabaseErr(err))
	}

	return slugMap, nil
}

func SortUsersInToGroups(slugMap SlugMap, pagination Pagination) (*dtos.PaginationWrapper, error) {
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
		groupWithNoUsers := !hasUserSlug
		sameGroupAsPrev := false
		if j > 0 {
			sameGroupAsPrev = slugMap[j].GroupSlug == slugMap[j-1].GroupSlug
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
		} else if firstItem && groupWithNoUsers {
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
		} else if otherItem && diffGroup && groupWithNoUsers {
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
			tempGroupMap = dtos.UserGroupAdminView{
				Slug:    slugMap[j].GroupSlug,
				Name:    slugMap[j].GroupName,
				Deleted: slugMap[j].Deleted.Valid,
			}
		} else if isLastItem && sameGroupAsPrev && hasUserSlug {
			tempGroupMap.UserSlugs = append(tempGroupMap.UserSlugs, slugMap[j].UserSlug.String)
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
		} else if isLastItem && groupWithNoUsers {
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
			tempGroupMap = dtos.UserGroupAdminView{
				Slug:    slugMap[j].GroupSlug,
				Name:    slugMap[j].GroupName,
				Deleted: slugMap[j].Deleted.Valid,
			}
			userGroupsDTO = append(userGroupsDTO, tempGroupMap)
		}
	}

	prevLastIndex := (pagination.Page - 1) * pagination.PageSize
	groupLength := len(userGroupsDTO)
	totalPages := math.Ceil(float64(groupLength) / float64(pagination.PageSize))
	remainingItemsCount := (groupLength - int(prevLastIndex)) % int(pagination.PageSize)

	currLastIndex := int(pagination.Page * pagination.PageSize)
	pageSize := pagination.PageSize
	if pagination.Page == int64(totalPages) {
		currLastIndex = int(prevLastIndex) + remainingItemsCount
		pageSize = int64(remainingItemsCount)
	}

	paginatedResults := userGroupsDTO[prevLastIndex:currLastIndex]
	paginatedData := &dtos.PaginationWrapper{
		PageNumber: pagination.Page,
		PageSize:   pageSize,
		Content:    paginatedResults,
		TotalCount: int64(groupLength),
		TotalPages: int64(totalPages),
	}
	return paginatedData, nil
}

// Lists all user groups for an operation; op admins and sys admins can view
func ListUserGroupsForOperation(ctx context.Context, db *database.Connection, i ListUserGroupsForOperationInput) ([]*dtos.UserGroupOperationRole, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err := policyRequireWithAdminBypass(ctx, policy.CanListUserGroupsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, backend.WrapError("Unwilling to list usergroups", backend.UnauthorizedReadErr(err))
	}

	query := sq.Select("slug", "name", "role").
		From("user_group_operation_permissions").
		LeftJoin("user_groups ON user_group_operation_permissions.group_id = user_groups.id").
		Where(sq.Eq{"operation_id": operation.ID, "user_groups.deleted_at": nil}).
		OrderBy("user_group_operation_permissions.created_at ASC")

	i.UserGroupFilter.AddWhere(&query)

	var userGroups []userGroupAndRole
	err = db.Select(&userGroups, query)
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

// lists all user groups that can be added to an operation
// no pagination, because this is used for the search bar
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
}
