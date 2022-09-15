// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package helpers

// Map is a generic function that converts a list of type T into a list of type U, along with
// a function that converts a T to a U.
// This is essentially the same as a `map` function in other languages, like javascript
func Map[T any, U any](slice []T, mapFn func(T) U) []U {
	result := make([]U, len(slice))

	for i, t := range slice {
		result[i] = mapFn(t)
	}

	return result
}

// Find is a generic function that searches through a list searching for all items that match
// the given predicate. This returns the elements that matched the predicate in the order they were
// encountered.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	filteredData := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			filteredData = append(filteredData, v)
		}
	}
	return filteredData
}

// Find is a generic function that searches through a list searching for an item that matches
// the given predicate. If found, returns the index where it was found, and a pointer to the
// actual data. If not found, return (-1, nil)
// Note: the search is sequential, but terminated once the element is found
func Find[T any](slice []T, predicate func(T) bool) (int, *T) {
	for i, v := range slice {
		if predicate(v) {
			return i, &v
		}
	}
	return -1, nil
}

// FindMatch is a minor optimization of Find that restricts finds to only comparable elements.
func FindMatch[T comparable](slice []T, value T) (int, *T) {
	for i, v := range slice {
		if v == value {
			return i, &v
		}
	}
	return -1, nil
}

// Contains checks that a value exists within a slice. If not, returns false
func Contains[T any](slice []T, predicate func(T) bool) bool {
	index, _ := Find(slice, predicate)
	return index != -1
}

// ContainsMatch is an optimization on Contains that restricts checks to comparable elements.
func ContainsMatch[T comparable](slice []T, value T) bool {
	index, _ := FindMatch(slice, value)
	return index != -1
}
