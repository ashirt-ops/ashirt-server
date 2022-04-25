package helpers

import "fmt"

// StringPtr converts a string into a *string
func StringPtr(s string) *string {
	return &s
}

// SprintfPtr is a wrapper around Sprintf that returns the result as a string pointer
func SprintfPtr(s string, vals ...any) *string {
	whatever := fmt.Sprintf(s, vals...)
	return StringPtr(whatever)
}

// I64Ptr converts a int64 into an *int64
func I64Ptr(i int64) *int64 {
	return &i
}

// BoolPtr converts a bool into an *bool
func BoolPtr(b bool) *bool {
	return &b
}
