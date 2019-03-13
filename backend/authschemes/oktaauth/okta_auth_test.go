package oktaauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/theparanoids/ashirt/backend/authschemes"
	"github.com/theparanoids/ashirt/backend/config"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/integration"
	"github.com/theparanoids/ashirt/backend/session"

	sq "github.com/Masterminds/squirrel"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/stretchr/testify/require"
)

var generalOktaConfig = config.AuthInstanceConfig{
	ClientID:                 "magicClientID",
	ClientSecret:             "magicClientSecret",
	Issuer:                   "http://me.com/oauth2/default",
	BackendURL:               "http://localhost:8080/web",
	SuccessRedirectURL:       "http://localhost:8080",
	FailureRedirectURLPrefix: "http://localhost:8080",
	ProfileToShortnameField:  "preferred_username",
}

var everyone = func(map[string]string) bool { return true }

func dynamicOktaConfig(remoteServer string) config.AuthInstanceConfig {
	return config.AuthInstanceConfig{
		ClientID:                 "magicClientID",
		ClientSecret:             "magicClientSecret",
		Issuer:                   remoteServer,
		BackendURL:               "http://localhost:8080/web",
		SuccessRedirectURL:       "http://localhost:8080",
		FailureRedirectURLPrefix: "http://localhost:8080",
		ProfileToShortnameField:  "preferred_username",
	}
}

func TestNewFromConfig(t *testing.T) {
	config := generalOktaConfig
	inst := NewFromConfig(config, func(map[string]string) bool { return true })
	require.Equal(t, config.ClientID, inst.clientID)
	require.Equal(t, config.ClientSecret, inst.clientSecret)
	require.Equal(t, config.Issuer, inst.issuer)
	require.Equal(t, config.BackendURL, inst.absoluteBackendPath)
	require.Equal(t, config.SuccessRedirectURL, inst.authSuccessRedirectPath)
	require.Equal(t, config.FailureRedirectURLPrefix, inst.authFailureRedirectPathPrefix)
	require.Equal(t, config.ProfileToShortnameField, inst.profileToShortnameField)
}

func TestAuthDone(t *testing.T) {
	browser := integration.TestBrowser{}

	// Test no error flow
	w, r := browser.NewRequest()
	expectedPath := "/go/to/here"
	rtnVal, rtnErr := authDone(w, r, expectedPath, nil)
	require.NoError(t, rtnErr)
	require.Nil(t, rtnVal)
	verifyRedirect(t, browser, url.URL{Path: expectedPath})

	// test error flow
	w, r = browser.NewRequest()
	expectedPath = "/go/to/there"
	expectedError := errors.New("Oops")
	rtnVal, rtnErr = authDone(w, r, expectedPath, expectedError)
	require.Equal(t, expectedError, rtnErr)
	require.Nil(t, rtnVal)
	verifyRedirect(t, browser, url.URL{Path: expectedPath})
}

func TestAuthFailure(t *testing.T) {
	inst := NewFromConfig(generalOktaConfig, everyone)
	browser := integration.TestBrowser{}
	w, r := browser.NewRequest()
	expectedPathSegment := "/bad"
	expectedError := errors.New("Oops")
	_, rtnErr := inst.authFailure(w, r, expectedError, expectedPathSegment)
	expectedURL, err := url.Parse(inst.authFailureRedirectPathPrefix + expectedPathSegment)
	require.NoError(t, err)
	require.Equal(t, expectedError, rtnErr)
	verifyRedirect(t, browser, *expectedURL)
}

func TestAuthSuccess(t *testing.T) {
	inst := NewFromConfig(generalOktaConfig, everyone)
	browser := integration.TestBrowser{}

	// no linking
	w, r := browser.NewRequest()
	_, rtnErr := inst.authSuccess(w, r, false)
	expectedURL, err := url.Parse(inst.authSuccessRedirectPath)
	require.NoError(t, rtnErr)
	require.NoError(t, err)
	verifyRedirect(t, browser, *expectedURL)

	// with linking
	w, r = browser.NewRequest()
	_, rtnErr = inst.authSuccess(w, r, true)
	require.NoError(t, rtnErr)
	verifyRedirect(t, browser, url.URL{Path: "/account/authmethods"})

}

