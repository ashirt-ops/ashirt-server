// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package evidence

import "testing"

func TestSubstr(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{
			input:    "this is a test",
			length:   5,
			expected: "this ",
		},
		{
			input:    "this is a test",
			length:   100,
			expected: "this is a test",
		},
		{
			input:    "this is a test",
			length:   0,
			expected: "",
		},
		{
			input:    "💩 pile of poo",
			length:   1,
			expected: "",
		},
		{
			input:    "💩 pile of poo",
			length:   4,
			expected: "💩",
		},
		{
			input:    "💩 pile of poo",
			length:   6,
			expected: "💩 p",
		},
	}
	for _, test := range tests {
		output := substr(test.input, test.length)
		if output != test.expected {
			t.Logf("expected '%s' and got '%s'\n", test.expected, output)
			t.Fail()
		}
	}
}
