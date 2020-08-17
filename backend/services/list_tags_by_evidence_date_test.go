// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestListTagsByEvidenceDate(t *testing.T) {
	db := initTest(t)
	HarryPotterSeedData.ApplyTo(t, db)

	masterOp := OpGanttChart
	input := services.ListTagsByEvidenceDateInput{
		OperationSlug: masterOp.Slug,
	}

	expectedData := HarryPotterSeedData.TagIDsUsageByDate(masterOp.ID)

	// test no-access
	_, err := services.ListTagsByEvidenceDate(contextForUser(UserDraco, db), db, input)
	require.Error(t, err)

	// test read
	actualData, err := services.ListTagsByEvidenceDate(contextForUser(UserGinny, db), db, input)
	require.NoError(t, err)
	validateTagUsageData(t, HarryPotterSeedData, expectedData, actualData)

	// test write
	actualData, err = services.ListTagsByEvidenceDate(contextForUser(UserHarry, db), db, input)
	require.NoError(t, err)
	validateTagUsageData(t, HarryPotterSeedData, expectedData, actualData)

	// test admin
	actualData, err = services.ListTagsByEvidenceDate(contextForUser(UserDumbledore, db), db, input)
	require.NoError(t, err)
	validateTagUsageData(t, HarryPotterSeedData, expectedData, actualData)
}

func validateTagUsageData(t *testing.T, seed TestSeedData, expectedTagIDUsage map[int64][]time.Time, actual []*dtos.TagByEvidenceDate) {
	require.Equal(t, len(expectedTagIDUsage), len(actual))
	for tagID, dates := range expectedTagIDUsage {
		tag := seed.GetTagFromID(tagID)

		var match *dtos.TagByEvidenceDate = nil
		for _, tagUsage := range actual {
			if tagUsage.ID == tag.ID {
				match = tagUsage
				break
			}
		}
		require.NotNil(t, match)
		require.Equal(t, tag.ColorName, match.ColorName)
		require.Equal(t, tag.Name, match.Name)

		mapDates := func(times []time.Time) map[int64]int {
			rtn := make(map[int64]int)
			for _, someTime := range times {
				timeCount := 0
				for _, matchTime := range times {
					if someTime == matchTime {
						timeCount += 1
					}
				}
				rtn[someTime.Unix()] = timeCount
			}
			return rtn
		}

		expectedDateUsage := mapDates(dates)
		actualDateUsage := mapDates(match.UsageDates)
		require.Equal(t, expectedDateUsage, actualDateUsage)
	}
}
