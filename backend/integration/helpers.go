// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend/authschemes"
	"github.com/ashirt-ops/ashirt-server/backend/authschemes/localauth"
	"github.com/ashirt-ops/ashirt-server/backend/contentstore"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/dtos"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/server"
	"github.com/ashirt-ops/ashirt-server/signer"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type Tester struct {
	t           *testing.T
	s           *httptest.Server
	DefaultUser *UserSession
}

func NewTester(t *testing.T) *Tester {
	db := database.NewTestConnection(t, "integration-test-db")

	doMinimalSeed(db)

	contentStore, err := contentstore.NewDevStore()
	require.NoError(t, err)
	commonLogger := logging.SetupStdoutLogging()

	s := chi.NewRouter()

	s.Route("/web", func(r chi.Router) {
		server.Web(r,
			db, contentStore, &server.WebConfig{
				CSRFAuthKey:     []byte("csrf-auth-key-for-integration-tests"),
				SessionStoreKey: []byte("session-store-key-for-integration-tests"),
				AuthSchemes: []authschemes.AuthScheme{localauth.LocalAuthScheme{
					RegistrationEnabled: true,
				}},
				Logger: commonLogger,
			},
		)
	})

	s.Route("/api", func(r chi.Router) {
		server.API(r,
			db, contentStore, commonLogger,
		)
	})

	return &Tester{
		t: t,
		s: httptest.NewServer(s),
	}
}

func doMinimalSeed(db *database.Connection) {
	commonFindingCategories := []string{
		"Product",
		"Network",
		"Enterprise",
		"Vendor",
		"Behavioral",
		"Detection Gap",
	}
	db.BatchInsert("finding_categories", len(commonFindingCategories), func(i int) map[string]interface{} {
		return map[string]interface{}{
			"category": commonFindingCategories[i],
		}
	})
}

type UserSession struct {
	Client    *http.Client
	CSRFToken string
	UserSlug  string
}

type APIKey struct {
	AccessKey string `json:"accessKey"`
	SecretKey []byte `json:"secretKey"`
}

func (a *Tester) NewUser(slug string, firstName string, lastName string) *UserSession {
	a.t.Helper()

	jar, err := cookiejar.New(nil)
	require.NoError(a.t, err)
	session := &UserSession{Client: &http.Client{Jar: jar}}

	a.Get("/web/user").AsUser(session).Do()

	a.Post("/web/auth/local/register").AsUser(session).WithMarshaledJSONBody(map[string]interface{}{
		"firstName": firstName,
		"lastName":  lastName,
		"username":  slug,
		"email":     strings.ToLower(slug) + "@example.com",
		"password":  "password",
	}).Do().ExpectSuccess()
	a.Post("/web/auth/local/login").AsUser(session).WithMarshaledJSONBody(map[string]interface{}{
		"username": slug,
		"password": "password",
	}).Do().ExpectSuccess()

	profileBytes := a.Get("/web/user").AsUser(session).Do().ResponseBody()
	var profile dtos.UserOwnView
	err = json.Unmarshal(profileBytes, &profile)

	require.NoError(a.t, err)

	session.UserSlug = profile.Slug
	return session
}

type RequestBuilder struct {
	t           *testing.T
	req         *http.Request
	userSession *UserSession
	apiKey      *APIKey
}

func (a *Tester) Delete(path string) *RequestBuilder { a.t.Helper(); return a.buildReq("DELETE", path) }
func (a *Tester) Get(path string) *RequestBuilder    { a.t.Helper(); return a.buildReq("GET", path) }
func (a *Tester) Patch(path string) *RequestBuilder  { a.t.Helper(); return a.buildReq("PATCH", path) }
func (a *Tester) Post(path string) *RequestBuilder   { a.t.Helper(); return a.buildReq("POST", path) }
func (a *Tester) Put(path string) *RequestBuilder    { a.t.Helper(); return a.buildReq("PUT", path) }

func (a *Tester) TestingT() *testing.T {
	return a.t
}

func (a *Tester) buildReq(method string, path string) *RequestBuilder {
	a.t.Helper()
	reqURL, err := url.Parse(a.s.URL + path)
	require.NoError(a.t, err)

	rb := &RequestBuilder{
		t: a.t,
		req: &http.Request{
			Method: method,
			URL:    reqURL,
			Header: http.Header{},
		},
	}

	if a.DefaultUser != nil {
		rb = rb.AsUser(a.DefaultUser)
	}

	return rb
}

