// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package dissectors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCoercableOrDefault(t *testing.T) {
	mock := Coercable{}

	assert.Nil(t, mock.defaultValue)

	value := "tomato"
	mock.OrDefault(value)
	assert.Equal(t, value, mock.defaultValue)
}

func TestCoercableRequired(t *testing.T) {
	mock := Coercable{}

	assert.False(t, mock.required)
	mock.Required()
	assert.Equal(t, true, mock.required)
}

func TestAsBool_isSo(t *testing.T) {
	mock := Coercable{
		rawValue: "true",
	}
	expected := true
	actual := mock.AsBool()

	assert.Equal(t, expected, actual)
}

func TestAsBool_isNotSo(t *testing.T) {
	mock := Coercable{
		rawValue: "7",
	}
	expected := false
	actual := mock.AsBool()

	assert.Equal(t, expected, actual)
}

func TestAsInt64_isSo(t *testing.T) {
	mock := Coercable{
		rawValue: "128",
	}
	expected := int64(128)
	actual := mock.AsInt64()

	assert.Equal(t, expected, actual)
}

func TestAsInt64_isNotSo(t *testing.T) {
	mock := Coercable{
		rawValue: "banana",
	}
	expected := int64(0)
	actual := mock.AsInt64()

	assert.Equal(t, expected, actual)
}

func TestAsInt64_isNotSo_hasDefault(t *testing.T) {
	mock := Coercable{
		rawValue:     "banana",
		defaultValue: int64(42),
	}
	expected := int64(42)
	actual := mock.AsInt64()

	assert.Equal(t, expected, actual)
}

func TestAsString_isSo(t *testing.T) {
	mock := Coercable{
		rawValue: "tofu",
	}
	expected := "tofu"
	actual := mock.AsString()

	assert.Equal(t, expected, actual)
}

func TestAsString_isNotSo(t *testing.T) {
	mock := Coercable{
		rawValue: 9,
	}
	expected := ""
	actual := mock.AsString()

	assert.Equal(t, expected, actual)
}

func TestAsStringSlice_isSo(t *testing.T) {
	mock := Coercable{
		rawValue: []string{"this", "that"},
	}
	expected := []string{"this", "that"}
	actual := mock.AsStringSlice()

	assert.Equal(t, expected, actual)
}

func TestAsStringSlice_isNotSo(t *testing.T) {
	mock := Coercable{
		rawValue: "this",
	}
	var expected []string
	actual := mock.AsStringSlice()

	assert.Equal(t, expected, actual)
}

func TestAsInt64Slice_isSo(t *testing.T) {
	mock := Coercable{
		rawValue: []string{"4", "8", "15"},
	}
	expected := []int64{4, 8, 15}
	actual := mock.AsInt64Slice()

	assert.Equal(t, expected, actual)
}

func TestAsInt64Slice_isNotSo(t *testing.T) {
	mock := Coercable{
		rawValue: 4,
	}
	var expected []int64
	actual := mock.AsInt64Slice()

	assert.Equal(t, expected, actual)
}

func TestAsTime_isSo(t *testing.T) {
	timeStr := "2012-11-01T22:08:41+00:00"
	timeAsTime, _ := time.Parse(time.RFC3339, timeStr)

	mock := Coercable{
		rawValue: timeStr,
	}
	expected := timeAsTime
	actual := mock.AsTime()

	assert.Equal(t, expected, actual)
}

func TestAsTime_isNotSo(t *testing.T) {
	timeStr := "Thursday"

	mock := Coercable{
		rawValue: timeStr,
	}
	expected := time.Time{}
	actual := mock.AsTime()

	assert.Equal(t, expected, actual)
}

func TestAsUnixTime_isSo(t *testing.T) {
	timeAsInt := float64(0xCAFEBABE)
	timeAsTime := time.Unix(0, int64(timeAsInt))

	mock := Coercable{
		rawValue: timeAsInt,
	}
	expected := timeAsTime
	actual := mock.AsUnixTime()

	assert.Equal(t, expected, actual)
}

func TestAsUnixTime_isNotSo(t *testing.T) {
	mock := Coercable{
		rawValue: "Seven o' clock",
	}
	expected := time.Time{}
	actual := mock.AsUnixTime()

	assert.Equal(t, expected, actual)
}
