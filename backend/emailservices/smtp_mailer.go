package emailservices

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"

	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/logging"
)

// SMTPEmailAuthType indicates how the system should authenticate with the STMP server
// see: https://www.samlogic.net/articles/smtp-commands-reference-auth.htm
type SMTPEmailAuthType string

const (
	// LoginType indicates the login SMTP authentication flow
	LoginType SMTPEmailAuthType = "login"
	// PlainType indicates the plain SMTP authentication flow
	PlainType SMTPEmailAuthType = "plain"
	// CRAMMD5Type  indicates the CRAM-MD5 SMTP authentication flow
	CRAMMD5Type SMTPEmailAuthType = "crammd5"
)

// SMTPMailer is the struct that holds an email servicer that sends emails over SMTP
type SMTPMailer struct {
	logger logging.Logger
}

// MakeSMTPMailer constructs an SMTPMailer with the given logger
func MakeSMTPMailer(logger logging.Logger) SMTPMailer {
	return SMTPMailer{
		logger: logger,
	}
}

func (m *SMTPMailer) auth() smtp.Auth {
	switch config.EmailSMTPAuthType() {
	case string(LoginType):
		return smtpLoginAuth(config.EmailUserName(), config.EmailPassword())
	case string(PlainType):
		return smtp.PlainAuth(config.EmailIdentity(), config.EmailUserName(),
			config.EmailPassword(), config.EmailHost())
	case string(CRAMMD5Type):
		return smtp.CRAMMD5Auth(config.EmailUserName(), config.EmailSecret())
	}
	return nil
}

// AddToQueue attempts to send the provided email over smtp
func (m *SMTPMailer) AddToQueue(job EmailJob) error {
	msg, err := buildSMTPEmail(job)
	if err != nil {
		return err
	}
	err = smtp.SendMail(config.EmailHost(), m.auth(), job.From, []string{job.To}, msg)
	job.OnCompleted(err)
	return err
}

func buildSMTPEmail(job EmailJob) ([]byte, error) {
	smtpEmailContent, err := buildEmailContent(job)

	if err != nil {
		return []byte{}, err
	}

	return []byte(smtpEmailContent), nil
}

func buildEmailContent(job EmailJob) (string, error) {
	emailMessage := &bytes.Buffer{}
	emailMessageWriter := multipart.NewWriter(emailMessage)

	boundary := emailMessageWriter.Boundary()

	// check errs
	err := buildPlaintextPart(emailMessageWriter, job.Body)
	if err != nil {
		return "", err
	}
	err = buildHTMLPart(emailMessageWriter, job.HTMLBody)
	if err != nil {
		return "", err
	}
	err = emailMessageWriter.Close()
	if err != nil {
		return "", err
	}

	start := buildEmailHeader(job, boundary)
	return start + emailMessage.String(), nil
}

func buildPlaintextPart(partWriter *multipart.Writer, content string) error {
	return buildEmailBodyPart(partWriter, content, `text/plain; charset="utf-8"`)
}

func buildHTMLPart(partWriter *multipart.Writer, content string) error {
	return buildEmailBodyPart(partWriter, content, `text/html; charset="utf-8"`)
}

func buildEmailBodyPart(partWriter *multipart.Writer, content string, contentType string) error {
	mimeHeader := textproto.MIMEHeader{"Content-Type": {contentType}}
	mimeHeader.Add("Content-Transfer-Encoding", "quoted-printable")
	mimeHeader.Add("Content-Disposition", "inline")
	childWriter, err := partWriter.CreatePart(mimeHeader)
	if err != nil {
		return err
	}
	_, err = childWriter.Write([]byte(content))
	return err
}

func buildEmailHeader(job EmailJob, boundary string) string {
	to := fmt.Sprintf("To: %v\r\n", job.To)
	from := job.From
	if from == "" {
		from = "AShirt"
	}
	from = fmt.Sprintf("From: %v\r\n", from)
	subject := fmt.Sprintf("Subject: %v\r\n", job.Subject)
	mimeVersion := "MIME-Version: 1.0\r\n"
	contentType := fmt.Sprintf("Content-Type: multipart/alternative; boundary=%v\r\n", boundary)

	return to + from + subject + mimeVersion + contentType
}
