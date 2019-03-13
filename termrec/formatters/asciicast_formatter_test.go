package formatters

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/systemstate"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
)

func TestAsciiCastFormatterWriteFooter(t *testing.T) {
	formatter := asciiCast{clock: clockwork.NewFakeClock()}
	bytes, err := formatter.WriteFooter(Metadata{})

	assert.Equal(t, []byte{}, bytes)
	assert.Nil(t, err)
}

func TestAsciiCastDefaultUsesRealTime(t *testing.T) {
	now := time.Now()
	asciiNow := ASCIICast.clock.Now()
	// not sure if we can really do a better test -- this ensures that now is what asciiCast uses as now
	// are nearly the same, which is probably good enough.
	assert.True(t, asciiNow.Sub(now) < 1*time.Second)
}

func TestAsciiCastFormatterWriteEvent(t *testing.T) {
	formatter := asciiCast{clock: clockwork.NewFakeClock()}
	evt := common.Event{When: 1 * time.Second, Type: common.Input, Data: "Yep"}
	bytes, err := formatter.WriteEvent(evt)

	assert.Equal(t, "[1,\"i\",\"Yep\"]\n", string(bytes))
	assert.Nil(t, err)
}

func TestAddNewLine(t *testing.T) {
	text := []byte("Yo")
	var err error
	expected := []byte("Yo\n")

	actVal, actErr := AddNewline(text, err)
	assert.Equal(t, actVal, expected, "Adds a new line when no error is present")
	assert.Nil(t, actErr, "Does not add an error")

	err = fmt.Errorf("boo")
	actVal, actErr = AddNewline(text, err)
	assert.Equal(t, actVal, text, "Does nothing if error is present")
	assert.NotNil(t, actErr, "leaves error in place")
}

func TestAsciiCastFormatterWriteHeaderPlain(t *testing.T) {
	clock := clockwork.NewFakeClock()
	formatter := asciiCast{clock: clock}
	unixNow := clock.Now().Unix()

	systemstate.UpdateTermHeight('h') //converts to ascii code
	systemstate.UpdateTermWidth('w')

	bytes, err := formatter.WriteHeader(Metadata{})

	assert.Nil(t, err)
	assert.Equal(t, uint8('\n'), bytes[len(bytes)-1])

	bytes = bytes[:len(bytes)-1] // trim last byte
	var header ASCIICastHeader
	err = json.Unmarshal(bytes, &header)
	assert.Nil(t, err)

	assert.Equal(t, header.Version, 2)
	assert.Equal(t, header.Timestamp, int64(0))
	assert.Equal(t, header.Title, strconv.FormatInt(unixNow, 10))
	assert.Equal(t, header.Width, uint16('w'))
	assert.Equal(t, header.Height, uint16('h'))
	assert.Equal(t, header.Duration, float64(0))
	assert.Equal(t, header.Env, map[string]string{"SHELL": "", "TERM": ""})
}

func TestAsciiCastFormatterWriteHeaderWithContent(t *testing.T) {
	clock := clockwork.NewFakeClock()
	formatter := asciiCast{clock: clock}

	systemstate.UpdateTermHeight('h') //converts to ascii code
	systemstate.UpdateTermWidth('w')

	bytes, err := formatter.WriteHeader(Metadata{
		StartTimeUnix:   1234,
		DurationSeconds: 10,
		Title:           "SomeTitle",
		Shell:           "/bin/bash",
		Term:            "otherTerm",
	})

	assert.Nil(t, err)
	assert.Equal(t, uint8('\n'), bytes[len(bytes)-1])

	bytes = bytes[:len(bytes)-1] // trim last byte
	var header ASCIICastHeader
	err = json.Unmarshal(bytes, &header)
	assert.Nil(t, err)

	assert.Equal(t, header.Version, 2)
	assert.Equal(t, header.Timestamp, int64(1234))
	assert.Equal(t, header.Title, "SomeTitle")
	assert.Equal(t, header.Width, uint16('w'))
	assert.Equal(t, header.Height, uint16('h'))
	assert.Equal(t, header.Duration, float64(10))
	assert.Equal(t, header.Env, map[string]string{"SHELL": "/bin/bash", "TERM": "otherTerm"})
}
