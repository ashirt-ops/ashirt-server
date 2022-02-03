// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
)

// DateRange is a simple struct representing a slice of time From a point To a point
type DateRange struct {
	From time.Time
	To   time.Time
}

// TimelineFilters represents all of the parsed timeline configuraions
type TimelineFilters struct {
	UUID             string
	Text             []string
	Tags             []string
	Type             string
	Operator         string
	DateRange        *DateRange
	WithEvidenceUUID string
	Linked           *bool
	SortAsc          bool
}

// ParseTimelineQuery parses a query a user may type into the search box on the timeline page
// into a TimelineFilters struct that the events/evidence services expect
func ParseTimelineQuery(query string) (TimelineFilters, error) {
	timelineFilters := TimelineFilters{}

	for k, v := range tokenizeTimelineQuery(query) {
		switch k {
		case "":
			timelineFilters.Text = v
		case "tag":
			timelineFilters.Tags = v
		case "operator":
			if len(v) != 1 {
				errReason := "Only one operator can be specified"
				return timelineFilters, backend.BadInputErr(errors.New(errReason), errReason)
			}
			timelineFilters.Operator = v[0]
		case "range":
			dateRange, err := parseRangeQuery(v)
			if err != nil {
				return timelineFilters, err
			}
			timelineFilters.DateRange = dateRange
		case "uuid":
			if len(v) != 1 {
				errReason := "Only one uuid can be specified"
				return timelineFilters, backend.BadInputErr(errors.New(errReason), errReason)
			}
			timelineFilters.UUID = v[0]
		case "with-evidence":
			if len(v) != 1 {
				errReason := "Only one with-evidence can be specified"
				return timelineFilters, backend.BadInputErr(errors.New(errReason), errReason)
			}
			timelineFilters.WithEvidenceUUID = v[0]
		case "linked":
			if len(v) != 1 {
				errReason := "Linked can only be specified once"
				return timelineFilters, backend.BadInputErr(errors.New(errReason), errReason)
			}
			if strings.ToLower(v[0]) == "all" {
				continue // providing a 3rd option here to allow for easy filter-removal
			}
			val, err := strconv.ParseBool(v[0])
			if err != nil {
				errReason := "Linked value must be True or False"
				return timelineFilters, backend.BadInputErr(errors.New(errReason), errReason)
			}
			timelineFilters.Linked = &val
		case "sort":
			if len(v) != 1 {
				errReason := "Only one sorting flag can be specified"
				return timelineFilters, backend.BadInputErr(errors.New(errReason), errReason)
			}
			direction := strings.ToLower(v[0])
			if direction == "asc" || direction == "chronological" || direction == "ascending" {
				timelineFilters.SortAsc = true
			}
		case "type":
			if len(v) != 1 {
				errReason := "Only one evidence type can be specified"
				return timelineFilters, backend.BadInputErr(errors.New(errReason), errReason)
			}
			timelineFilters.Type = v[0]
		default:
			errReason := fmt.Sprintf("Unknown filter key '%s'", k)
			return timelineFilters, backend.BadInputErr(errors.New(errReason), errReason)
		}
	}

	return timelineFilters, nil
}

// Parses the raw query string into a map.
//
// Examples:
// tokenizeTimelineQuery(`hello world tag:foo tag:bar is:event`)
// becomes
// map[string][]string{
//   "": []string{"hello", "world"},
//   "tag": []string{"foo", "bar", "fizz buzz"},
//   "is": []string{"event"},
// }
//
// Quotes act like they do in most shells (quotes prevent spaces from becoming splits):
// tokenizeTimelineQuery(`foo "bar baz" tag:"fizz buzz"`)
// becomes
// map[string][]string{
//   "": []string{"foo", "bar baz"},
//   "tag": []string{"fizz buzz"},
// }
func tokenizeTimelineQuery(query string) map[string][]string {
	parsed := map[string][]string{}
	currentToken := ""
	inQuote := false
	currentKey := ""

	for _, char := range query {
		switch char {
		case ' ':
			if !inQuote {
				if len(currentToken) > 0 {
					parsed[currentKey] = append(parsed[currentKey], currentToken)
				}
				currentToken = ""
				currentKey = ""
				continue
			}
		case '"':
			inQuote = !inQuote
			continue
		case ':':
			if currentKey == "" {
				currentKey = currentToken
				currentToken = ""
				continue
			}
		}
		currentToken += string(char)
	}
	if len(currentToken) > 0 {
		parsed[currentKey] = append(parsed[currentKey], currentToken)
	}
	return parsed
}

func parseRangeQuery(rangeQuery []string) (*DateRange, error) {
	if len(rangeQuery) != 1 {
		errReason := fmt.Sprintf("Query can only have one date ranges. (%d ranges supplied)", len(rangeQuery))
		return nil, backend.BadInputErr(errors.New(errReason), errReason)
	}
	split := strings.Split(rangeQuery[0], ",")
	if len(split) != 2 {
		errReason := fmt.Sprintf("Query range must be in the format [date],[date]. (Got '%s')", rangeQuery[0])
		return nil, backend.BadInputErr(errors.New(errReason), errReason)
	}
	from, err := parseTime(split[0], false)
	if err != nil {
		return nil, err
	}
	to, err := parseTime(split[1], true)
	if err != nil {
		return nil, err
	}
	return &DateRange{from, to}, nil
}

func parseTime(str string, useEndOfDayIfTimeIsMissing bool) (time.Time, error) {
	t, rfc3339Err := time.Parse(time.RFC3339, str)
	if rfc3339Err == nil {
		return t, nil
	}

	t, iso8601Err := time.Parse("2006-01-02", str)
	if iso8601Err == nil {
		if useEndOfDayIfTimeIsMissing {
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
		return t, nil
	}

	return time.Now(), backend.BadInputErr(
		fmt.Errorf("Failed to parse time. (RFC3339: %v) (ISO8601: %v)", rfc3339Err.Error(), iso8601Err.Error()),
		fmt.Sprintf("Query ranges must be in ISO8601 or RFC3339 format. (Got '%s')", str),
	)
}
