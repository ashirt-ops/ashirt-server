package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/theparanoids/ashirt/cmd/termrec/appdialogs"
	"github.com/theparanoids/ashirt/cmd/termrec/config"
	"github.com/theparanoids/ashirt/termrec/fancy"
	"github.com/theparanoids/ashirt/termrec/isthere"
	"github.com/theparanoids/ashirt/termrec/network"
	"github.com/theparanoids/ashirt/termrec/systemstate"
)

func main() {
	cfg := readConfig()

	network.SetBaseURL(cfg.APIURL)
	network.SetAccessKey(cfg.AccessKey)
	network.SetSecretKey(cfg.SecretKey)

	// switch to raw input, to stream to pty
	initialTerminalState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	defer func() {
		// defering in case we end up panicing (though we want to end this mode earlier)
		terminal.Restore(int(os.Stdin.Fd()), initialTerminalState)
	}()

	ptyReader, ptyWriter := io.Pipe()
	dialogReader, dialogWriter := io.Pipe()
	writeTarget := 0
	appdialogs.SetBasicUploadData(cfg.OperationID, dialogReader)

	recOpts := RecordingInput{
		FileDir:   cfg.OutputDir,
		FileName:  cfg.OutputFileName,
		Shell:     cfg.RecordingShell,
		TermInput: ptyReader,
		OnRecordingStart: func(output RecordingOutput) {
			fmt.Println("Recording to " + fancy.WithPizzazz(output.FilePath, fancy.Bold) + "\n\r")
			fmt.Println(fancy.WithPizzazz("Recording now live!\r", fancy.Bold|fancy.Reverse|fancy.LightGreen))
		},
	}

	if size, err := pty.GetsizeFull(os.Stdin); isthere.No(err) {
		systemstate.UpdateTermHeight(size.Rows)
		systemstate.UpdateTermWidth(size.Cols)
	}

	go func() {
		copyRouter([]io.Writer{ptyWriter, dialogWriter}, os.Stdin, &writeTarget)
	}()
	output, err := record(recOpts)
	terminal.Restore(int(os.Stdin.Fd()), initialTerminalState) // Restore terminal state for dialog

	if err != nil {
		fmt.Printf("%+v\n", err) // this will print a stacktrace, which may be better at this point
		return
	}
	writeTarget = 1

	appdialogs.SetDefaultData(appdialogs.UploadDefaults{FilePath: output.FilePath})
	appdialogs.ShowUploadMainMenu()

	fmt.Println("Bye!")
}

func readConfig() config.TermRecorderConfig {
	shouldContinue := true

	appConfig := config.Parse()
	issues := validateConfig(appConfig)

	allIssues := append(convertConfigIssuesToStartupIssues(appConfig), issues...)
	for _, iss := range allIssues {
		// coloring := fancy.Yellow
		if iss.Level == fatal {
			shouldContinue = false
			fmt.Println(fancy.Fatal("Fatal: "+iss.Msg, nil))
		} else {
			fmt.Println(fancy.Caution("Warning: "+iss.Msg, nil))
		}
	}
	if !shouldContinue {
		fmt.Println(fancy.WithPizzazz("Fatal issues encountered", fancy.Red|fancy.Reverse|fancy.Bold))
		fmt.Println("Please correct issues above, then rerun. Also note you can bring up the help menu with -h")
		os.Exit(-1)
	}

	return appConfig
}

func convertConfigIssuesToStartupIssues(cfg config.TermRecorderConfig) []StartupIssue {
	issues := make([]StartupIssue, len(cfg.Issues()))
	for i, msg := range cfg.Issues() {
		issues[i] = NewStartupIssue(warning, msg, true)
	}
	return issues
}

func validateConfig(cfg config.TermRecorderConfig) []StartupIssue {
	errors := make([]StartupIssue, 0)
	if cfg.APIURL == "" {
		errors = append(errors, NewStartupIssue(fatal, "API URL was not specified.", false))
	}

	if cfg.AccessKey == "" {
		errors = append(errors, NewStartupIssue(fatal, "Access key not set", false))
	} else if cfg.SecretKeyBase64 == "" {
		errors = append(errors, NewStartupIssue(fatal, "Secret key not set", false))
	} else {
		var err error
		cfg.SecretKey, err = base64.StdEncoding.DecodeString(cfg.SecretKeyBase64)
		if err != nil {
			errors = append(errors, NewStartupIssue(warning, "Unable to parse Secret Key. Please double check your configuration", false))
		}
	}

	return errors
}

type startupLevel string

const (
	warning startupLevel = "Warning"
	fatal   startupLevel = "Fatal"
)

// StartupIssue is a small structure for capturing issues during start up (in particular, reading the
// configuration data). Errors have two levels: Warnings and Fatals. Warning events signify that
// execution can continue, but with limitations (e.g. not being able to upload content) while Fatals
// signify that execution cannot continue.
// In addition to issue levels, this struct also exposes "Msg" which captures a friendly message,
// and CanUpload, which indicates to the underlying services if uploading can be done with the
// current configuration
type StartupIssue struct {
	Level     startupLevel
	Msg       string
	CanUpload bool
}

// NewStartupIssue is a constructor for StartupIssue.
func NewStartupIssue(lvl startupLevel, msg string, canUpload bool) StartupIssue {
	return StartupIssue{Level: lvl, Msg: msg, CanUpload: canUpload}
}

func copyRouter(dsts []io.Writer, src io.Reader, target *int) (written int64, err error) {
	size := 32 * 1024
	if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
		if l.N < 1 {
			size = 1
		} else {
			size = int(l.N)
		}
	}
	buf := make([]byte, size)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dsts[*target].Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
