package helpers

import "strings"

func StrToUpperCaseUnderscore(str string) string {
	return strings.Replace(strings.ToUpper(str), " ", "_", -1)
}

func StrToLowerCaseUnderscore(str string) string {
	return strings.Replace(strings.ToLower(str), " ", "_", -1)
}
