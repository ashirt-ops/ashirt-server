// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/helpers/filter"
)

func TestTokenizeTimelineQuery(t *testing.T) {
	// normal tests
	plainTextPart := []string{"some", "text", "plain"}
	tokenValue := "token"
	normalToken := "normal:" + tokenValue
	notToken := "not:!" + tokenValue

	query := strings.Join([]string{strings.Join(plainTextPart, " "), normalToken, notToken}, " ")
	result := tokenizeTimelineQuery(query)

	for i := range plainTextPart {
		require.Equal(t, filter.Value{Value: plainTextPart[i], Modifier: filter.Normal}, result[""][i])
	}

	require.Equal(t, filter.Value{Value: tokenValue, Modifier: filter.Normal}, result["normal"][0])
	require.Equal(t, filter.Value{Value: tokenValue, Modifier: filter.Not}, result["not"][0])

	// complex tests

	tokenValueTwo := "double-token"
	tokenValueThree := `test three`
	morePlainText := []string{"some", "text", "plain", "!plain"}
	normalSecondToken := "normal:" + tokenValueTwo
	notSecondToken := "not:!" + tokenValueTwo
	normalThirdToken := `normal:!"` + tokenValueThree + `"`

	query = strings.Join([]string{
		strings.Join(morePlainText, " "),
		normalToken, notToken,
		normalSecondToken, notSecondToken,
		normalThirdToken,
	}, " ")
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

	query = strings.Join([]string{
		normalToken, notToken,
		normalSecondToken, notSecondToken,
		strings.Join(morePlainText[0:2], " "),
		normalThirdToken,
		strings.Join(morePlainText[2:], " "),
	}, " ")
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
