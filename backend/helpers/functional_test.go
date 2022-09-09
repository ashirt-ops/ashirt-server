// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theparanoids/ashirt-server/backend/helpers"
)

func TestMap(t *testing.T) {
	numbers := []int{1, 2, 3, 4}
	doubled := make([]int, len(numbers))
	for i, v := range numbers {
		doubled[i] = v * 2
	}
	doubleFn := func(i int) int {
		return i * 2
	}

	mappedNumbers := helpers.Map(numbers, doubleFn)
	require.Equal(t, doubled, mappedNumbers)

	numbersAsStr := make([]string, len(numbers))
	for i, v := range numbers {
		numbersAsStr[i] = fmt.Sprintf("%v", v)
	}
	asStrFn := func(i int) string {
		return fmt.Sprintf("%v", i)
	}
	mappedStrings := helpers.Map(numbers, asStrFn)
	require.Equal(t, numbersAsStr, mappedStrings)
}

func TestFind(t *testing.T) {
	// possible todo: check if this stops after finding the first match
	values := []int{1, 4, 9, 16, 25}
	target := 10
	var expectedIndex int = -1
	var expectedValue *int = nil

	for i := range values {
		if values[i] > target {
			expectedIndex = i
			expectedValue = &values[i]
			break
		}
	}

	index, foundValue := helpers.Find(values, func(i int) bool { return i > 10 })
	require.NotNil(t, foundValue)
	require.Equal(t, expectedIndex, index)
	require.Equal(t, expectedValue, foundValue)

	// check not-found branch
	index, foundValue = helpers.Find(values, func(i int) bool { return i < 0 })
	require.Nil(t, foundValue)
	require.Equal(t, -1, index)
}

func TestFindMatch(t *testing.T) {
	// possible todo: check if this stops after finding the first match
	values := []int{1, 4, 9, 16, 25}
	expectedIndex := 3
	expectedValue := &values[expectedIndex]

	index, foundValue := helpers.FindMatch(values, values[expectedIndex])
	require.NotNil(t, foundValue)
	require.Equal(t, expectedIndex, index)
	require.Equal(t, expectedValue, foundValue)

	// check not-found branch
	index, foundValue = helpers.FindMatch(values, -1)
	require.Nil(t, foundValue)
	require.Equal(t, -1, index)
}

func TestFilter(t *testing.T) {
	values := []int{1, 4, 9, 16, 25}
	expectedValues := []int{1, 9, 25}

	oddNumberFn := func(i int) bool {
		return i%2 == 1
	}
	greaterThan100Fn := func(i int) bool {
		return i > 100
	}

	actualResult := helpers.Filter(values, oddNumberFn)
	require.Equal(t, expectedValues, actualResult)

	actualResult = helpers.Filter(values, greaterThan100Fn)
	require.Equal(t, []int{}, actualResult)
}