func TestCallbackURI(t *testing.T) {
	inst := NewFromConfig(generalOktaConfig, everyone)
	result := inst.callbackURI()
	require.Equal(t, inst.absoluteBackendPath+"/auth/okta/callback", result)
}

func TestMakeUserProfile(t *testing.T) {
	inst := NewFromConfig(generalOktaConfig, everyone)
	slug := "Wow"
	firstName := "Such"
	lastName := "Testing"
	email := "Very@Secure.com"
	expectedProfile := map[string]string{
		"given_name":                 firstName,
		"family_name":                lastName,
		"email":                      email,
		inst.profileToShortnameField: slug,
	}
	result := inst.makeUserProfile(expectedProfile)
	require.Equal(t, firstName, result.FirstName)
	require.Equal(t, lastName, result.LastName)
	require.Equal(t, slug, result.Slug)
	require.Equal(t, email, result.Email)
}

func TestRedirectLogin(t *testing.T) {
	inst := NewFromConfig(generalOktaConfig, everyone)
	browser := integration.TestBrowser{}
	bridge := initBridge(t)

	w, r := browser.NewRequest()
	expectedMode := "AnyMode"
	inst.redirectLogin(w, r, bridge, expectedMode)
	sess, ok := bridge.ReadAuthSchemeSession(r).(*preLoginAuthSession)
	require.True(t, ok)

	require.Equal(t, expectedMode, sess.OktaMode)

	expectedURL, _ := url.Parse(inst.issuer + "/v1/authorize")
	expectedValues := url.Values{
		"nonce":        []string{sess.Nonce},
		"state":        []string{sess.StateChallengeCSRF},
		"client_id":    []string{inst.clientID},
		"redirect_uri": []string{inst.callbackURI()},
	}
	expectedURL.RawQuery = expectedValues.Encode()

	verifyRedirect(t, browser, *expectedURL)
}

func TestExchangeCode(t *testing.T) {
	// fixtures
	magicCode := "magicCode"
	goodExchange := Exchange{
		AccessToken: "abc123",
		TokenType:   "masterToken",
		ExpiresIn:   1,
		Scope:       "openid profile email",
		IDToken:     "idToken",
	}
	badFormatExchange := `{
				"errorCode": "E0000021",
				"errorSummary": "Bad request.  Accept and/or Content-Type headers likely do not match supported values.",
				"errorLink": "E0000021",
				"errorId": "???",
				"errorCauses": []
			}`

	// set up
	browser := integration.TestBrowser{}
	s := http.NewServeMux()
	testServer := httptest.NewServer(s)
	inst := NewFromConfig(dynamicOktaConfig(testServer.URL), everyone)

	s.HandleFunc("/v1/token", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		values := r.URL.Query()
		require.Equal(t, "authorization_code", values.Get("grant_type"))
		require.Equal(t, magicCode, values.Get("code"))
		require.Equal(t, inst.callbackURI(), values.Get("redirect_uri"))

		// determine how to exit...
		errCode := values.Get("ctrl")
		//... no-error
		switch errCode {
		case "badFormat":
			w.Write([]byte(badFormatExchange))
			w.WriteHeader(400)
		default:
			encodedBody, err := json.Marshal(goodExchange)
			require.NoError(t, err)
			w.Write(encodedBody)
			w.WriteHeader(200)
		}
		return
	})

	// test 1: Exchange works
	_, testRequest := browser.NewRequest()

	exchange := inst.exchangeCode(magicCode, testRequest)
	require.Equal(t, goodExchange.AccessToken, exchange.AccessToken)
	require.Equal(t, goodExchange.TokenType, exchange.TokenType)
	require.Equal(t, goodExchange.ExpiresIn, exchange.ExpiresIn)
	require.Equal(t, goodExchange.Scope, exchange.Scope)
	require.Equal(t, goodExchange.IDToken, exchange.IDToken)
	require.NoError(t, exchange.WrappedError)

	// test 2: Exchange fails (unexpected format)
	_, testRequest = browser.NewRequest()
	query := testRequest.URL.Query()
	query.Add("ctrl", "badFormat")
	testRequest.URL.RawQuery = query.Encode()
	exchange = inst.exchangeCode(magicCode, testRequest)
	require.Error(t, exchange.WrappedError)
}