func (a *Tester) APIKeyForUser(u *UserSession) *APIKey {
	a.t.Helper()
	url := fmt.Sprintf("/web/user/%v/apikeys", u.UserSlug)
	body := a.Post(url).AsUser(u).Do().ResponseBody()
	var key APIKey
	require.NoError(a.t, json.Unmarshal(body, &key), "Failed to unmarshal response JSON")
	return &key
}

func (b *RequestBuilder) AsUser(u *UserSession) *RequestBuilder {
	b.userSession = u
	b.req.Header.Set("X-CSRF-Token", u.CSRFToken)
	return b
}

func (b *RequestBuilder) WithAPIKey(k *APIKey) *RequestBuilder {
	b.req.Header.Set("Date", time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123))
	b.apiKey = k
	return b
}

func (b *RequestBuilder) WithMarshaledJSONBody(body map[string]interface{}) *RequestBuilder {
	b.t.Helper()
	marshaledBody, err := json.Marshal(body)
	require.NoError(b.t, err)
	return b.WithJSONBody(string(marshaledBody))
}

func (b *RequestBuilder) WithJSONBody(body string) *RequestBuilder {
	b.t.Helper()
	b.req.Header.Add("Content-Type", "application/json")
	b.req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(body)), nil }
	b.req.Body, _ = b.req.GetBody()
	return b
}

func (b *RequestBuilder) WithMultipartBody(fields map[string]string, files map[string]*os.File) *RequestBuilder {
	b.t.Helper()
	body := &bytes.Buffer{}
	mp := multipart.NewWriter(body)
	for k, v := range fields {
		err := mp.WriteField(k, v)
		require.NoError(b.t, err)
	}
	for k, v := range files {
		f, err := mp.CreateFormFile(k, v.Name())
		require.NoError(b.t, err)
		io.Copy(f, v)
	}
	require.NoError(b.t, mp.Close())
	b.req.Header.Add("Content-Type", mp.FormDataContentType())
	b.req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body.Bytes())), nil }
	b.req.Body, _ = b.req.GetBody()
	return b
}

func (b *RequestBuilder) WithURLEncodedBody(fields map[string]string) *RequestBuilder {
	b.t.Helper()
	b.req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	data := url.Values{}
	for k, v := range fields {
		data.Set(k, v)
	}
	b.req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(data.Encode())), nil }
	b.req.Body, _ = b.req.GetBody()
	return b
}

func (b *RequestBuilder) Do() *ResponseTester {
	b.t.Helper()
	client := &http.Client{}
	if b.userSession != nil {
		client = b.userSession.Client
	}

	if b.apiKey != nil {
		authorization, err := signer.BuildClientRequestAuthorization(b.req, b.apiKey.AccessKey, b.apiKey.SecretKey)
		require.NoError(b.t, err)
		b.req.Header.Set("Authorization", authorization)
	}

	res, err := client.Do(b.req)
	require.NoError(b.t, err)

	returnedCSRFToken := res.Header.Get("X-CSRF-Token")
	if b.userSession != nil && returnedCSRFToken != "" {
		b.userSession.CSRFToken = returnedCSRFToken
	}

	return &ResponseTester{t: b.t, res: res}
}

type ResponseTester struct {
	t    *testing.T
	res  *http.Response
	body []byte
}

func (rt *ResponseTester) ExpectSuccess() *ResponseTester {
	rt.t.Helper()
	require.Conditionf(
		rt.t,
		func() bool { return rt.res.StatusCode >= 200 && rt.res.StatusCode < 300 },
		"Expected status code to be success but got %d", rt.res.StatusCode,
	)
	return rt
}

func (rt *ResponseTester) ExpectUnauthorized() *ResponseTester {
	rt.t.Helper()
	require.Equal(rt.t, "401 Unauthorized", rt.res.Status)
	require.JSONEq(rt.t, `{"error": "Unauthorized"}`, string(rt.ResponseBody()))
	return rt
}

func (rt *ResponseTester) ExpectNotFound() *ResponseTester {
	rt.t.Helper()
	rt.ExpectStatus(404)
	return rt
}

