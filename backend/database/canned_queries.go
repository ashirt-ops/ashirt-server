// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package database

import (
	"strconv"
	"strings"
	"time"

	"github.com/theparanoids/ashirt/backend/models"

	sq "github.com/Masterminds/squirrel"
)

// RetrieveUserByID retrieves a full user from the users table given a user ID
func (c *Connection) RetrieveUserByID(userID int64) (models.User, error) {
	var rtn models.User
	err := c.Get(&rtn, sq.Select("*").From("users").Where(sq.Eq{"id": userID}))
	return rtn, err
}

// RetrieveUserBySlug retrieves a full user from the users table give a user slug
func (c *Connection) RetrieveUserBySlug(slug string) (models.User, error) {
	var rtn models.User
	err := c.Get(&rtn, sq.Select("*").From("users").Where(sq.Eq{"slug": slug}))
	return rtn, err
}

// RetrieveUserIDBySlug retrieves a user's ID from a given slug. Likely faster than retriving the
// full record, so this is preferred if all you need to the slug/id conversion
func (c *Connection) RetrieveUserIDBySlug(slug string) (int64, error) {
	var rtn int64
	err := c.Get(&rtn, sq.Select("id").From("users").Where(sq.Eq{"slug": slug}))
	return rtn, err
}

// RetrieveUserWithAuthDataBySlug retrieves a full user from the users table given a slug. Includes
// data from the auth_scheme_data table (namely, scheme names)
func (c *Connection) RetrieveUserWithAuthDataBySlug(slug string) (models.UserWithAuthData, error) {
	return c.retrieveUserWithAuthData(models.User{Slug: slug})
}

// RetrieveUserWithAuthDataByID retrieves a full user from the users table given that user's ID.
// Includes data from the auth_scheme_data table (namely, scheme names)
func (c *Connection) RetrieveUserWithAuthDataByID(userID int64) (models.UserWithAuthData, error) {
	return c.retrieveUserWithAuthData(models.User{ID: userID})
}

func (c *Connection) retrieveUserWithAuthData(user models.User) (models.UserWithAuthData, error) {
	var protoUserWithAuthData struct {
		models.User
		AuthSchemes *string `db:"auth_schemes"`
		LastLogins  *string `db:"last_logins"`
	}

	query := sq.Select("users.*",
		"GROUP_CONCAT(auth_scheme) AS auth_schemes",
		"GROUP_CONCAT(UNIX_TIMESTAMP(IFNULL(last_login, 0)) ) as last_logins").
		From("users").
		LeftJoin("auth_scheme_data ON users.id = user_id").
		GroupBy("users.id")

	if user.Slug != "" {
		query = query.Where(sq.Eq{"slug": user.Slug})
	}
	if user.ID != 0 {
		query = query.Where(sq.Eq{"users.id": user.ID})
	}

	err := c.Get(&protoUserWithAuthData, query)
	if err != nil {
		return models.UserWithAuthData{}, err
	}

	var schemes []models.LimitedAuthSchemeData
	if protoUserWithAuthData.AuthSchemes != nil {
		schemeNames := strings.Split(*protoUserWithAuthData.AuthSchemes, ",")
		loginDates := strings.Split(*protoUserWithAuthData.LastLogins, ",")

		schemes = make([]models.LimitedAuthSchemeData, len(schemeNames))
		for i, schemeName := range schemeNames {
			secs, err := strconv.ParseInt(loginDates[i], 10, 64)
			if err != nil {
				secs = 0
			}
			var loginDate *time.Time
			if secs != 0 {
				t := time.Unix(secs, 0)
				loginDate = &t
			}
			schemes[i] = models.LimitedAuthSchemeData{
				AuthScheme: schemeName,
				LastLogin:  loginDate,
			}
		}
	}

	rtn := models.UserWithAuthData{
		User:           protoUserWithAuthData.User,
		AuthSchemeData: schemes,
	}

	return rtn, nil
}
