// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers

import "fmt"

// SprintfPtr is a wrapper around Sprintf that returns the result as a string pointer
func SprintfPtr(s string, vals ...any) *string {
	whatever := fmt.Sprintf(s, vals...)
	return Ptr(whatever)
}

// Ptr is a small helper to convert a real value into a pointer to that value.
// Most useful as a way to turn a literal into a pointer to that literal
func Ptr[T any](t T) *T {
	return &t
}

// PTrue returns a pointer to a true value
func PTrue() *bool {
	return Ptr(true)
}

// PFalse returns a pointer to a false value
func PFalse() *bool {
	return Ptr(false)
}
