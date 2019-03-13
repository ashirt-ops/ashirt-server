package recorders

import (
	"os"
	"time"

	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/formatters"
	"github.com/theparanoids/ashirt/termrec/write"

	"github.com/jonboulle/clockwork"
)

// BufferedRecorder is an alternative Recorder that can be used when users are confident that a
// recording will be short. It provides some extra metadata for recording files. Not currently
// recommended for use.
type BufferedRecorder struct {
	startTime time.Time
	events    []common.Event
	metadata  formatters.Metadata
}

// NewBufferedRecorder is a constructor for a BufferedRecorder
func NewBufferedRecorder(clock clockwork.Clock, shell string) BufferedRecorder {
	now := clock.Now()
	return BufferedRecorder{
		startTime: now,
		events:    make([]common.Event, 0),
		metadata: formatters.Metadata{
			StartTimeUnix: now.Unix(),
			Term:          os.Getenv("TERM"),
			Shell:         shell,
		},
	}
}

// AddEvent records an event from the provided input
func (r *BufferedRecorder) AddEvent(eType common.EventType, data string, evtTime time.Time) {
	r.events = append(r.events, common.Event{
		When: evtTime.Sub(r.startTime),
		Type: eType,
		Data: data,
	})
}

// GetEventCount returns the number of encoutnered events in this stream
func (r *BufferedRecorder) GetEventCount() int {
	return len(r.events)
}

// GetDurationInSeconds returns the difference between the very first action and the very last action.
// If not events have been recorded yet, returns 0
func (r *BufferedRecorder) GetDurationInSeconds() float64 {
	if len(r.events) == 0 {
		return 0
	}
	return (r.events[len(r.events)-1]).When.Seconds()
}

// GetStartTime returns the start of the recording, in unix time
func (r *BufferedRecorder) GetStartTime() int64 {
	return r.startTime.Unix()
}

// Output writes the recorded events to the provided TerminalWriter
func (r *BufferedRecorder) Output(dst write.TerminalWriter) {
	r.metadata.DurationSeconds = r.GetDurationInSeconds()
	dst.WriteHeader(r.metadata)
	for _, evt := range r.events {
		dst.WriteEvent(evt)
	}
	dst.WriteFooter(r.metadata)
}

// GetEventsForTesting is a method that simply returns the unxpected events field. To be used only
// for testing purposes
func (r *BufferedRecorder) GetEventsForTesting() []common.Event {
	return r.events
}
