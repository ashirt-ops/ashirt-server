package write

import (
	"bufio"
	"io"

	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/formatters"
	"github.com/theparanoids/ashirt/termrec/isthere"
)

// StreamingFileWriter is the preferred TerminalWriter that writes to a file as soon as
// an event/header/footer comes in. The file is kept open until Close is called.
//
// Note: This currently ignores any write errors.
type StreamingFileWriter struct {
	outStream   io.Writer
	backingFile FileLike
	formatter   formatters.Formatter
}

// NewStreamingFileWriter generates a new StreamingFileWriter, which can be used as a TerminalWriter
// Note: this will open up a file (at the provided path, or if nil, as a temporary file)
func NewStreamingFileWriter(filedir, filename string, formatter formatters.Formatter, buffered bool) (StreamingFileWriter, error) {
	file, err := NewFile(filedir, filename)

	var writer io.Writer = file
	if buffered {
		writer = bufio.NewWriter(file)
	}

	return StreamingFileWriter{
		formatter:   formatter,
		backingFile: file,
		outStream:   writer,
	}, err
}

// WriteHeader attempts to write a header (per the provided formatter) directly to the backing file.
// Note that since this is streamed, this must be called before footer or event (i.e. first)
func (fw StreamingFileWriter) WriteHeader(m formatters.Metadata) {
	encoded, err := fw.formatter.WriteHeader(m)
	if len(encoded) > 0 {
		noErrorWrite(encoded, err, fw.outStream.Write)
	}
}

// WriteFooter attempts to write a footer (per the provided formatter) directly to the backing file.
// Note that since this is streamed, this must be called after header and all events (i.e. last)
func (fw StreamingFileWriter) WriteFooter(m formatters.Metadata) {
	encoded, err := fw.formatter.WriteFooter(m)
	if len(encoded) > 0 {
		noErrorWrite(encoded, err, fw.outStream.Write)
	}
}

// WriteEvent attempts to write out a single event to the stream, per the provided formatter.
func (fw StreamingFileWriter) WriteEvent(evt common.Event) {
	encoded, err := fw.formatter.WriteEvent(evt)
	noErrorWrite(encoded, err, fw.outStream.Write)
}

// Close is a required call to both flush the buffered writer (if signaled via the constructor/hand created)
// and to close the file itself.
func (fw StreamingFileWriter) Close() error {
	if w, ok := fw.outStream.(*bufio.Writer); ok {
		w.Flush()
	}
	return fw.backingFile.Close()
}

// Filepath retrieves the path to the streamed file.
func (fw StreamingFileWriter) Filepath() string {
	return fw.backingFile.Name()
}

// noErrorWrite actually just passes the errors along now.
func noErrorWrite(bytes []byte, err error, Write func(b []byte) (n int, err error)) {
	if isthere.No(err) {
		Write(bytes)
	}
}
