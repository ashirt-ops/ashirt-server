// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package dissectors_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	extract "github.com/theparanoids/ashirt/backend/server/dissectors"
)

type dummyRequest struct {
	StringAsString           string    `json:"stringAsString"`
	IntAsNumber              float64   `json:"intAsNumber"`
	IntAsString              string    `json:"intAsString"`
	FloatAsNumber            float64   `json:"floatAsNumber"`
	FloatAsString            string    `json:"floatAsString"`
	IntSliceAsNumberArray    []float64 `json:"intSliceAsNumberArray"`
	StringSliceAsStringArray []string  `json:"stringSliceAsStringArray"`
	BoolAsBool               bool      `json:"boolAsBool"`
	BoolAsString             string    `json:"boolAsString"`
	TimeAsRFC3339            time.Time `json:"timeAsRFC3339"`
	TimeAsRFC3339WithTZ      time.Time `json:"timeAsRFC3339WithTZ"`
}

const rawJSON = `{
	"stringAsString": "someString",
	"intAsNumber": 2048,
	"intAsString": "1024",
	"floatAsNumber": 3.14,
	"floatAsString": "2.71",
	"intSliceAsNumberArray": [4, 8, 15, 16, 23, 42],
	"stringSliceAsStringArray": ["do", "you", "want", "to", "play", "a", "game?"],
	"boolAsBool": true,
	"boolAsString": "true",
	"timeAsRFC3339": "2001-01-31T11:22:33Z",
	"timeAsRFC3339WithTZ": "2001-01-31T11:22:33-07:00"
}`

func TestGetRequiredAfterError(t *testing.T) {
	rawBody, _ := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	str := pq.FromBody("DoesNotExit").Required().AsString()

	assert.NotNil(t, pq.Error, "Should have an error")
	someErr := pq.Error

	str2 := pq.FromBody("OtherEmptyField").Required().AsString()
	otherErr := pq.Error

	assert.Equal(t, str, "")
	assert.Equal(t, str2, "")

	assert.Equal(t, someErr, otherErr)
}

// Verify negative tests

