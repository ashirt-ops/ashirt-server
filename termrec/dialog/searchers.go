package dialog

import "strings"

// SearcherContainsCI provides a basic searcher with the logic that a given label will match
// if, **ignoring case**, the given input matches some portion of the given label;
// specifically, if the input is a substring of the label
func SearcherContainsCI(optionLabels []string) func(input string, index int) bool {
	return func(input string, index int) bool {
		return strings.Contains(strings.ToLower(optionLabels[index]), strings.ToLower(input))
	}
}
