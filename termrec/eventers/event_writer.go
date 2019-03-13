package eventers

import (
	"time"

	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/recorders"
	"github.com/jonboulle/clockwork"
)

// EventWriter is structure that satisfies the io.Writer interface. Its objective is to interpret
// written bytes and form them into events that can be recorded to the provided Recorder.
// EventWriter supports the concept of middleware, which enables some transformations to occur
// to the raw event.
type EventWriter struct {
	rec        recorders.Recorder
	clock      clockwork.Clock
	eventType  common.EventType
	middleware []EventMiddleware
}

// RawEvent is a structure that can be sent to middleware to review/modify. The fields included
// are the basic components of an event, unmodified (specifically, Data, EventType, EventTime)
// TODO: should these be replaced with common.Event?
type RawEvent struct {
	Data      []byte
	EventTime time.Time
	EventType common.EventType
	Error     error
	rec       recorders.Recorder
}

// Dispatch provides a mechanism to send the event off to its Recorder. This is present so that a
// caching mechanism can be enabled, or alternatively, allowing some middleware to dispatch an
// event so that other middleware cannot inspect it (if such a feature is needed)
func (evt RawEvent) Dispatch() {
	evt.rec.AddEvent(evt.EventType, string(evt.Data), evt.EventTime)
}

// EventMiddleware is a type alias for easier middleware construction
type EventMiddleware func(RawEvent) RawEvent

func (e EventWriter) Write(p []byte) (n int, err error) {
	evt := RawEvent{Data: p, EventTime: e.clock.Now(), EventType: e.eventType, rec: e.rec}

	for _, mid := range e.middleware {
		evt = mid(evt)
	}
	if evt.Error != nil || len(evt.Data) == 0 {
		return 0, evt.Error // will still be nil if no error occurred, which is intended
	}

	evt.Dispatch()

	return len(evt.Data), nil
}

// NewEventWriter is a constructor for an EventWriter
func NewEventWriter(rec recorders.Recorder, eventType common.EventType, middleware ...EventMiddleware) EventWriter {
	return EventWriter{
		rec:        rec,
		clock:      clockwork.NewRealClock(),
		eventType:  eventType,
		middleware: middleware,
	}
}