func TestGetProfileData(t *testing.T) {
	// set up
	browser := integration.TestBrowser{}
	s := http.NewServeMux()
	testServer := httptest.NewServer(s)
	inst := NewFromConfig(dynamicOktaConfig(testServer.URL), everyone)
	expectedProfile := map[string]string{
		"User": "Profile",
	}

	s.HandleFunc("/v1/userinfo", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "GET", r.Method)
		require.Equal(t, "application/json", r.Header.Get("Accept"))
		require.Contains(t, r.Header.Get("Authorization"), "Bearer ")
		require.Greater(t, len(r.Header.Get("Authorization")), len("Bearer "))

		content, err := json.Marshal(expectedProfile)
		require.NoError(t, err)
		w.Write(content)
		w.WriteHeader(200)
		return
	})

	// verify success
	_, testRequest := browser.NewRequest()
	foundProfile := inst.getProfileData(testRequest, "abc123")
	require.Equal(t, expectedProfile, foundProfile)

	// verify no-access token profiled
	_, testRequest = browser.NewRequest()
	foundProfile = inst.getProfileData(testRequest, "")
	require.Equal(t, map[string]string{}, foundProfile)

}

func TestHandleOktaCallback(t *testing.T) {
	firstNameNotBob := func(values map[string]string) bool { return values["given_name"] != "Bob" }
	inst := NewFromConfig(generalOktaConfig, firstNameNotBob)
	browser := integration.TestBrowser{}
	bridge := initBridge(t)
	nonce := "idToken"

	type reqInput struct {
		// opting for negatives to keep the normal case just an empty structure
		omitQuery   bool
		omitSession bool
		linkAccount bool
	}

	newRequest := func(i reqInput) (http.ResponseWriter, *http.Request) {
		challenge := "challenege"
		code := "abc123"
		w, r := browser.NewRequest()
		if !i.omitQuery {
			q := r.URL.Query()
			q.Add("code", code)
			q.Add("state", challenge)
			r.URL.RawQuery = q.Encode()
		}
		if !i.omitSession {
			mode := modeLogin
			if i.linkAccount {
				mode = modeLink
			}
			bridge.SetAuthSchemeSession(w, r, &preLoginAuthSession{
				Nonce:              nonce,
				StateChallengeCSRF: challenge,
				OktaMode:           mode,
			})
		}

		return w, r
	}
	makeURL := func(errorBase bool, path string) url.URL {
		base := generalOktaConfig.SuccessRedirectURL
		if errorBase {
			base = generalOktaConfig.FailureRedirectURLPrefix
		}
		if path == "/account/authmethods" {
			base = "" // total hack, but it makes the tests read better below
		}
		parsed, err := url.Parse(base + path)
		require.NoError(t, err)
		return *parsed
	}

	goodExchange := func(code string, req *http.Request) Exchange {
		return Exchange{AccessToken: "abc123", TokenType: "masterToken", ExpiresIn: 1, Scope: "openid profile email", IDToken: nonce}
	}
	wrappedErrExchange := func(code string, r *http.Request) Exchange { return Exchange{WrappedError: errors.New("something")} }
	inherentErrExchange := func(code string, r *http.Request) Exchange { return Exchange{Error: "something else"} }

	validTokener := func(token, nonce string) (*verifier.Jwt, error) { return nil, nil }
	invalidTokener := func(token, nonce string) (*verifier.Jwt, error) { return nil, errors.New("placeholder") }

	goodProfiler := func(req *http.Request, accessToken string) map[string]string {
		return map[string]string{"given_name": "All", "family_name": "Goode", "email": "agoode@notbad", "preferred_username": "Goodie"}
	}
	altGoodProfiler := func(req *http.Request, accessToken string) map[string]string {
		return map[string]string{"given_name": "Still", "family_name": "Goode", "email": "sgoode@notbad", "preferred_username": "Oldie"}
	}
	disabledProfiler := func(req *http.Request, accessToken string) map[string]string {
		return map[string]string{"given_name": "Not", "family_name": "Goode", "email": "ngoode@notbad", "preferred_username": "Moldie"}
	}
	tooSmallProfiler := func(req *http.Request, accessToken string) map[string]string {
		return map[string]string{"given_name": "Maybe", "preferred_username": "Maybe_Not"}
	}
	emptyProfiler := func(req *http.Request, accessToken string) map[string]string { return map[string]string{} }
	noBobsTestProfiler := func(req *http.Request, accessToken string) map[string]string {
		return map[string]string{"given_name": "Bob", "preferred_username": "Bobbo"}
	}

	// verify access tests
	testScenario := func(input reqInput, expectError bool, exchangeCode codeExchanger, verifyToken tokenVerifier, getUserProfile profileGatherer, expectedRedirect string) {
		w, r := newRequest(input)
		_, err := inst.handleOktaCallback(w, r, bridge, exchangeCode, verifyToken, getUserProfile)
		if expectError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		verifyRedirect(t, browser, makeURL(expectError, expectedRedirect))
	}
	testScenario(reqInput{omitQuery: true}, true, wrappedErrExchange, validTokener, goodProfiler, "/autherror/noverify")
	testScenario(reqInput{omitSession: true}, true, wrappedErrExchange, validTokener, goodProfiler, "/autherror/noverify")
	testScenario(reqInput{}, true, wrappedErrExchange, validTokener, goodProfiler, "/autherror/noverify")
	testScenario(reqInput{}, true, inherentErrExchange, validTokener, goodProfiler, "/autherror/noverify")
	testScenario(reqInput{}, true, goodExchange, invalidTokener, goodProfiler, "/autherror/noverify")
	testScenario(reqInput{}, true, goodExchange, validTokener, emptyProfiler, "/autherror/noaccess")
	testScenario(reqInput{}, true, goodExchange, validTokener, noBobsTestProfiler, "/autherror/noaccess")
	testScenario(reqInput{}, true, goodExchange, validTokener, tooSmallProfiler, "/autherror/incomplete")
	testScenario(reqInput{}, false, goodExchange, validTokener, goodProfiler, "") // Create a new user!

	// link a user
	newProfile := authschemes.UserProfile{
		FirstName: "Newie",
		LastName:  "User",
		Slug:      altGoodProfiler(nil, "")["preferred_username"],
		Email:     "newie@user",
	}
	makeFullAcount(t, bridge, newProfile)
	testScenario(reqInput{linkAccount: true}, false, goodExchange, validTokener, altGoodProfiler, "/account/authmethods") // link user

	// verify disabled users have no access
	disabledSlug := disabledProfiler(nil, "")["preferred_username"]
	makeFullAcount(t, bridge, authschemes.UserProfile{
		FirstName: "Disabled",
		LastName:  "User",
		Slug:      disabledSlug,
		Email:     "disabled@user",
	})
	bridge.GetDatabase().Update(sq.Update("users").Set("disabled", true).Where(sq.Eq{"slug": disabledSlug}))
	testScenario(reqInput{}, true, goodExchange, validTokener, disabledProfiler, "/autherror/disabled")
}

