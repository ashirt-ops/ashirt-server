// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers

// Clamp is a small helper that ensures that a given number v lies between bounds m, n. If it is outside
// of those bounds, it sets v to the closer of (m, n)
func Clamp(v, min, max int64) int64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
