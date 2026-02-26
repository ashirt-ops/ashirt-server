package filter

type DateValue struct {
	Value    DateRange
	Modifier FilterModifier
}

type DateValues []DateValue

// DateVal is a shorthand method for creating a standard, un-modified value.
func DateVal(val DateRange) DateValue {
	return DateValue{Value: val}
}

// NotDateVal is a shorthand method for creating a filter value with the Not modification.
func NotDateVal(val DateRange) DateValue {
	return DateValue{Value: val, Modifier: Not}
}

// Values converts a filter.Values into a string slice for easier consumption
// essentially [].map(fv => fv.Value) in javascript
func (f DateValues) Values() []DateRange {
	values := make([]DateRange, len(f))
	for i, v := range f {
		values[i] = v.Value
	}
	return values
}

type DateMap = map[string][]DateRange

// SplitValues divides a Values into a map, based on the given partitioning function
func (f DateValues) SplitValues(partitionFn func(DateValue) string) DateMap {
	splitValues := make(DateMap)

	for _, v := range f {
		key := partitionFn(v)
		splitValues[key] = append(splitValues[key], v.Value)
	}
	return splitValues
}

// Value is a shorthand method for retriving the string at a particular index
func (f DateValues) Value(index int) DateRange {
	return f[index].Value
}

// SplitByModifier divides up a Values into pieces based on what its modifier
// is.
func (f DateValues) SplitByModifier() map[FilterModifier][]DateRange {
	splitValues := make(map[FilterModifier][]DateRange)

	for _, v := range f {
		splitValues[v.Modifier] = append(splitValues[v.Modifier], v.Value)
	}
	return splitValues
}
