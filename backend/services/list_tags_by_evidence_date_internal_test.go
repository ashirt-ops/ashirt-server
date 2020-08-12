// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func timeToDBTimestamp(d1 time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", d1.Year(), d1.Month(), d1.Day(), d1.Hour(), d1.Minute(), d1.Second())
}

func TestSliceStrDatesToSliceDates(t *testing.T) {
	d1 := time.Date(2014, 9, 14, 9, 45, 3, 0, time.UTC)
	d2 := time.Date(2016, 11, 14, 8, 23, 3, 0, time.UTC)

	// test no dates
	result, err := sliceStrDatesToSliceDates([]string{})
	require.NoError(t, err)
	require.Equal(t, []time.Time{}, result)

	// test one date
	result, err = sliceStrDatesToSliceDates([]string{timeToDBTimestamp(d1)})
	require.NoError(t, err)
	require.Equal(t, []time.Time{d1}, result)

	// test two dates
	result, err = sliceStrDatesToSliceDates([]string{timeToDBTimestamp(d1), timeToDBTimestamp(d2)})
	require.NoError(t, err)
	require.Equal(t, []time.Time{d1, d2}, result)

	// test bad dates
	result, err = sliceStrDatesToSliceDates([]string{timeToDBTimestamp(d1), "Aug 12, 2012"})
	require.Error(t, err)

}
