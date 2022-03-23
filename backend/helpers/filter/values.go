// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package filter

type StringModifier int64

const (
	Normal StringModifier = 0
	Not    StringModifier = 1 << (iota)
)

type Value struct {
	Value    string
	Modifier StringModifier
}

type Values []Value

// Val is a shorthand method for creating a standard, un-modified value.
func Val(n string) Value {
	return Value{Value: n}
}

// NotVal is a shorthand method for creating a filter value with the Not modification.
func NotVal(n string) Value {
	return Value{Value: n, Modifier: Not}
}

// Values converts a filter.Values into a string slice for easier consumption
// essentially [].map(fv => fv.Value) in javascript
func (f Values) Values() []string {
	values := make([]string, len(f))
	for i, v := range f {
		values[i] = v.Value
	}
	return values
}

// SplitValues divides a Values into a map, based on the given partitioning function
func (f Values) SplitValues(partitionFn func(Value) string) map[string][]string {
	splitValues := make(map[string][]string)

	for _, v := range f {
		key := partitionFn(v)
		splitValues[key] = append(splitValues[key], v.Value)
	}
	return splitValues
}

// Value is a shorthand method for retriving the string at a particular index
func (f Values) Value(index int) string {
	return f[index].Value
}

// SplitByModifier divides up a Values into pieces based on what its modifier
// is. 
func (f Values) SplitByModifier() map[StringModifier][]string {
	splitValues := make(map[StringModifier][]string)

	for _, v := range f {
		splitValues[v.Modifier] = append(splitValues[v.Modifier], v.Value)
	}
	return splitValues
}