func makeFullAcount(t *testing.T, bridge authschemes.AShirtAuthBridge, newProfile authschemes.UserProfile) {
	output, err := bridge.CreateNewUser(newProfile)
	require.NoError(t, err)
	require.Equal(t, newProfile.Slug, output.RealSlug)
	bridge.CreateNewAuthForUser(authschemes.UserAuthData{UserID: output.UserID, UserKey: output.RealSlug})
}

func verifyRedirect(t *testing.T, browser integration.TestBrowser, expectedURL url.URL) {
	require.Equal(t, http.StatusFound, browser.LastResponseRecorder.Code)
	location := browser.LastResponseRecorder.HeaderMap.Get("Location")
	parsedURL, err := url.Parse(location)
	parsedQuery, err := url.ParseQuery(parsedURL.RawQuery)
	expectedQuery, _ := url.ParseQuery(expectedURL.RawQuery)
	require.NoError(t, err)
	require.Equal(t, expectedURL.Scheme, parsedURL.Scheme)
	require.Equal(t, expectedURL.Path, parsedURL.Path)
	require.Equal(t, expectedURL.Host, parsedURL.Host)

	for key, expectedValue := range expectedQuery {
		parsedValue, ok := parsedQuery[key]
		require.True(t, ok)
		require.Equal(t, parsedValue, expectedValue)
	}
}

func initBridge(t *testing.T) authschemes.AShirtAuthBridge {
	db := database.NewTestConnectionFromNonStandardMigrationPath(t, "okta-test-db", "../../migrations")
	sessionStore, err := session.NewStore(db, session.StoreOptions{SessionDuration: time.Hour, Key: []byte{}})
	require.NoError(t, err)
	return authschemes.MakeAuthBridge(db, sessionStore, "test")
}
