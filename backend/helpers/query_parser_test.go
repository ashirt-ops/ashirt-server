// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/helpers"
)

func testTimelineQueryCase(t *testing.T, input string, expectedOutput helpers.TimelineFilters) {
	t.Helper()
	actualOutput, err := helpers.ParseTimelineQuery(input)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, actualOutput)
}

func testTimelineQueryExpectErr(t *testing.T, input string) {
	t.Helper()
	_, err := helpers.ParseTimelineQuery(input)
	require.NotNil(t, err)
	require.Error(t, err)
}

func TestParseTimelineQuery(t *testing.T) {
	testTimelineQueryCase(t, "", helpers.TimelineFilters{})
	testTimelineQueryCase(t, "some text string", helpers.TimelineFilters{
		Text: []string{"some", "text", "string"},
	})
	testTimelineQueryCase(t, `Text without quotes "text with quotes"`, helpers.TimelineFilters{
		Text: []string{"Text", "without", "quotes", "text with quotes"},
	})
	testTimelineQueryCase(t, "tag:MyTag", helpers.TimelineFilters{
		Tags: []string{"MyTag"},
	})
	testTimelineQueryCase(t, "tag:MyTag tag:OtherTag", helpers.TimelineFilters{
		Tags: []string{"MyTag", "OtherTag"},
	})
	testTimelineQueryCase(t, `tag:"Tag with spaces"`, helpers.TimelineFilters{
		Tags: []string{"Tag with spaces"},
	})
	testTimelineQueryCase(t, `"Some text" search tag:"First tag" more "text search" tag:SecondTag`, helpers.TimelineFilters{
		Text: []string{"Some text", "search", "more", "text search"},
		Tags: []string{"First tag", "SecondTag"},
	})
	testTimelineQueryCase(t, "Text   with        extra spaces   tag:tag", helpers.TimelineFilters{
		Text: []string{"Text", "with", "extra", "spaces"},
		Tags: []string{"tag"},
	})
	testTimelineQueryCase(t, `operator:alice`, helpers.TimelineFilters{
		Operator: "alice",
	})
	testTimelineQueryCase(t, `Date range example range:2019-05-01,2019-08-05`, helpers.TimelineFilters{
		Text:      []string{"Date", "range", "example"},
		DateRange: &helpers.DateRange{time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 8, 5, 23, 59, 59, 0, time.UTC)},
	})
	testTimelineQueryCase(t, `Time range example range:2019-05-01T08:00:00Z,2019-08-05T19:30:00Z`, helpers.TimelineFilters{
		Text:      []string{"Time", "range", "example"},
		DateRange: &helpers.DateRange{time.Date(2019, 5, 1, 8, 0, 0, 0, time.UTC), time.Date(2019, 8, 5, 19, 30, 0, 0, time.UTC)},
	})
	testTimelineQueryCase(t, `uuid:00000000-1234-5678-ABCD-000000000000`, helpers.TimelineFilters{
		UUID: "00000000-1234-5678-ABCD-000000000000",
	})
	testTimelineQueryCase(t, `with-evidence:00000000-1234-5678-ABCD-000000000000`, helpers.TimelineFilters{
		WithEvidenceUUID: "00000000-1234-5678-ABCD-000000000000",
	})

	True := true
	False := false
	testTimelineQueryCase(t, `linked:true`, helpers.TimelineFilters{
		Linked: &True,
	})
	testTimelineQueryCase(t, `linked:false`, helpers.TimelineFilters{
		Linked: &False,
	})
	testTimelineQueryCase(t, `linked:all`, helpers.TimelineFilters{
		Linked: nil,
	})
	testTimelineQueryCase(t, `sort:asc`, helpers.TimelineFilters{
		SortAsc: true,
	})
	testTimelineQueryCase(t, `sort:chronological`, helpers.TimelineFilters{
		SortAsc: true,
	})
	testTimelineQueryCase(t, `sort:ascending`, helpers.TimelineFilters{
		SortAsc: true,
	})

	testTimelineQueryCase(t, `sort:desc`, helpers.TimelineFilters{
		SortAsc: false,
	})
	testTimelineQueryCase(t, ``, helpers.TimelineFilters{
		SortAsc: false,
	})

	testTimelineQueryExpectErr(t, `invalid keys cause error invalid:value`)
	testTimelineQueryExpectErr(t, `multiple operators       cause error operator:alice operator:bob`)
	testTimelineQueryExpectErr(t, `multiple uuids           cause error uuid:ABC123 uuid:XYZ789`) // actual uuid doesn't currently matter
	testTimelineQueryExpectErr(t, `multiple with_evidence   cause error with-evidence:ABC123 with-evidence:XYZ789`)
	testTimelineQueryExpectErr(t, `multiple linked          cause error linked:all linked:true`)
	testTimelineQueryExpectErr(t, `multiple sort_directions cause error sort:desc sort:asc`)
	testTimelineQueryExpectErr(t, `unparsable bool/not all cause error linked:maybe`)
	testTimelineQueryExpectErr(t, `unparsable date cause error range:2021-01-01,2021-02-31`)
	testTimelineQueryExpectErr(t, `unparsable date cause error (alt) range:2021-01-01`)
}
