// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package dissectors

import (
	"fmt"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
)

// Coercable is a simple builder pattern to help transform the underlying value (if present)
// into the desired type. The logic aims to be relatively simple, while trying to handle the possible
// input variations. The underlying value types can be: string, []string, []interface{}, float64,
// bool
//
// Because of the variations for values, some conversions are done to make this feature extensible.
// However, no effort has been made to support converting from any type to any other type.
type Coercable struct {
	backReference *DissectedRequest
	fieldName     string
	rawValue      interface{}
	defaultValue  interface{}
	required      bool
	fieldPresent  bool
}

// makeCoercable constructs a basic Coercable
func makeCoercable(backref *DissectedRequest, key string, keyFound bool, value interface{}) *Coercable {
	return &Coercable{
		required:      false,
		rawValue:      value,
		backReference: backref,
		fieldName:     key,
		fieldPresent:  keyFound,
	}
}

// makeErrCoercable creates an empty Coercable in the case that DissectedRequest encounters an error
func makeErrCoercable(backref *DissectedRequest) *Coercable {
	return &Coercable{backReference: backref}
}

// OrDefault provides a default value for a given convesion, if the conversion would have
// otherwise failed. This is really only useful for optional values. Note: the provided value
// will itself be coerced into the target type. If the type cannot be coerced into the target type
// then the appropriate zero value will be returned.
//
// # This is a non-terminal Coercable action
//
// Examples:
//
//	assert(parsedJSON.FromBody("anInt").OrDefault(12).AsInt64(), int64(12))
//	assert(parsedJSON.FromBody("anInt").OrDefault("twelve").AsInt64(), 0)
func (c *Coercable) OrDefault(defaultValue interface{}) *Coercable {
	c.defaultValue = defaultValue
	return c
}

// Required marks the field as required, which, on some coercion issue, will set the
// DissectedRequest.Error field, indicating that an error has occurred
//
// By deafult, every field is Optional
func (c *Coercable) Required() *Coercable {
	c.required = true
	return c
}

// storeError attempts to record an encountered error (determined externally to this method).
// It will only do this if the field has been marked as required
func (c *Coercable) storeError(nameOfType, friendlyName string) {
	if c.required && c.backReference.Error == nil {
		if !c.fieldPresent {
			c.backReference.Error = backend.MissingValueErr(c.fieldName)
		} else {
			c.backReference.Error = backend.BadInputErr(
				fmt.Errorf("Unable to coerce into %v", nameOfType),
				fmt.Sprintf("%v must be a %v", c.fieldName, friendlyName),
			)
		}
	}
}

// AsString converts the Coercable into a string type.
// If this is impossible, then the zero value will be returned
//
// Note: if the underlying object is a []string, then the first value will be returned
func (c *Coercable) AsString() string {
	value, ok := maybeStringToString(c.rawValue)

	if !ok {
		c.storeError("string", "string")
		value, _ = c.defaultValue.(string)
	}
	return value
}

// AsStringPtr converts the Coercable into a *string type.
// If the underlying value is nil, then this will return nil.
// if the underlying value is not nil, but also not a string,
// then the zero value will be returned
func (c *Coercable) AsStringPtr() *string {
	if c.rawValue == nil {
		return nil
	}
	v := c.AsString()
	return &v
}

// AsTime converts the Coercable into a time.Time type.
// AsTime will look to parse an RFC3339 datetime string. If the string does not match the
// format, then it will not parse, and the zero value will be returned.
func (c *Coercable) AsTime() time.Time {
	value, ok := maybeTimeToTime(c.rawValue)

	if !ok {
		c.storeError("datetime", "datetime in RFC3339 format")
		value, _ = c.defaultValue.(time.Time)
	}
	return value
}

// AsUnixTime converts the Coercable into a time.Time type.
// AsUnixTime will look to convert the integer (nanoseconds since 1970) into the
// appropriate time format. If this fails then the zero value will be returned.
func (c *Coercable) AsUnixTime() time.Time {
	var value time.Time
	iValue, ok := maybeIntToInt64(c.rawValue, false)

	if ok {
		value = time.Unix(0, iValue)
	} else {
		c.storeError("datetime", "datetime in unix format")
		value, _ = c.defaultValue.(time.Time)
	}
	return value
}

// AsInt64 converts the Coercable into an int64 type.
// If this is impossible, then the zero value will be returned
//
// Note: if the underlying object is a []string, then the first value will be converted
// into an int64, then returned
func (c *Coercable) AsInt64() int64 {
	value, ok := maybeIntToInt64(c.rawValue, true)
	if !ok {
		c.storeError("int64", "int")
		value, _ = c.defaultValue.(int64)
	}
	return value
}

// AsInt64Slice converts the Coercable into an []int64 type.
// If this is impossible, then the zero value will be returned
//
// Note: if the underlying object is a []string, then each value will be converted
// into an int64, then returned
func (c *Coercable) AsInt64Slice() []int64 {
	value, ok := maybeIntSliceToInt64Slice(c.rawValue, true)

	if !ok {
		c.storeError("[]int64", "integer array")
		value, _ = c.defaultValue.([]int64)
	}
	return value
}

// AsStringSlice converts the Coercable into an []string type.
// If this is impossible, then the zero value will be returned
func (c *Coercable) AsStringSlice() []string {
	value, ok := maybeStringSliceToStringSlice(c.rawValue)

	if !ok {
		c.storeError("[]string", "string array")
		value, _ = c.defaultValue.([]string)
	}

	return value
}

// AsBool converts the Coercable into a bool type.
// If this is impossible, then the zero value will be returned
func (c *Coercable) AsBool() bool {
	value, ok := maybeBoolToBool(c.rawValue, c.fieldPresent)

	if !ok {
		c.storeError("bool", "boolean")
		value, _ = c.defaultValue.(bool)
	}

	return value
}

// AsBoolPtr converts the Coercable into a *bool type.
// If the underlying value is nil, then this will return nil.
// if the underlying value is not nil, but also not a bool,
// then the zero value will be returned
func (c *Coercable) AsBoolPtr() *bool {
	if c.rawValue == nil {
		return nil
	}

	v := c.AsBool()
	return &v
}
