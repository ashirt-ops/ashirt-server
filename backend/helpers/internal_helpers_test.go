// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/helpers/filter"
)

type StringSlice []string

func (s StringSlice) Join(between string) string {
	if len(s) == 0 {
		return ""
	}

	rtn := s[0]
	for _, v := range s[1:] {
		rtn += between + v
	}

	return rtn
}
func (s StringSlice) AsSlice() []string {
	return s
}

func TestTokenizeTimelineQuery(t *testing.T) {
	// normal tests
	plainTextPart := StringSlice{"some", "text", "plain"}
	tokenValue := "token"
	normalToken := "normal:" + tokenValue
	notToken := "not:!" + tokenValue

	query := StringSlice{plainTextPart.Join(" "), normalToken, notToken}.Join(" ")
	result := tokenizeTimelineQuery(query)

	for i := range plainTextPart {
		require.Equal(t, filter.Value{Value: plainTextPart[i], Modifier: filter.Normal}, result[""][i])
	}

	require.Equal(t, filter.Value{Value: tokenValue, Modifier: filter.Normal}, result["normal"][0])
	require.Equal(t, filter.Value{Value: tokenValue, Modifier: filter.Not}, result["not"][0])

	// complex tests

	tokenValueTwo := "double-token"
	tokenValueThree := `test three`
	morePlainText := StringSlice{"some", "text", "plain", "!plain"}
	normalSecondToken := "normal:" + tokenValueTwo
	notSecondToken := "not:!" + tokenValueTwo
	normalThirdToken := `normal:!"` + tokenValueThree + `"`

	query = StringSlice{
		morePlainText.Join(" "),
		normalToken, notToken,
		normalSecondToken, notSecondToken,
		normalThirdToken,
	}.Join(" ")
	result = tokenizeTimelineQuery(query)

	for i := range morePlainText {
		require.Equal(t, filter.Value{Value: morePlainText[i], Modifier: filter.Normal}, result[""][i])
	}

	require.Equal(t, filter.Value{Value: tokenValue, Modifier: filter.Normal}, result["normal"][0])
	require.Equal(t, filter.Value{Value: tokenValueTwo, Modifier: filter.Normal}, result["normal"][1])
	require.Equal(t, filter.Value{Value: tokenValueThree, Modifier: filter.Not}, result["normal"][2])
	require.Equal(t, filter.Value{Value: tokenValue, Modifier: filter.Not}, result["not"][0])
	require.Equal(t, filter.Value{Value: tokenValueTwo, Modifier: filter.Not}, result["not"][1])

	// mixed up

	query = StringSlice{
		normalToken, notToken,
		normalSecondToken, notSecondToken,
		morePlainText[0:2].Join(" "),
		normalThirdToken,
		morePlainText[2:].Join(" "),
	}.Join(" ")
	result = tokenizeTimelineQuery(query)

	for i := range morePlainText {
		require.Equal(t, filter.Value{Value: morePlainText[i], Modifier: filter.Normal}, result[""][i])
	}

	require.Equal(t, filter.Value{Value: tokenValue, Modifier: filter.Normal}, result["normal"][0])
	require.Equal(t, filter.Value{Value: tokenValueTwo, Modifier: filter.Normal}, result["normal"][1])
	require.Equal(t, filter.Value{Value: tokenValueThree, Modifier: filter.Not}, result["normal"][2])
	require.Equal(t, filter.Value{Value: tokenValue, Modifier: filter.Not}, result["not"][0])
	require.Equal(t, filter.Value{Value: tokenValueTwo, Modifier: filter.Not}, result["not"][1])
}
