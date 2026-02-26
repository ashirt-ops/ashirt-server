package dissectors

import (
	"strconv"
	"time"
)

// maybeBoolToBool tries to convert string, bool or []string into
// (bool, true). If not possible, then (false, false) is returned
//
// Note: for strings, in order to accomodate the &flag query string parameter,
// we will also assume the value is true if the key is present, even if the value is ""
func maybeBoolToBool(val interface{}, keyPresent bool) (bool, bool) {
	switch coercedValue := val.(type) {
	case bool:
		return coercedValue, true
	case string:
		bValue, err := strconv.ParseBool(coercedValue)
		return bValue, (err == nil)
	case []string:
		if len(coercedValue) > 0 {
			if coercedValue[0] == "" && keyPresent {
				return true, true
			}
			bValue, err := strconv.ParseBool(coercedValue[0])
			return bValue, (err == nil)
		}
		return false, false

	default:
		return false, false
	}
}

// maybeStringSliceToStringSlice tries to convert []string or []interface{} into
// ([]string, true). If not possible, then ([]string{}, false) is returned
func maybeStringSliceToStringSlice(val interface{}) ([]string, bool) {
	switch coercedValue := val.(type) {
	case []string: // query param
		return coercedValue, true

	case []interface{}:
		allStrings := make([]string, len(coercedValue))
		for idx, maybeString := range coercedValue {
			someString, convertOk := maybeString.(string)
			if !convertOk {
				return []string{}, false
			}
			allStrings[idx] = someString
		}
		return allStrings, true

	default:
		return []string{}, false
	}
}

// maybeStringSliceToStringSlice tries to convert []string, []interface{} or []float64 into
// ([]int64, true). If not possible, then ([]int64{}, false) is returned.
// The strict parameter specifies whether the converted floats must have the same integer
// value as the original floats. If not, then ([]int64{}, false) is returned. Otherwise,
// the converted values are returned
func maybeIntSliceToInt64Slice(val interface{}, strict bool) ([]int64, bool) {
	mapFloatSliceToIntSlice := func(in []float64) ([]int64, bool) {
		allInts := make([]int64, len(in))
		for idx, someFloat := range in {
			intValue := int64(someFloat)
			if strict && someFloat != float64(intValue) {
				return []int64{}, false
			}
			allInts[idx] = intValue
		}
		return allInts, true
	}

	mapMaybeFloatSliceToIntSlice := func(in []interface{}) ([]int64, bool) {
		allFloats := make([]float64, len(in))
		for idx, maybeFloat := range in {
			someFloat, convertOk := maybeFloat.(float64)
			if !convertOk {
				return []int64{}, false
			}
			allFloats[idx] = someFloat
		}
		return mapFloatSliceToIntSlice(allFloats)
	}

	mapStringSliceToIntSlice := func(in []string) ([]int64, bool) {
		allInts := make([]int64, len(in))
		for idx, someStr := range in {
			intValue, err := strconv.ParseInt(someStr, 10, 64)
			if err != nil {
				return []int64{}, false
			}
			allInts[idx] = intValue
		}
		return allInts, true
	}

	switch coercedValue := val.(type) {
	case []float64: //possibly dead code -- what will send a []float64?
		return mapFloatSliceToIntSlice(coercedValue)
	case []string: // query param
		return mapStringSliceToIntSlice(coercedValue)
	case []interface{}:
		return mapMaybeFloatSliceToIntSlice(coercedValue)

	default:
		return []int64{}, false
	}
}

// maybeIntToInt64 tries to convert []string, string or float64 into
// (int64, true). If not possible, then (0, false) is returned.
// The strict parameter specifies whether the converted float must have the same integer
// value as the original float. If not, then (0, false) is returned. Otherwise,
// the converted value is returned
func maybeIntToInt64(val interface{}, strict bool) (int64, bool) {
	switch coercedValue := val.(type) {
	case float64: // number in json
		iValue := int64(coercedValue)
		if strict && coercedValue != float64(iValue) {
			return 0, false
		}
		return iValue, true

	case string: // url param / string in json
		iValue, err := strconv.ParseInt(coercedValue, 10, 64)
		return iValue, (err == nil)

	case []string: // query param
		if len(coercedValue) > 0 {
			iValue, err := strconv.ParseInt(coercedValue[0], 10, 64)
			return iValue, (err == nil)
		}
		return 0, false

	default:
		return 0, false
	}
}

// maybeStringToString tries to convert []string or string into
// (string, true). If not possible, then ("", false) is returned.
func maybeStringToString(val interface{}) (string, bool) {
	switch coercedValue := val.(type) {
	case string: // url param / string in json
		return coercedValue, true
	case []string: // query param
		if len(coercedValue) > 0 {
			return coercedValue[0], true
		}
		return "", false
	default:
		return "", false
	}
}

// maybeTimeToTime tries to convert []string or string into
// (time.Time, true). If not possible, then (time.Time{}, false) is returned.
func maybeTimeToTime(val interface{}) (time.Time, bool) {
	switch coercedValue := val.(type) {
	case string: // url param / string in json
		tValue, err := time.Parse(time.RFC3339, coercedValue)
		return tValue, (err == nil)
	case []string: // query param
		if len(coercedValue) > 0 {
			tValue, err := time.Parse(time.RFC3339, coercedValue[0])
			return tValue, (err == nil)
		}
		return time.Time{}, false
	default:
		return time.Time{}, false
	}
}
