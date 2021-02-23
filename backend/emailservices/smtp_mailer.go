package emailservices

import (
	"bytes"
	"net/smtp"
	"text/template"

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
	return err
}

func buildSMTPEmail(job EmailJob) ([]byte, error) {
	buff := bytes.NewBuffer(make([]byte, 0, 1024))
	err := emailSMTPFormat.Execute(buff, job)
	if err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

var emailSMTPFormat = template.Must(template.New("smtp-email").Parse(
	"To: {{ .To }}\r\n" +
		"Subject: | {{ .Subject }}\r\n" +
		"\r\n{{ .Body }}\r\n"))
