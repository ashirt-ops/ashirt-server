package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/jonboulle/clockwork"
	"github.com/pkg/errors"
	"github.com/theparanoids/ashirt/termrec/common"
	"github.com/theparanoids/ashirt/termrec/eventers"
	"github.com/theparanoids/ashirt/termrec/formatters"
	"github.com/theparanoids/ashirt/termrec/recorders"
	"github.com/theparanoids/ashirt/termrec/write"
)

// RecordingInput is a small structure for holding all configuration details for starting up
// a recording.
//
// This structure contains the following fields:
// FileName: The name of the file to be written
// FileDir: Where the file should be stored
// Shell: What shell to use for the PTY
// EventMiddleware: How to transform events that come through
// OnRecordingStart: A hook into the recording process just before actual recording starts
//   This is intended allow the user to provide messaging to the user
type RecordingInput struct {
	FileName         string
	FileDir          string
	Shell            string
	TermInput        io.Reader
	EventMiddleware  []eventers.EventMiddleware
	OnRecordingStart func(RecordingOutput)
}

// RecordingOutput is a small structure for communicating in-progress or completed recording details
type RecordingOutput struct {
	FilePath string
}

func record(ri RecordingInput) (RecordingOutput, error) {
	var result RecordingOutput
	tw, err := write.NewStreamingFileWriter(ri.FileDir, ri.FileName, formatters.ASCIICast, true)

	if err != nil {
		return result, errors.Wrap(err, "Unable to create file writer")
	}
	result.FilePath = tw.Filepath()

	recorder := recorders.NewStreamingRecorder(tw, clockwork.NewRealClock(), ri.Shell)
	eventWriter := eventers.NewEventWriter(&recorder, common.Output, ri.EventMiddleware...)
	wrappedStdOut := io.MultiWriter(os.Stdout, eventWriter)

	tracker := NewPtyTracker(wrappedStdOut, ioutil.Discard, ri.TermInput, func() { ri.OnRecordingStart(result) })

	err = tracker.Run(ri.Shell)
	if err != nil {
		return result, errors.Wrap(err, `Unable to start the recording. Shell path: "`+ri.Shell+`"`)
	}
	err = tw.Close()
	return result, errors.Wrap(err, "Issue closing file writer")
}
