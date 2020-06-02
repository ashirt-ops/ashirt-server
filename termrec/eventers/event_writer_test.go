package eventers

import (
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/recorders"
)

func newTestEventWriter(middleware ...EventMiddleware) (EventWriter, *recorders.BufferedRecorder) {
	rec := recorders.NewBufferedRecorder(clockwork.NewFakeClock(), "whatever")
	return EventWriter{
		rec:        &rec,
		clock:      clockwork.NewRealClock(),
		eventType:  common.Input,
		middleware: middleware,
	}, &rec
}

func TestEventWriterConstructor(t *testing.T) {
	now := time.Now()
	rec := recorders.NewBufferedRecorder(clockwork.NewRealClock(), "whatever")
	ew := NewEventWriter(&rec, common.Input)
	diff := ew.clock.Now().Sub(now).Seconds()
	assert.True(t, diff < 1, "EventWriter constructor is using real time")
}

func TestEventWriterPlainWrite(t *testing.T) {
	ew, rec := newTestEventWriter()
	msg := []byte("New Event!")
	n, err := ew.Write(msg)
	assert.Nil(t, err)
	assert.Equal(t, n, len(msg))
	assert.Equal(t, len(rec.GetEventsForTesting()), 1)
	assert.Equal(t, rec.GetEventsForTesting()[0].Data, string(msg))
}

func TestEventWriterMiddlewareWrite(t *testing.T) {
	ew, rec := newTestEventWriter(withEmphasisMiddleware)
	msg := []byte("New Event!")
	n, err := ew.Write(msg)
	assert.Nil(t, err)
	assert.Equal(t, n, len(msg)+1)
	assert.Equal(t, len(rec.GetEventsForTesting()), 1)
	assert.Equal(t, rec.GetEventsForTesting()[0].Data, string(append(msg, '!')))
}

func TestEventWriterMiddlewareError(t *testing.T) {
	ew, _ := newTestEventWriter(withErrorMiddleware)
	msg := []byte("New Event!")
	_, err := ew.Write(msg)
	assert.NotNil(t, err)
}

func withEmphasisMiddleware(evt RawEvent) RawEvent {
	evt.Data = append(evt.Data, '!')
	return evt
}

func withErrorMiddleware(evt RawEvent) RawEvent {
	evt.Error = fmt.Errorf("boo")
	return evt
}
