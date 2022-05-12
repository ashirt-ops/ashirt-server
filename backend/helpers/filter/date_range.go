// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package filter

import "time"

// DateRange is a simple struct representing a slice of time From a point To a point
type DateRange struct {
	From time.Time
	To   time.Time
}
