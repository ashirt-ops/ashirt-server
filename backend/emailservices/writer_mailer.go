package emailservices

import (
	"io"
	"text/template"

	"github.com/theparanoids/ashirt-server/backend/logging"
)

// WriterMailer acts as a no-email email server for monitoring when running locally. All emails are
// printed to the provided writer. This is intended to be used with os.Stdout, but any writer should
// work
type WriterMailer struct {
	writer io.Writer
	logger logging.Logger
}

// MakeWriterMailer constructs a WriterMailer
func MakeWriterMailer(w io.Writer, logger logging.Logger) WriterMailer {
	return WriterMailer{
		writer: w,
		logger: logger,
	}
}

// AddToQueue writes the email to the writer provided in MakeWriterMailer
// returns an error if the underlying template cannot be executed
func (m *WriterMailer) AddToQueue(job EmailJob) error {
	err := m.yellEmail(job)
	job.OnCompleted(err)

	return err
}

func (m *WriterMailer) yellEmail(job EmailJob) error {
	return emailPrintFormat.Execute(m.writer, job)
}

var emailPrintFormat = template.Must(template.New("yell-email").Parse(
	"\n" +
		`======================================================
To:      | {{ .To }}
From:    | {{ .From }}
Subject: | {{ .Subject }}
Body:
{{ .Body }}
` +
		"======================================================\n"))
