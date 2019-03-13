package common

import (
	"fmt"
	"time"
)

// EventType is an enum for the kinds of events that take can take place (really just input/output)
type EventType string

const (
	// Input signals the EventType for input-related events
	Input EventType = "i"
	// Output signals the EventType for output-related events
	Output EventType = "o"
)

// Event is the structure of a generic terminal event. An event is comprised of 3 core components
// _When_ the event occurred, what kind of event it was (EventType), and what _Data_ was in the event
type Event struct {
	When time.Duration
	Type EventType
	Data string
}

func (evt Event) String() string {
	return fmt.Sprintf("%v %v %v", evt.When, evt.Type, evt.Data)
}
