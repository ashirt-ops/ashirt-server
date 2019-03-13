package recorders

import (
	"os"
	"testing"
	"time"

	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/formatters"
	"github.com/theparanoids/ashirt/termrec/write"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

func makeStreamingRecorder() (StreamingRecorder, write.SaveTermWriter, clockwork.FakeClock) {
	clock := clockwork.NewFakeClock()
	writer := write.NewSaveTermWrier()
	return NewStreamingRecorder(writer, clock, "someShell"), writer, clock
}

func TestStreamingRecorderConstructor(t *testing.T) {
	clock := clockwork.NewFakeClock()
	writer := write.NewSaveTermWrier()
	rec := NewStreamingRecorder(writer, clock, "someShell")

	expectedMetadata := formatters.Metadata{
		StartTimeUnix: clock.Now().Unix(),
		Term:          os.Getenv("TERM"),
		Shell:         "someShell",
	}

	assert.Equal(t, rec.metadata, expectedMetadata)
	assert.Equal(t, *writer.HeaderMetadata, expectedMetadata)
}

func TestStreamingRecorderAddEvent(t *testing.T) {
	rec, write, clock := makeStreamingRecorder()

	evt := common.Event{
		Type: "i",
		Data: "someData",
		When: 2 * time.Second,
	}
	clock.Advance(2 * time.Second)
	rec.AddEvent(evt.Type, evt.Data, clock.Now())

	assert.Equal(t, len(*write.AllEvents), 1)
	assert.Equal(t, (*write.AllEvents)[0], evt)
}

func TestStreamingRecorderGetEventCount(t *testing.T) {
	rec, _, clock := makeStreamingRecorder()

	rec.AddEvent("i", "Data", clock.Now())

	assert.Equal(t, rec.GetEventCount(), -1, "Verify that streaming recorder drops count")
}

func TestStreamingRecorderGetDurationInSeconds(t *testing.T) {
	rec, _, clock := makeStreamingRecorder()

	clock.Advance(2 * time.Second)
	rec.AddEvent("i", "someData", clock.Now())

	clock.Advance(2 * time.Second)
	rec.AddEvent("i", "someMoreData", clock.Now())

	assert.Equal(t, rec.GetDurationInSeconds(), float64(4))
}

func TestStreamingRecorderGetStartTime(t *testing.T) {
	rec, _, clock := makeStreamingRecorder()
	now := clock.Now().Unix()

	clock.Advance(2 * time.Second)
	rec.AddEvent("i", "someData", clock.Now())

	assert.Equal(t, now, rec.GetStartTime())
}

func TestStreamingRecorderOutput(t *testing.T) {
	clock := clockwork.NewFakeClock()
	writer := write.NewSaveTermWrier()
	rec := NewStreamingRecorder(writer, clock, "someShell")

	expectedMetadata := formatters.Metadata{
		StartTimeUnix: clock.Now().Unix(),
		Term:          os.Getenv("TERM"),
		Shell:         "someShell",
	}

	rec.Output(write.NilTermWriter{})
	assert.Equal(t, *writer.FooterMetadata, expectedMetadata)
}
