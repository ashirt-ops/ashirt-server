package emailservices

import (
	"bytes"
	"net/smtp"
	"text/template"

	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/logging"
)

type SMTPEmailAuthType string

const (
	LoginType   SMTPEmailAuthType = "login"
	PlainType   SMTPEmailAuthType = "plain"
	CRAMMD5Type SMTPEmailAuthType = "cramda5"
)

type SMTPMailer struct {
	logger logging.Logger
}

func MakeSMTPMailer(logger logging.Logger) SMTPMailer {
	return SMTPMailer{
		logger: logger,
	}
}

func (m *SMTPMailer) Auth() smtp.Auth {
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

func (m *SMTPMailer) AddToQueue(job EmailJob) error {
	msg, err := buildSMTPEmail(job)
	if err != nil {
		return err
	}
	err = smtp.SendMail(config.EmailHost(), m.Auth(), job.From, []string{job.To}, msg)
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
