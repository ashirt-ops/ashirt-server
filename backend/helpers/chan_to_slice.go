// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers

// ChanToSlice is a generic function that converts a buffered channel into a slice. This
// consumes the channel in the process.
func ChanToSlice[T any](channel *chan T) []T {
	result := make([]T, len(*channel))
	for i := range result {
		result[i] = <-(*channel)
	}
	return result
}
