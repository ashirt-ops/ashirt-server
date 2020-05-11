package main

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/theparanoids/ashirt/termrec/systemstate"
)

// PtyTracker is here to help collect all of the pty related items that need to be passed around
type PtyTracker struct {
	Pty                *os.File
	WindowListenerChan chan os.Signal
	termOut            io.Writer
	readOut            io.Writer
	termInput          io.Reader
	startupScript      string
	OnReady            func()
}

// NewPtyTracker generates an initial tracker.
// termOut io.Writer // The writer that will be used to output events to what would be considered stdout
//                   // Note: In many cases, you will want to pass a multiplexed writer (see io.MultiWriter)
// readOut io.Writer // A writer that will forward along stdin events. This is the result of calling io.TeeReader.
//                   // Note: if you are not interested in receiving stdin events, send along a io.Discard / no-op writer
// onReady func()    // A hook into the run state that provides a pre-recording area to display messages
// startupScript string // The path to the (optional) startup script. Ignored if startupScript is an empty string
//
// Note: This method assumes all inputs come from stdin. This will adjust the stdin on the calling terminal
// We will attempt to restore to the original state on exit.
func NewPtyTracker(termOut, readOut io.Writer, termInput io.Reader, onReadyFunc func(), startupScript string) PtyTracker {
	return PtyTracker{
		OnReady:       onReadyFunc,
		termOut:       termOut,
		readOut:       readOut,
		termInput:     termInput,
		startupScript: startupScript,
	}
}

// Run starts the pty session
func (t *PtyTracker) Run(shell string) error {
	defer t.close()

	c := exec.Command(shell)
	var err error
	t.Pty, err = pty.Start(c)
	if err != nil {
		return err
	}

	if t.startupScript != "" {
		t.Pty.WriteString(t.startupScript + "\n")
	}

	t.WindowListenerChan = startResizeListener(t.Pty)

	t.OnReady()

	wrappedStdin := io.TeeReader(t.termInput, t.readOut)
	go func() { io.Copy(t.Pty, wrappedStdin) }()
	io.Copy(t.termOut, t.Pty)

	return nil
}

// Close performs all of the closes necessary to restore the system back to a good state.
func (t *PtyTracker) close() {
	if t.WindowListenerChan != nil {
		close(t.WindowListenerChan)
	}
	if t.Pty != nil {
		t.Pty.Close()
	}
}

func startResizeListener(ptmx *os.File) chan os.Signal {
	// TODO: this expects stdin, but everything else here does not; we should move this elsewhere
	// (perhaps provide a function generate this logic, given the ptmx?)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			size, err := pty.GetsizeFull(os.Stdin)
			if err == nil {
				systemstate.UpdateTermHeight(size.Rows)
				systemstate.UpdateTermWidth(size.Cols)
				pty.InheritSize(os.Stdin, ptmx)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.
	return ch
}
