package services

import (
	"context"
	"database/sql"
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ashirt-ops/ashirt-server/internal/database"
	"github.com/ashirt-ops/ashirt-server/internal/dtos"
	"github.com/ashirt-ops/ashirt-server/internal/errors"
	"github.com/ashirt-ops/ashirt-server/internal/models"
	"github.com/ashirt-ops/ashirt-server/internal/policy"
	"github.com/ashirt-ops/ashirt-server/internal/server/middleware"

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
	IncludeDeleted bool
}

type ListUserGroupsForOperationInput struct {
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
		return errors.MissingValueErr("Slug")
	}
	return nil
}

func (i CreateUserGroupInput) validateUserGroupInput() error {
	if i.Slug == "" {
		return errors.MissingValueErr("Slug")
	}
	if i.Name == "" {
		return errors.MissingValueErr("Name")
	}
	return nil
}

func AddUsersToGroup(tx *database.Transactable, userSlugs []string, groupID int64) error {
	if len(userSlugs) == 0 {
		return nil
	}

	peopleToAdd := sq.Select("id", strconv.FormatInt(groupID, 10)).
		From("users").
		Where(sq.Eq{"slug": userSlugs})

	insertQuery := sq.
		Insert("group_user_map").
		Options("Ignore").
		Columns("user_id", "group_id").Select(peopleToAdd)

	err := tx.Exec(insertQuery)

	if err != nil {
		return errors.WrapError("Unable to add users to group", errors.DatabaseErr(err))
	}
	return nil
}

func CreateUserGroup(ctx context.Context, db *database.Connection, i CreateUserGroupInput) (*dtos.UserGroup, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, errors.WrapError("Unwilling to create a user group", errors.UnauthorizedReadErr(err))
	}

	cleanSlug := SanitizeSlug(i.Slug)
	if cleanSlug == "" {
		return nil, errors.BadInputErr(stderrors.New("Unable to create operation. Invalid operation slug"), "Slug must contain english letters or numbers")
	}

	err := db.WithTx(context.Background(), func(tx *database.Transactable) {
		id, _ := tx.Insert("user_groups", map[string]interface{}{
			"slug": cleanSlug,
			"name": i.Name,
		})
		AddUsersToGroup(tx, i.UserSlugs, id)
	})

	if err != nil {
		return nil, errors.WrapError("Error creating user group", errors.BadInputErr(err, "A user group with this slug already exists; please choose another name"))
	}
	return &dtos.UserGroup{
		Slug: i.Slug,
		Name: i.Name,
	}, nil
}

func ModifyUserGroup(ctx context.Context, db *database.Connection, i ModifyUserGroupInput) (*dtos.UserGroup, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, errors.WrapError("Unwilling to modify a user group", errors.UnauthorizedReadErr(err))
	}

	if err := i.validateUserGroupInput(); err != nil {
		return nil, errors.WrapError("Unable to modify user group", errors.BadInputErr(err, "Unable to modify user group due to bad input"))
	}

	userGroup, err := lookupUserGroup(db, i.Slug)
	if err != nil {
		return nil, errors.WrapError("Unable to modify user group", errors.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(context.Background(), func(tx *database.Transactable) {
		if i.Name != "" {
			tx.Update(sq.Update("user_groups").Set("name", i.Name).Where(sq.Eq{"id": userGroup.ID}))
		}
		if len(i.UsersToRemove) > 0 {
			interfaceSlice := make([]interface{}, len(i.UsersToRemove))
			questionMarks := "("

			for i, v := range i.UsersToRemove {
				questionMarks += "?, "
				interfaceSlice[i] = v
			}

			questionMarks = strings.TrimSuffix(questionMarks, ", ")
			questionMarks += ")"

			sqlStatement := fmt.Sprintf(`DELETE gm FROM group_user_map gm JOIN users u on gm.user_id = u.id WHERE u.slug in %s;`, questionMarks)
			tx.Exec(sq.Expr(sqlStatement, interfaceSlice...))
		}
		AddUsersToGroup(tx, i.UsersToAdd, userGroup.ID)
	})
	if err != nil {
		return nil, errors.WrapError("Error creating user group", errors.BadInputErr(err, "A user group with this name already exists; please choose another name"))
	}

	return &dtos.UserGroup{
		Slug: i.Slug,
		Name: i.Name,
	}, nil
}

func DeleteUserGroup(ctx context.Context, db *database.Connection, slug string) error {
	if err := isAdmin(ctx); err != nil {
		return errors.WrapError("Unwilling to delete a user group", errors.UnauthorizedReadErr(err))
	}
	userGroup, err := lookupUserGroup(db, slug)
	if err != nil {
		return errors.WrapError("Unable to delete user group", errors.UnauthorizedWriteErr(err))
	}

	err = db.WithTx(context.Background(), func(tx *database.Transactable) {
		tx.Delete(sq.Delete("user_group_operation_permissions").Where(sq.Eq{"group_id": userGroup.ID}))
		tx.Update(sq.Update("user_groups").Set("deleted_at", time.Now()).Where(sq.Eq{"slug": slug}))
	})
	if err != nil {
		return errors.WrapError("Cannot delete user group", errors.DatabaseErr(err))
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
func ListUserGroupsForAdmin(ctx context.Context, db *database.Connection, i ListUserGroupsForAdminInput) ([]dtos.UserGroupAdminView, error) {
	if err := isAdmin(ctx); err != nil {
		return nil, errors.WrapError("Unwilling to list user groups", errors.UnauthorizedReadErr(err))
	}

	slugMap, _ := GetSlugMap(db, i)

	sortedUser, err := SortUsersInToGroups(slugMap)

	if err != nil {
		return nil, errors.WrapError("Unable to list user groups", errors.DatabaseErr(err))
	}

	return sortedUser, nil
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
		return nil, errors.WrapError("unable to get map of user IDs to group IDs from database", errors.DatabaseErr(err))
	}

	return slugMap, nil
}

func SortUsersInToGroups(slugMap SlugMap) ([]dtos.UserGroupAdminView, error) {
	userGroupsDTO := []dtos.UserGroupAdminView{}
	tempGroupMap := dtos.UserGroupAdminView{}

	if len(slugMap) == 0 {
		return []dtos.UserGroupAdminView{}, nil
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

	return userGroupsDTO, nil
}

// Lists all user groups for an operation; op admins and sys admins can view
func ListUserGroupsForOperation(ctx context.Context, db *database.Connection, i ListUserGroupsForOperationInput) ([]*dtos.UserGroupOperationRole, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err := policyRequireWithAdminBypass(ctx, policy.CanListUserGroupsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, errors.WrapError("Unwilling to list usergroups", errors.UnauthorizedReadErr(err))
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
		return nil, errors.WrapError("Cannot list user groups for operation", errors.DatabaseErr(err))
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
func ListUserGroups(ctx context.Context, db *database.Connection, i ListUserGroupsInput) ([]*dtos.UserGroupAdminView, error) {
	operation, err := lookupOperation(db, i.OperationSlug)
	if err := policyRequireWithAdminBypass(ctx, policy.CanListUserGroupsOfOperation{OperationID: operation.ID}); err != nil {
		return nil, errors.WrapError("Unwilling to list usergroups", errors.UnauthorizedReadErr(err))
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
		return nil, errors.WrapError("Cannot list user groups", errors.DatabaseErr(err))
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
