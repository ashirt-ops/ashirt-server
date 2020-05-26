package recorders

import (
	"os"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/formatters"
	"github.com/theparanoids/ashirt/termrec/write"
)

// StreamingRecorder controls writes to a TerminalWriter. Events are added to the recorder as they
// are received. It is expected to be paried with a TerminalWriter that will keep the stream open
type StreamingRecorder struct {
	startTime time.Time
	clock     clockwork.Clock
	writer    write.TerminalWriter
	metadata  formatters.Metadata
}

// NewStreamingRecorder makes a new StreamingRecorder. The provided terminal writer should allow for
// an open session (see StreamingFileWriter)
//
// Parameter Notes:
// clock: use clockwork.NewRealClock() unless you are testing
// shell: Passed along in metadata
//
// Note: start time is set to now. Unfortunately, this cannot be made lazy, due to a shell prompt
// coming up immediately
func NewStreamingRecorder(writer write.TerminalWriter, clock clockwork.Clock, shell string) StreamingRecorder {
	rtn := StreamingRecorder{
		startTime: clock.Now(),
		clock:     clock,
		writer:    writer,
		metadata: formatters.Metadata{
			StartTimeUnix: clock.Now().Unix(),
			Term:          os.Getenv("TERM"),
			Shell:         shell,
		},
	}

	rtn.writer.WriteHeader(rtn.metadata)

	return rtn
}

// AddEvent adds an event to the stream with an arbitrary timestamp
func (r *StreamingRecorder) AddEvent(eType common.EventType, data string, evtTime time.Time) {
	r.writer.WriteEvent(common.Event{
		When: evtTime.Sub(r.startTime),
		Type: eType,
		Data: data,
	})
}

// GetEventCount returns -1, as we don't know how many events we have, or will, process.
func (r *StreamingRecorder) GetEventCount() int {
	return -1
}

// GetDurationInSeconds returns the elapsed time from the start of the stream
func (r *StreamingRecorder) GetDurationInSeconds() float64 {
	return r.clock.Now().Sub(r.startTime).Seconds()
}

// GetStartTime returns the unix time that represents the start of this stream
func (r *StreamingRecorder) GetStartTime() int64 {
	return r.startTime.Unix()
}

// Output writes the footer to the StreamingRecorder's TerminalWriter (and ignores the passed parameter)
func (r *StreamingRecorder) Output(_ write.TerminalWriter) {
	r.writer.WriteFooter(r.metadata)
}
