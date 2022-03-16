// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package middleware

import (
	"crypto/hmac"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/signer"

	sq "github.com/Masterminds/squirrel"
)

// Max allowed difference between the current time and passed Date header
const maxDateDelta = time.Hour

type UserData struct {
	ID       int64
	Headless bool
}

func authenticateAPI(db *database.Connection, r *http.Request, requestBody io.Reader) (UserData, error) {
	emptyUserData := UserData{ID: -1, Headless: false}
	if err := checkDateHeader(r.Header.Get("Date")); err != nil {
		return emptyUserData, backend.WrapError("Unable to parse date header (for api auth)", err)
	}

	// Check HMAC
	accessKey, headerHMAC, err := parseAuthorizationHeader(r.Header.Get("Authorization"))
	if err != nil {
		return emptyUserData, backend.WrapError("Unable to parse (api) authorization header", err)
	}

	var apiKey struct {
		models.APIKey
		DisabledFlag bool `db:"disabled"`
		Headless     bool `db:"headless"`
	}

	// var apiKey models.APIKey
	// Defer checking error here to avoid timing attacks to discover valid access keys
	err = db.Get(&apiKey, sq.Select("secret_key", "user_id", "disabled", "headless").
		From("api_keys").
		LeftJoin("users ON users.id = user_id").
		Where(sq.Eq{"access_key": accessKey}))
	expectedHMAC := signer.BuildRequestHMAC(r, requestBody, apiKey.SecretKey)
	if !hmac.Equal(headerHMAC, expectedHMAC) {
		return emptyUserData, errors.New("Bad HMAC")
	}
	if err != nil {
		return emptyUserData, backend.WrapError("Unable to retrieve API key data", err)
	}
	if apiKey.DisabledFlag {
		return emptyUserData, backend.DisabledUserError()
	}

	err = db.Update(sq.Update("api_keys").Set("last_auth", time.Now()).Where(sq.Eq{"access_key": accessKey}))
	if err != nil {
		logging.Log(r.Context(), "msg", "Failed to update last_auth", "access_key", accessKey, "error", err)
	}

	return UserData{ID: apiKey.UserID, Headless: apiKey.Headless}, nil
}

// parseAuthorizationHeader parses the authorization header and returns the access key and HMAC
func parseAuthorizationHeader(authorizationStr string) (string, []byte, error) {
	if authorizationStr == "" {
		return "", []byte{}, errors.New("Missing required Authorization header")
	}

	split := strings.SplitN(authorizationStr, ":", 2)
	if len(split) != 2 {
		return "", []byte{}, errors.New("Missing required HMAC signature in Authorization header")
	}
	accessKey, base64HMAC := split[0], split[1]

	headerHMAC, err := base64.StdEncoding.DecodeString(base64HMAC)
	if err != nil {
		return accessKey, headerHMAC, backend.WrapError("Unable to decode base64 HMAC", err)
	}
	return accessKey, headerHMAC, nil
}

// checkDateHeader verifies that the passed Date header is valid and within the maxDateDelta of the current time
func checkDateHeader(dateStr string) error {
	if dateStr == "" {
		return errors.New("Missing required Date header")
	}

	parsedDate, err := time.Parse(time.RFC1123, dateStr)
	if err != nil {
		return err
	}

	if parsedDate.Location().String() != "GMT" {
		// RFC7231 specifies the Date header must always be in GMT
		// Enforcing this avoids bugs where go silently converts unknown timestamps to
		// UTC which may happen in docker containers where tzdata isn't installed
		return fmt.Errorf("Date header must be in GMT (got %s)", dateStr)
	}

	delta := time.Since(parsedDate)
	dateIsWithinMaxDelta := (0 < delta && delta < maxDateDelta) || (-maxDateDelta < delta && delta < 0)
	if !dateIsWithinMaxDelta {
		return fmt.Errorf("Date %s is not within max delta (%s) from current time %s", parsedDate, maxDateDelta, time.Now())
	}

	return nil
}