func (rt *ResponseTester) ExpectStatus(statusCode int) *ResponseTester {
	rt.t.Helper()
	require.Equal(rt.t, statusCode, rt.res.StatusCode)
	return rt
}

var timestampRegexp = regexp.MustCompile(`"2\d{3}-(?:0\d|1[012])-(?:[012]\d|3[01])T(?:[01]\d|2[0-3]):[0-6]\d:[0-6]\d(?:\.\d+)?(?:Z|[\d-:]+)"`)

func (rt *ResponseTester) ExpectJSON(expected string) *ResponseTester {
	rt.t.Helper()
	rt.ExpectSuccess()

	actual := string(rt.ResponseBody())
	if regexp.MustCompile(`"_TIMESTAMP_"`).MatchString(expected) {
		actual = timestampRegexp.ReplaceAllString(actual, `"_TIMESTAMP_"`)
	}
	require.JSONEq(rt.t, expected, actual, fmt.Sprintf("Expected json %s but got %s", expected, actual))
	return rt
}

func (rt *ResponseTester) ExpectSubsetJSON(expected string) *ResponseTester {
	rt.t.Helper()
	rt.ExpectSuccess()

	actual := map[string]interface{}{}
	require.NoError(rt.t, json.Unmarshal(rt.ResponseBody(), &actual), "Failed to unmarshal response JSON")
	requireSubsetJSONEqual(rt.t, expected, actual, "Subset json does not match")

	return rt
}

func (rt *ResponseTester) ExpectSubsetJSONArray(expected []string) *ResponseTester {
	rt.t.Helper()
	rt.ExpectSuccess()

	actual := []map[string]interface{}{}
	require.NoError(rt.t, json.Unmarshal(rt.ResponseBody(), &actual), "Failed to unmarshal response JSON")
	require.Equal(rt.t, len(actual), len(expected))
	for i := range expected {
		requireSubsetJSONEqual(rt.t, expected[i], actual[i], fmt.Sprintf("Subset json in array index %d does not match", i))
	}

	return rt
}

func (rt *ResponseTester) ExpectResponse(expectedStatus int, expectedBody []byte) *ResponseTester {
	rt.t.Helper()
	require.Equal(rt.t, expectedStatus, rt.res.StatusCode)
	require.Equal(rt.t, expectedBody, rt.ResponseBody())
	return rt
}

func (rt *ResponseTester) ResponseBody() []byte {
	rt.t.Helper()
	if rt.body == nil {
		body, err := io.ReadAll(rt.res.Body)
		require.NoError(rt.t, err)
		rt.body = body
	}
	return rt.body
}

func (rt *ResponseTester) ResponseUUID() string {
	rt.t.Helper()
	var responseObject struct {
		UUID string `json:"uuid"`
	}
	require.NoError(rt.t, json.Unmarshal(rt.ResponseBody(), &responseObject), "Failed to unmarshal response JSON")
	return responseObject.UUID
}

func (rt *ResponseTester) ResponseID() int64 {
	rt.t.Helper()
	var responseObject struct {
		ID int64 `json:"id"`
	}
	require.NoError(rt.t, json.Unmarshal(rt.ResponseBody(), &responseObject), "Failed to unmarshal response JSON")
	return responseObject.ID
}

func requireSubsetJSONEqual(t *testing.T, expectedStr string, actual map[string]interface{}, message string) {
	t.Helper()

	expected := map[string]interface{}{}
	require.NoError(t, json.Unmarshal([]byte(expectedStr), &expected), "Failed to unmarshal expected JSON. Make sure the JSON in the test case is valid")
	for k := range expected {
		require.Equal(t, expected[k], actual[k], `%s: Values for key "%s" differ`, message, k)
	}
}

// TestBrowser generates test requests/responsewriters and saves cookies for all future requests
type TestBrowser struct {
	LastResponseRecorder *httptest.ResponseRecorder
	cookies              []*http.Cookie
}

func (b *TestBrowser) NewRequest() (http.ResponseWriter, *http.Request) {
	// Save cookies from last recorded response
	if b.LastResponseRecorder != nil {
		cookiesToAdd := b.LastResponseRecorder.Result().Cookies()
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

	b.LastResponseRecorder = responseRecorder

	return responseRecorder, r
}
