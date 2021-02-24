package emailtemplates

import (
	"bytes"
	"text/template"

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
	"AddRecoveryAuth": func(data EmailTemplateData, label string) (string, error) {
		// create recovery URL
		appFrontendRoot := config.FrontendIndexURL()
		recoveryCode, err := recoveryHelpers.GenerateRecoveryCodeForUser(data.DB, data.UserRecord.ID)
		recoveryURL := appFrontendRoot + "/web/auth/recovery/login?code=" + recoveryCode

		return `<a href="` + recoveryURL + `">` + label + `</a>`, err
	},
	"FullName": func(data EmailTemplateData) string {
		if data.UserRecord != nil {
			return data.UserRecord.FirstName + " " + data.UserRecord.LastName
		}
		return ""
	},
})

// BuildEmailContent constructs an email subject and body from the given template and template data. Returns
// (body, subject, nil) if there is no error generating the content, otherwise ("", "", error)
func BuildEmailContent(emailTemplate EmailTemplate, templateData EmailTemplateData) (string, string, error) {
	w := bytes.NewBuffer(make([]byte, 0))
	var err error
	subject := ""

	switch emailTemplate {
	case EmailRecoveryTemplate:
		err = recoveryEmail.Execute(w, templateData)
		subject = "Recover your AShirt account"
	case EmailRecoveryDeniedTemplate:
		err = recoveryDeniedDisabledEmail.Execute(w, templateData)
		subject = "Recover your AShirt account"
	}

	if err != nil {
		return "", "", err
	}
	return string(w.Bytes()), subject, nil
}
