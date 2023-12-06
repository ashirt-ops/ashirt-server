package filter

import "time"

// DateRange is a simple struct representing a slice of time From a point To a point
type DateRange struct {
	From time.Time
	To   time.Time
}
