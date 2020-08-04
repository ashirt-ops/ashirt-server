// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"fmt"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"math/rand"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
)

type CreateUserInput struct {
	FirstName string
	LastName  string
	Slug      string
	Email     string
	Headless  bool
}

type CreateUserOutput struct {
	RealSlug string
	UserID   int64
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
func CreateHeadlessUser(ctx context.Context, db *database.Connection, i CreateUserInput) (CreateUserOutput, error) {
	if err := isAdmin(ctx); err != nil {
		return CreateUserOutput{}, backend.WrapError("Unable to create new headless user", backend.UnauthorizedWriteErr(err))
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
func CreateUser(db *database.Connection, i CreateUserInput) (CreateUserOutput, error) {
	validationErr := i.validate()
	if validationErr != nil {
		return CreateUserOutput{}, backend.WrapError("Unable to create new user", validationErr)
	}

	var userID int64
	var err error
	slugSuffix := ""
	var attemptedSlug string
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
				// an account with this slug already exists, attempt creating it again with a suffix
				// TODO: There's a possible, but impractical infinite loop here. We need some way to escape this
				slugSuffix = fmt.Sprintf("-%d", rand.Intn(99999))
				continue
			}
			return CreateUserOutput{}, backend.WrapError("Unable to insert new user", backend.DatabaseErr(err))
		}
		break
	}
	if userID == 1 {
		err := db.Update(sq.Update("users").Set("admin", true).Where(sq.Eq{"id": userID}))
		if err != nil {
			logging.GetSystemLogger().Log("msg", "Unable to make the first user an admin", "error", err.Error())
		}
	}
	return CreateUserOutput{
		RealSlug: attemptedSlug,
		UserID:   userID,
	}, nil
}
