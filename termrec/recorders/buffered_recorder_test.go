package recorders

import (
	"os"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/formatters"
	"github.com/theparanoids/ashirt/termrec/write"
)

func TestBufferedRecorderConstructor(t *testing.T) {
	rec := NewBufferedRecorder(clockwork.NewFakeClock(), "someShell")

	assert.Equal(t, rec.metadata.Shell, "someShell")
	assert.Equal(t, len(rec.events), 0)
}

func TestBufferedRecorderAddEvent(t *testing.T) {
	clock := clockwork.NewFakeClock()
	rec := NewBufferedRecorder(clock, "someShell")

	evt := common.Event{
		Type: "i",
		Data: "someData",
		When: 2 * time.Second,
	}
	clock.Advance(2 * time.Second)
	rec.AddEvent(evt.Type, evt.Data, clock.Now())

	assert.Equal(t, len(rec.events), 1)
	assert.Equal(t, rec.events[0], evt)
}

func TestBufferedRecorderGetEventCount(t *testing.T) {
	clock := clockwork.NewFakeClock()
	rec := NewBufferedRecorder(clock, "someShell")

	rec.AddEvent("i", "Data", clock.Now())

	assert.Equal(t, rec.GetEventCount(), 1)
}

func TestBufferedRecorderGetDurationInSeconds(t *testing.T) {
	clock := clockwork.NewFakeClock()
	rec := NewBufferedRecorder(clock, "someShell")

	clock.Advance(2 * time.Second)
	rec.AddEvent("i", "someData", clock.Now())

	clock.Advance(2 * time.Second)
	rec.AddEvent("i", "someMoreData", clock.Now())

	assert.Equal(t, len(rec.events), 2)
	assert.Equal(t, rec.GetDurationInSeconds(), float64(4))
}

func TestBufferedRecorderGetDurationInSecondsForNoEvents(t *testing.T) {
	clock := clockwork.NewFakeClock()
	rec := NewBufferedRecorder(clock, "someShell")

	assert.Equal(t, rec.GetDurationInSeconds(), float64(0))
}

func TestBufferedRecorderGetStartTime(t *testing.T) {
	clock := clockwork.NewFakeClock()
	rec := NewBufferedRecorder(clock, "someShell")
	now := clock.Now().Unix()

	clock.Advance(2 * time.Second)
	rec.AddEvent("i", "someData", clock.Now())

	assert.Equal(t, now, rec.GetStartTime())
}

func TestBufferedRecorderOutput(t *testing.T) {
	writer := write.NewSaveTermWrier()
	clock := clockwork.NewFakeClock()
	rec := NewBufferedRecorder(clock, "someShell")

	evt1 := common.Event{Type: "i", Data: "someData", When: 2 * time.Second}
	evt2 := common.Event{Type: "o", Data: "moreData", When: 4 * time.Second}

	expectedMetadata := formatters.Metadata{
		StartTimeUnix:   clock.Now().Unix(),
		Term:            os.Getenv("TERM"),
		Shell:           "someShell",
		DurationSeconds: 4,
	}

	clock.Advance(2 * time.Second)
	rec.AddEvent(evt1.Type, evt1.Data, clock.Now())

	clock.Advance(2 * time.Second)
	rec.AddEvent(evt2.Type, evt2.Data, clock.Now())

	expectedEvents := []common.Event{evt1, evt2}

	rec.Output(writer)
	assert.Equal(t, *writer.HeaderMetadata, expectedMetadata)
	assert.Equal(t, *writer.AllEvents, expectedEvents)
	assert.Equal(t, *writer.FooterMetadata, expectedMetadata)
}

func TestBufferedRecorderGetEventsForTesting(t *testing.T) {
	clock := clockwork.NewFakeClock()
	rec := NewBufferedRecorder(clock, "someShell")

	evt1 := common.Event{Type: "i", Data: "someData", When: 2 * time.Second}
	evt2 := common.Event{Type: "o", Data: "moreData", When: 4 * time.Second}

	clock.Advance(2 * time.Second)
	rec.AddEvent(evt1.Type, evt1.Data, clock.Now())

	clock.Advance(2 * time.Second)
	rec.AddEvent(evt2.Type, evt2.Data, clock.Now())

	expectedEvents := []common.Event{evt1, evt2}

	assert.Equal(t, rec.events, expectedEvents)
	assert.Equal(t, expectedEvents, rec.GetEventsForTesting())
}
