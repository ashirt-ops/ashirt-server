// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers

import "strings"

func StrToUpperCaseUnderscore(str string) string {
	return strings.Replace(strings.ToUpper(str), " ", "_", -1)
}
