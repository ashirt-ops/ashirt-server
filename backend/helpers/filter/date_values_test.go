// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package filter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/helpers/filter"
)

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestDateValues(t *testing.T) {
	val1, val2, values := testSetup()

	require.Equal(t, []filter.DateRange{val1, val2}, values.Values())
}

func TestSplitDateValuesByModifier(t *testing.T) {
	val1, val2, values := testSetup()

	require.Equal(t, map[filter.FilterModifier][]filter.DateRange{
		filter.Normal: {val1},
		filter.Not:    {val2},
	}, values.SplitByModifier())
}

func TestSplitDateValues(t *testing.T) {
	val1, val2, values := testSetup()

	require.Equal(t, map[string][]filter.DateRange{
		"in-2020": {val1},
		"outside": {val2},
	}, values.SplitValues(func(v filter.DateValue) string {
		if v.Value.From.Year() == 2020 {
			return "in-2020"
		}
		return "outside"
	}))
}

func TestDateValue(t *testing.T) {
	_, val2, values := testSetup()

	require.Equal(t, val2, values.Value(1))
}

func testSetup() (filter.DateRange, filter.DateRange, filter.DateValues) {
	val1 := filter.DateRange{
		From: date(2020, 1, 10),
		To:   date(2020, 1, 12),
	}
	val2 := filter.DateRange{
		From: date(2021, 4, 1),
		To:   date(2021, 4, 20),
	}
	values := filter.DateValues{
		filter.DateVal(val1),
		filter.NotDateVal(val2),
	}
	return val1, val2, values
}