func TestNoBody(t *testing.T) {
	req := makeRequest(nil, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	assert.Nil(t, pq.Error, "Should have no error")
}

func TestNotJSONBody(t *testing.T) {
	notJson := "[}"
	req := makeRequest(&notJson, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	assert.NotNil(t, pq.Error, "Should have a parse error")
}

func TestNoContentAsString(t *testing.T) {
	req := makeRequest(nil, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("DoesNotExist").AsString()
	expected := ""
	assert.Equal(t, expected, actual, "On Unparsable, zero value should be returned")

	assert.Nil(t, pq.Error, "Should have no error")
}

func TestNoContentAsRequiredString(t *testing.T) {
	req := makeRequest(nil, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("DoesNotExist").Required().AsString()
	expected := ""
	assert.Equal(t, expected, actual, "On Unparsable, zero value should be returned")

	assert.NotNil(t, pq.Error, "Should have an error")
}

func TestNoContentAsInt64(t *testing.T) {
	req := makeRequest(nil, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("DoesNotExist").AsInt64()
	expected := int64(0)
	assert.Equal(t, expected, actual, "On Unparsable, zero value should be returned")

	assert.Nil(t, pq.Error, "Should have no error")
}

func TestStringAsInt64(t *testing.T) {
	// designed to fail, because we cannot convert a random string into an integer
	rawBody, _ := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("stringAsString").Required().AsInt64()
	expected := int64(0)
	assert.Equal(t, expected, actual, "On Unparsable, zero value should be returned")

	assert.NotNil(t, pq.Error, "Reqquired values that do not coerce should have an error")
}

func TestNoContentAsRequiredInt64(t *testing.T) {
	req := makeRequest(nil, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("DoesNotExist").Required().AsInt64()
	expected := int64(0)
	assert.Equal(t, expected, actual, "On Unparsable, zero value should be returned")

	assert.NotNil(t, pq.Error, "Should have an error")
}

// Verify Required / Default values

func TestDefaultValue(t *testing.T) {
	rawBody, _ := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	defaultValue := "It's Okay"
	actual := pq.FromBody("DoesNotExist").OrDefault(defaultValue).AsString()
	expected := defaultValue
	assert.Nil(t, pq.Error)
	assert.Equal(t, expected, actual, "Default value should be used")
}

func TestRequiredValue_doesNotExist(t *testing.T) {
	rawBody, _ := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("DoesNotExist").Required().AsString()
	expected := ""
	assert.NotNil(t, pq.Error)
	assert.Equal(t, expected, actual, "On Unparsable, zero value should be returned")
}

// Verify reading/converting from Body

func TestParseStringFromBody(t *testing.T) {
	rawBody, asStruct := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("stringAsString").Required().AsString()
	expected := asStruct.StringAsString
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseIntFromBody(t *testing.T) {
	rawBody, asStruct := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("intAsNumber").Required().AsInt64()
	expected := int64(asStruct.IntAsNumber)
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseIntSliceFromBody(t *testing.T) {
	rawBody, asStruct := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("intSliceAsNumberArray").Required().AsInt64Slice()
	expected := toInt64Slice(asStruct.IntSliceAsNumberArray)
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseStringSliceFromBody(t *testing.T) {
	rawBody, asStruct := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("stringSliceAsStringArray").Required().AsStringSlice()
	expected := asStruct.StringSliceAsStringArray
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseTimeFromBody(t *testing.T) {
	rawBody, asStruct := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("timeAsRFC3339").Required().AsTime()
	expected := asStruct.TimeAsRFC3339
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseFullTimeFromBody(t *testing.T) {
	rawBody, asStruct := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("timeAsRFC3339WithTZ").Required().AsTime()
	expected := asStruct.TimeAsRFC3339WithTZ
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseBoolFromBody(t *testing.T) {
	rawBody, asStruct := prepInputData()
	req := makeRequest(&rawBody, "")
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromBody("boolAsBool").Required().AsBool()
	expected := asStruct.BoolAsBool
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

// Verify reading/converting from Query

func TestParseStringFromQuery(t *testing.T) {
	key, value := "key", "singleValue"
	req := makeRequest(nil, fmt.Sprintf("%v=%v", key, value))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromQuery(key).Required().AsString()
	expected := value
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseIntFromQuery(t *testing.T) {
	key, value := "key", int64(123)
	req := makeRequest(nil, fmt.Sprintf("%v=%v", key, value))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromQuery(key).Required().AsInt64()
	expected := value
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseIntSliceFromQuery(t *testing.T) {
	key, value1 := "key", int64(123)
	key, value2 := "key", int64(456)
	req := makeRequest(nil, fmt.Sprintf("%v=%v&%v=%v", key, value1, key, value2))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromQuery(key).Required().AsInt64Slice()
	expected := []int64{value1, value2}
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseStringSliceFromQuery(t *testing.T) {
	key, value1 := "key", "dog"
	key, value2 := "key", "cat"
	req := makeRequest(nil, fmt.Sprintf("%v=%v&%v=%v", key, value1, key, value2))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromQuery(key).Required().AsStringSlice()
	expected := []string{value1, value2}
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseTimeFromQuery(t *testing.T) {
	key, value1 := "key", "2001-01-31T11:22:33Z"
	req := makeRequest(nil, fmt.Sprintf("%v=%v", key, value1))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromQuery(key).Required().AsTime()
	expected, _ := time.Parse(time.RFC3339, value1)
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseFullTimeFromQuery(t *testing.T) {
	key, value1 := "key", "2001-01-31T11:22:33-07:00"
	req := makeRequest(nil, fmt.Sprintf("%v=%v", key, value1))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromQuery(key).Required().AsTime()
	expected, _ := time.Parse(time.RFC3339, value1)
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseMultipleValuesFromQuery(t *testing.T) {
	key, value1, value2 := "key", int64(123), int64(456)
	altKey, altValue := "mischief", "managed"
	req := makeRequest(nil, fmt.Sprintf("%v=%v&%v=%v&%v=%v", key, value1, key, value2, altKey, altValue))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual1 := pq.FromQuery(key).Required().AsInt64Slice()
	expected1 := []int64{value1, value2}
	assert.Equal(t, expected1, actual1)

	actual2 := pq.FromQuery(altKey).Required().AsString()
	expected2 := altValue
	assert.Equal(t, expected2, actual2)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseBoolFromQuery(t *testing.T) {
	key, value := "key", true
	req := makeRequest(nil, fmt.Sprintf("%v=%v", key, value))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromQuery(key).Required().AsBool()
	expected := value
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseBoolFlagFromQuery(t *testing.T) {
	key := "key"
	req := makeRequest(nil, fmt.Sprintf("%v", key))
	pq := extract.DissectJSONRequest(req, map[string]string{})

	actual := pq.FromQuery(key).Required().AsBool()
	expected := true
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

// Verify reading/converting from URL

func TestParseStringFromURL(t *testing.T) {
	req := makeRequest(nil, "")
	key, value := "key", "value"
	secondKey, secondValue := "key2", "value2"
	pq := extract.DissectJSONRequest(req, makeUrlParamMap(key, value, secondKey, secondValue))

	actual := pq.FromURL(key).Required().AsString()
	expected := value
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseInt64FromURL(t *testing.T) {
	req := makeRequest(nil, "")
	key, value := "key", "12"
	secondKey, secondValue := "key2", "value2"
	pq := extract.DissectJSONRequest(req, makeUrlParamMap(key, value, secondKey, secondValue))

	actual := pq.FromURL(key).Required().AsInt64()
	expected, _ := strconv.ParseInt(value, 10, 64)
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseMultipleValuesFromURL(t *testing.T) {
	req := makeRequest(nil, "")
	key, value := "key", "12"
	key2, value2 := "key2", "value2"
	pq := extract.DissectJSONRequest(req, makeUrlParamMap(key, value, key2, value2))

	actualValue1 := pq.FromURL(key).Required().AsInt64()
	expectedValue1, _ := strconv.ParseInt(value, 10, 64)
	assert.Equal(t, expectedValue1, actualValue1)

	actualValue2 := pq.FromURL(key2).Required().AsString()
	expectedValue2 := value2
	assert.Equal(t, expectedValue2, actualValue2)

	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseTimeFromURL(t *testing.T) {
	key, value := "key", "2001-01-31T11:22:33Z"
	req := makeRequest(nil, "")
	pq := extract.DissectJSONRequest(req, makeUrlParamMap(key, value))

	actual := pq.FromURL(key).Required().AsTime()
	expected, _ := time.Parse(time.RFC3339, value)
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseFullTimeFromURL(t *testing.T) {
	key, value := "key", "2001-01-31T11:22:33-07:00"
	req := makeRequest(nil, "")
	pq := extract.DissectJSONRequest(req, makeUrlParamMap(key, value))

	actual := pq.FromURL(key).Required().AsTime()
	expected, _ := time.Parse(time.RFC3339, value)
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

func TestParseBoolFromURL(t *testing.T) {
	req := makeRequest(nil, "")
	key, value := "key", "true"
	secondKey, secondValue := "key2", "value2"
	pq := extract.DissectJSONRequest(req, makeUrlParamMap(key, value, secondKey, secondValue))

	actual := pq.FromURL(key).Required().AsBool()
	expected, _ := strconv.ParseBool(value)
	assert.Equal(t, expected, actual)
	assert.Nil(t, pq.Error, "Should have no error")
}

// test helpers

func makeRequest(body *string, query string) *http.Request {
	rtn := http.Request{
		URL: makeURL(query),
	}
	if body == nil {
		rtn.Body = http.NoBody
	} else {
		rtn.Body = makeReadCloser(*body)
	}
	return &rtn
}

func makeReadCloser(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}

func makeURL(queryString string) *url.URL {
	rtn := url.URL{
		RawQuery: queryString,
	}
	return &rtn
}

func makeUrlParamMap(values ...string) map[string]string {
	if len(values)%2 == 1 {
		panic("I need an even number of arguments!")
	}
	rtn := make(map[string]string)
	for i := 0; i < len(values); i += 2 {
		rtn[values[i]] = values[i+1]
	}
	return rtn
}

func prepInputData() (string, dummyRequest) {
	var rtn dummyRequest
	err := json.Unmarshal([]byte(rawJSON), &rtn)
	if err != nil {
		panic("Couldn't decode json! " + err.Error())
	}
	return rawJSON, rtn
}

func toInt64Slice(before []float64) []int64 {
	after := make([]int64, len(before))
	for i, v := range before {
		after[i] = int64(v)
	}
	return after
}
