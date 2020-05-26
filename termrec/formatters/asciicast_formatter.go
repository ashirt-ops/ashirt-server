package formatters

import (
	"encoding/json"
	"strconv"

	"github.com/jonboulle/clockwork"
	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/systemstate"
)

type asciiCast struct {
	clock clockwork.Clock
}

// ASCIICast is an empty struct to help format into Asciicast/Asciinema format. See
// https://github.com/asciinema/asciinema/blob/develop/doc/asciicast-v2.md for more details on the
// format
var ASCIICast = asciiCast{
	clock: clockwork.NewRealClock(),
}

// WriteHeader constructs an Asciicinema/Asciicast header. A basic header is constructed, and if
// information is present in the recorder, more details can be added.
func (f asciiCast) WriteHeader(m Metadata) ([]byte, error) {
	title := m.Title
	if title == "" {
		title = strconv.FormatInt(f.clock.Now().Unix(), 10)
	}

	header := ASCIICastHeader{
		Version:   2,
		Width:     systemstate.TermWidth(),
		Height:    systemstate.TermHeight(),
		Timestamp: m.StartTimeUnix,
		Env:       map[string]string{"SHELL": m.Shell, "TERM": m.Term},
		Title:     title,
	}

	if m.DurationSeconds != 0 {
		header.Duration = m.DurationSeconds
	}

	if theme() != nil {
		header.Theme = theme()
	}

	return AddNewline(json.Marshal(header))
}

// AddNewline is a small helper function to append a newline character (\n) to the end of a byte
// slice, if the provided error is not nil. The input is the same as the output for json.Marshal
// so that they can be chained easily.
func AddNewline(bytes []byte, err error) ([]byte, error) {
	if err != nil {
		return bytes, err
	}
	bytes = append(bytes, byte('\n'))
	return bytes, nil
}

// WriteEvent serializes an event into a json array, terminated by a new line (\n)
func (f asciiCast) WriteEvent(evt common.Event) ([]byte, error) {
	encodedEvent := []interface{}{evt.When.Seconds(), evt.Type, evt.Data}
	return AddNewline(json.Marshal(encodedEvent))
}

// WriteFooter is a no-op here, as the asciicast file format does not contain a footer
func (f asciiCast) WriteFooter(m Metadata) ([]byte, error) {
	return []byte{}, nil
}

// TODO: Not sure how to really get this information, nor where it fits in structurly
func theme() *ASCIICastTheme { return nil }
