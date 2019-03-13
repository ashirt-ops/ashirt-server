package recorders

import (
	"time"

	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/write"
)

// Recorder is an interface for tracking I/O events
type Recorder interface {
	AddEvent(common.EventType, string, time.Time)
	GetEventCount() int
	GetDurationInSeconds() float64
	GetStartTime() int64
	Output(write.TerminalWriter)
}
