// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package dissectors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Notes:
// Url Parameters are always strings, so a true conversion needs to take place
// Query Parameters are always []string (at least of len 1), so a true conversion takes place either over the first element, or over all elements
// JSON bodies are either float64 (numbers), strings, nil (null), bool, []interface{} (arrays), or map[string]interface{} (objects)
//   For JSON, we don't worry about nested documents or null values
// MultipartForms are always strings, and so effectively behave as UrlParameters

func TestMaybeBoolFromString_valid(t *testing.T) {
	inputValue := "true"
	expected := true
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromString_false_valid(t *testing.T) {
	inputValue := "false"
	expected := false
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromString_invalid(t *testing.T) {
	inputValue := "???"
	expected := false
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.False(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromStringSlice_valid(t *testing.T) {
	inputValue := []string{"true", "shouldNotMatter"}
	expected := true
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromStringSlice_false_valid(t *testing.T) {
	inputValue := []string{"false", "shouldNotMatter"}
	expected := false
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromStringSlice_asFlag(t *testing.T) {
	inputValue := []string{"", "shouldNotMatter"}
	expected := true
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromStringSlice_asFlag_negative(t *testing.T) {
	inputValue := []string{"", "shouldNotMatter"}
	expected := false
	actual, ok := maybeBoolToBool(inputValue, false)

	assert.False(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromStringSlice_invalid(t *testing.T) {
	inputValue := []string{"???", "shouldNotMatter"}
	expected := false
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.False(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromBool_true(t *testing.T) {
	inputValue := true
	expected := true
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromBool_false(t *testing.T) {
	inputValue := false
	expected := false
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.True(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeBoolFromDefault(t *testing.T) {
	inputValue := float64(1.0)
	expected := false
	actual, ok := maybeBoolToBool(inputValue, true)

	assert.False(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeInt64_string_nonint(t *testing.T) {
	inputValue := "1.21"
	var expected int64
	actual, ok := maybeIntToInt64(inputValue, true)

	assert.False(t, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeInt64_string_int(t *testing.T) {
	inputValue := "1"
	expected, expectedOk := int64(1), true
	actual, ok := maybeIntToInt64(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeInt64_stringSlice_int(t *testing.T) {
	inputValue := []string{"1", "shouldNotMatter"}
	expected, expectedOk := int64(1), true
	actual, ok := maybeIntToInt64(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeInt64_stringSlice_nonint(t *testing.T) {
	inputValue := []string{"???", "shouldNotMatter"}
	expected, expectedOk := int64(0), false
	actual, ok := maybeIntToInt64(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeInt64_default(t *testing.T) {
	inputValue := true
	expected, expectedOk := int64(0), false
	actual, ok := maybeIntToInt64(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeInt64_float64_strict(t *testing.T) {
	inputValue := float64(1.21)
	expected, expectedOk := int64(0), false
	actual, ok := maybeIntToInt64(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeInt64_float64_nonstrict(t *testing.T) {
	inputValue := float64(1.21)
	expected, expectedOk := int64(1), true
	actual, ok := maybeIntToInt64(inputValue, false)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeStringSlice_stringSlice(t *testing.T) {
	inputValue := []string{"one", "two", "three"}
	expected, expectedOk := inputValue, true
	actual, ok := maybeStringSliceToStringSlice(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeStringSlice_nonSlice(t *testing.T) {
	inputValue := "tomato"
	expected, expectedOk := []string{}, false
	actual, ok := maybeStringSliceToStringSlice(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeStringSlice_interfaceSlice(t *testing.T) {
	inputValue := []interface{}{"1", "2", "3"}
	expected, expectedOk := []string{"1", "2", "3"}, true
	actual, ok := maybeStringSliceToStringSlice(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeStringSlice_interfaceSlice_mixed(t *testing.T) {
	inputValue := []interface{}{"1", 2, "3"}
	expected, expectedOk := []string{}, false
	actual, ok := maybeStringSliceToStringSlice(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeString_string(t *testing.T) {
	inputValue := "potato"
	expected, expectedOk := inputValue, true
	actual, ok := maybeStringToString(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeString_stringSlice(t *testing.T) {
	inputValue := []string{"one", "two"}
	expected, expectedOk := inputValue[0], true
	actual, ok := maybeStringToString(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeString_nonString(t *testing.T) {
	inputValue := 1234
	expected, expectedOk := "", false
	actual, ok := maybeStringToString(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeString_emptyStringSlice(t *testing.T) {
	inputValue := []string{}
	expected, expectedOk := "", false
	actual, ok := maybeStringToString(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeTime_stringSlice(t *testing.T) {
	timeStr := "2012-11-01T22:08:41+00:00"
	timeAsTime, _ := time.Parse(time.RFC3339, timeStr)

	inputValue := []string{timeStr, "does not matter"}
	expected, expectedOk := timeAsTime, true
	actual, ok := maybeTimeToTime(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeTime_string(t *testing.T) {
	timeStr := "2012-11-01T22:08:41+00:00"
	timeAsTime, _ := time.Parse(time.RFC3339, timeStr)

	inputValue := timeStr
	expected, expectedOk := timeAsTime, true
	actual, ok := maybeTimeToTime(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeTime_nontime(t *testing.T) {
	inputValue := 12
	expected, expectedOk := time.Time{}, false
	actual, ok := maybeTimeToTime(inputValue)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_floatSlice_nonstrict(t *testing.T) {
	inputValue := []float64{3.14, 2.71}
	expected, expectedOk := []int64{3, 2}, true
	actual, ok := maybeIntSliceToInt64Slice(inputValue, false)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_floatSlice_strict(t *testing.T) {
	inputValue := []float64{3.14, 2.71}
	expected, expectedOk := []int64{}, false
	actual, ok := maybeIntSliceToInt64Slice(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_interfaceSlice_strict(t *testing.T) {
	inputValue := []interface{}{3.14, 2.71}
	expected, expectedOk := []int64{}, false
	actual, ok := maybeIntSliceToInt64Slice(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_interfaceSlice_nonstrict(t *testing.T) {
	inputValue := []interface{}{3.14, 2.71}
	expected, expectedOk := []int64{3, 2}, true
	actual, ok := maybeIntSliceToInt64Slice(inputValue, false)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_interfaceSlice(t *testing.T) {
	inputValue := []interface{}{float64(2), float64(3)}
	expected, expectedOk := []int64{2, 3}, true
	actual, ok := maybeIntSliceToInt64Slice(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_interfaceSlice_mixed(t *testing.T) {
	inputValue := []interface{}{float64(2), "3"}
	expected, expectedOk := []int64{}, false
	actual, ok := maybeIntSliceToInt64Slice(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_nonslice(t *testing.T) {
	inputValue := float64(3)
	expected, expectedOk := []int64{}, false
	actual, ok := maybeIntSliceToInt64Slice(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_stringSlice(t *testing.T) {
	inputValue := []string{"2", "3"}
	expected, expectedOk := []int64{2, 3}, true
	actual, ok := maybeIntSliceToInt64Slice(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}

func TestMaybeIntSlice_stringSlice_mixed(t *testing.T) {
	inputValue := []string{"2", "b"}
	expected, expectedOk := []int64{}, false
	actual, ok := maybeIntSliceToInt64Slice(inputValue, true)

	assert.Equal(t, expectedOk, ok)
	assert.Equal(t, expected, actual)
}
