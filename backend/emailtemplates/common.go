package emailtemplates

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/jaytaylor/html2text"
	recoveryHelpers "github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth/helpers"
	"github.com/theparanoids/ashirt-server/backend/config"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/models"
)

// EmailTemplate is an enum describing each of the possible email types
type EmailTemplate = string

const (
	// EmailRecoveryTemplate contains a message indicating that a user can recover their account
	EmailRecoveryTemplate EmailTemplate = "self-service-recovery-email"

	// EmailRecoveryDeniedTemplate contains a message indicating that a user CANNOT recover their
	// account because it's disabled
	EmailRecoveryDeniedTemplate EmailTemplate = "self-service-recovery-denied-email"
)

type EmailTemplateData struct {
	UserRecord *models.User
	DB         *database.Connection
}

var templateFuncs = template.New("base").Funcs(template.FuncMap{
	"RecoveryURL": func(data EmailTemplateData) (string, error) { // create recovery URL
		appFrontendRoot := config.FrontendIndexURL()
		recoveryCode, err := recoveryHelpers.GenerateRecoveryCodeForUser(data.DB, data.UserRecord.ID)
		recoveryURL := appFrontendRoot + "/web/auth/recovery/login?code=" + recoveryCode
		return recoveryURL, err
	},
	"FullName": func(data EmailTemplateData) string {
		if data.UserRecord != nil {
			return data.UserRecord.FirstName + " " + data.UserRecord.LastName
		}
		return ""
	},
})

type EmailContent struct {
	Subject           string
	PlaintTextContent string
	HTMLContent       string
}

// BuildEmailContent constructs an email subject and body from the given template and template data. Returns
// (body, subject, nil) if there is no error generating the content, otherwise ("", "", error)
func BuildEmailContent(emailTemplate EmailTemplate, templateData EmailTemplateData) (EmailContent, error) {
	w := bytes.NewBuffer(make([]byte, 0))
	var err error
	rtn := EmailContent{}

	switch emailTemplate {
	case EmailRecoveryTemplate:
		err = recoveryEmail.Execute(w, templateData)
		rtn.Subject = "Recover your AShirt account"
	case EmailRecoveryDeniedTemplate:
		err = recoveryDeniedDisabledEmail.Execute(w, templateData)
		rtn.Subject = "Recover your AShirt account"
	default:
		err = errors.New("unsupported email template")
	}

	if err != nil {
		return EmailContent{}, err
	}
	rtn.HTMLContent = w.String()
	rtn.PlaintTextContent, err = html2text.FromString(rtn.HTMLContent)

	if err != nil {
		return EmailContent{}, nil
	}

	return rtn, nil
}
