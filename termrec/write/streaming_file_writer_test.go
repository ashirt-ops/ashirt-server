package write

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/formatters"
	"github.com/stretchr/testify/assert"
)

func makeTestStreamingFileWriter() (StreamingFileWriter, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	writer := StreamingFileWriter{
		outStream: buf,
		formatter: PlainFormatter{},
	}
	return writer, buf
}

func TestStreamingFileWriterWriteHeader(t *testing.T) {
	writer, buf := makeTestStreamingFileWriter()

	m := formatters.Metadata{Title: "bonanza"}

	writer.WriteHeader(m)

	assert.Equal(t, m.String(), string(buf.Bytes()))
}

func TestStreamingFileWriterWriteFooter(t *testing.T) {
	writer, buf := makeTestStreamingFileWriter()

	m := formatters.Metadata{Title: "bonanza"}

	writer.WriteFooter(m)

	assert.Equal(t, m.String(), string(buf.Bytes()))
}

func TestStreamingFileWriterWriteEvent(t *testing.T) {
	writer, buf := makeTestStreamingFileWriter()

	evt1 := common.Event{Type: "i", Data: "someData", When: 2 * time.Second}
	evt2 := common.Event{Type: "o", Data: "moreData", When: 4 * time.Second}

	writer.WriteEvent(evt1)
	assert.Equal(t, []byte(evt1.String()), buf.Bytes())
	buf.Reset()

	writer.WriteEvent(evt2)
	assert.Equal(t, []byte(evt2.String()), buf.Bytes())
}

func TestStreamingFileWriterFilepath(t *testing.T) {
	writer := StreamingFileWriter{
		backingFile: &DummyFile{DummyPath: "/path/to/file"},
	}

	assert.Equal(t, writer.Filepath(), "/path/to/file")
}

func TestStreamingFileWriterCloseNonBuffered(t *testing.T) {
	f := DummyFile{DummyPath: "/path/to/file"}
	writer := StreamingFileWriter{
		backingFile: &f,
	}
	writer.Close()

	assert.True(t, f.IsClosed)
}

func TestStreamingFileWriterCloseBuffered(t *testing.T) {
	f := DummyFile{DummyPath: "/path/to/file"}
	buf := new(bytes.Buffer)
	reBuf := bufio.NewWriter(buf)
	writer := StreamingFileWriter{
		backingFile: &f,
		outStream:   reBuf,
		formatter:   PlainFormatter{},
	}
	sampleMetadata := formatters.Metadata{Title: "Boo!"}
	assert.Equal(t, buf.Len(), 0)
	writer.WriteHeader(sampleMetadata)
	assert.Equal(t, buf.Len(), 0, "Check that buffered writer has actually buffered some data (not important)")
	writer.Close()
	assert.Equal(t, buf.Len(), len([]byte(sampleMetadata.String())), "Check that BufferedWriter has been flushed on close")
}

type DummyFile struct {
	DummyPath string
	IsClosed  bool
}

func (d *DummyFile) Name() string {
	return d.DummyPath
}

func (d *DummyFile) Close() error {
	d.IsClosed = true
	return nil
}

type PlainFormatter struct {
}

func (p PlainFormatter) WriteEvent(evt common.Event) ([]byte, error) {
	return []byte(evt.String()), nil
}
func (p PlainFormatter) WriteFooter(m formatters.Metadata) ([]byte, error) {
	return []byte(m.String()), nil

}
func (p PlainFormatter) WriteHeader(m formatters.Metadata) ([]byte, error) {
	return []byte(m.String()), nil
}
