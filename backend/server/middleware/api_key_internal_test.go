package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/signer"

	"github.com/stretchr/testify/require"
)

func TestCheckDateHeader(t *testing.T) {
	require.NotNil(t, checkDateHeader(""))
	require.NotNil(t, checkDateHeader("not a date"))

	gmtErr := checkDateHeader("Mon, 02 Jan 2006 15:04:05 MST")
	require.NotNil(t, gmtErr)
	require.True(t, strings.HasPrefix(gmtErr.Error(), "Date header must be in GMT"))

	tooLateErr := checkDateHeader("Mon, 02 Jan 2006 15:04:05 GMT")
	require.NotNil(t, tooLateErr)
	require.True(t, strings.Contains(tooLateErr.Error(), "not within max delta"))

	// Alternate format tests
	withTimezoneErr := checkDateHeader(time.Now().Format(time.RFC1123Z))
	require.NotNil(t, withTimezoneErr)

	asRFC3339Err := checkDateHeader(time.Now().Format(time.RFC3339))
	require.NotNil(t, asRFC3339Err)

	// As target format
	err := checkDateHeader(nowInGMT())
	require.NoError(t, err)
}

func TestParseAuthorizationHeader(t *testing.T) {
	key, hmac, err := parseAuthorizationHeader("")
	require.Equal(t, "", key)
	require.Equal(t, []byte{}, hmac)
	require.NotNil(t, err)

	key, hmac, err = parseAuthorizationHeader("Not A Valid Header")
	require.NotNil(t, err)

	key, hmac, err = parseAuthorizationHeader("SomeKey: Not~~Base~~64")
	require.NotNil(t, err)

	key, hmac, err = parseAuthorizationHeader("SomeKey:c3VjY2Vzcw==")
	require.NoError(t, err)
}

func TestAuthenticateAPI(t *testing.T) {
	//set up
	db := initTestDB(t)
	userID := createDummyUser(t, db, models.User{Slug: "slug", FirstName: "fn", LastName: "ln", Email: "normalUser@example.com", Disabled: false})
	keyData := createAPIKey(t, db, userID)

	disabledUser := createDummyUser(t, db, models.User{Slug: "snail", FirstName: "fn", LastName: "ln", Email: "disabledUser@example.com", Disabled: true})
	disabledUsernames := createAPIKey(t, db, disabledUser)

	browser := testBrowser{}
	newReq := func() (*http.Request, io.Reader) {
		_, r := browser.newRequest()
		return r, strings.NewReader("")
	}

	// actual tests
	req, reader := newReq()
	_, badDateErr := authenticateAPI(db, req, reader)
	require.Error(t, badDateErr)

	req, reader = newReq()
	req.Header.Add("Date", nowInGMT())
	_, badAuth := authenticateAPI(db, req, reader)
	require.Error(t, badAuth)

	req, reader = newReq()
	addGoodHeaders(t, req, "badAccessKey", []byte("badSecretKey"))
	_, badKeys := authenticateAPI(db, req, reader)
	require.Error(t, badKeys)

	req, reader = newReq()
	addGoodHeaders(t, req, disabledUsernames.AccessKey, disabledUsernames.SecretKey)
	_, disabledUserError := authenticateAPI(db, req, reader)
	require.Equal(t, backend.DisabledUserError(), disabledUserError)

	req, reader = newReq()
	addGoodHeaders(t, req, keyData.AccessKey, keyData.SecretKey)
	_, shouldWorkErr := authenticateAPI(db, req, reader)
	require.NoError(t, shouldWorkErr)
}

func addGoodHeaders(t *testing.T, r *http.Request, accessKey string, secretKey []byte) {
	r.Header.Add("Date", nowInGMT())

	authorization, err := signer.BuildClientRequestAuthorization(r, accessKey, secretKey)
	require.NoError(t, err)

	r.Header.Add("Authorization", authorization)
}

func nowInGMT() string {
	utcDate := time.Now().In(time.UTC).Format(time.RFC1123)
	return strings.Replace(utcDate, "UTC", "GMT", -1)
}

func initTestDB(t *testing.T) *database.Connection {
	db := database.NewTestConnectionFromNonStandardMigrationPath(t, "middleware-test-db", "../../migrations")
	return db
}

func createDummyUser(t *testing.T, db *database.Connection, usr models.User) int64 {
	userID, err := db.Insert("users", map[string]interface{}{
		"slug":       usr.Slug,
		"first_name": usr.FirstName,
		"last_name":  usr.LastName,
		"email":      usr.Email,
		"disabled":   usr.Disabled,
	})

	require.NoError(t, err)
	return userID
}

func createAPIKey(t *testing.T, db *database.Connection, userID int64) *dtos.APIKey {
	const accessKeyLength = 18
	const secretKeyLength = 64

	accessKey := make([]byte, accessKeyLength)
	_, err := rand.Read(accessKey)
	require.NoError(t, err)

	accessKeyStr := base64.URLEncoding.EncodeToString(accessKey)

	secretKey := make([]byte, secretKeyLength)
	_, err = rand.Read(secretKey)
	require.NoError(t, err)

	_, err = db.Insert("api_keys", map[string]interface{}{
		"user_id":    userID,
		"access_key": accessKeyStr,
		"secret_key": secretKey,
	})

	return &dtos.APIKey{
		AccessKey: accessKeyStr,
		SecretKey: secretKey,
	}
}

// testBrowser generates test requests/responsewriters and saves cookies for all future requests
type testBrowser struct {
	lastResponseRecorder *httptest.ResponseRecorder
	cookies              []*http.Cookie
}

func (b *testBrowser) newRequest() (http.ResponseWriter, *http.Request) {
	// Save cookies from last recorded response
	if b.lastResponseRecorder != nil {
		cookiesToAdd := b.lastResponseRecorder.Result().Cookies()
		for _, cookie := range cookiesToAdd {
			b.cookies = append(b.cookies, cookie)
		}
	}

	r := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()

	// Add all saved cookies to the request
	for _, cookie := range b.cookies {
		r.AddCookie(cookie)
	}

	b.lastResponseRecorder = responseRecorder

	return responseRecorder, r
}
