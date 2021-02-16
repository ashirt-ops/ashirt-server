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
	SendToAddress string
	UserRecord    *models.User
	DB            *database.Connection
}

type OutgoingEmail struct {
	Body  string
	To    string
	Error error
}

var templateFuncs = template.New("base").Funcs(template.FuncMap{
	"AddRecoveryAuth": func(data EmailTemplateData, label string) (string, error) {
		// create recovery URL
		appFrontendRoot := config.FrontendIndexURL()
		recoveryCode, err := recoveryHelpers.GenerateRecoveryCodeForUser(data.DB, data.UserRecord.ID)
		recoveryURL := appFrontendRoot + "/login/recovery/" + recoveryCode

		return `<a href="` + recoveryURL + `">` + label + `</a>`, err
	},
	"FullName": func(data EmailTemplateData) string {
		if data.UserRecord != nil {
			return data.UserRecord.FirstName + " " + data.UserRecord.LastName
		}
		return ""
	},
})

// BuildEmail constructs an email message from the given template and template data
func BuildEmail(emailTemplate EmailTemplate, templateData EmailTemplateData) OutgoingEmail {
	buf := make([]byte, 0)
	w := bytes.NewBuffer(buf)
	var err error

	switch emailTemplate {
	case EmailRecoveryTemplate:
		err = recoveryEmail.Execute(w, templateData)
	case EmailRecoveryDeniedTemplate:
		err = recoveryDeniedDisabledEmail.Execute(w, templateData)
	}

	rtn := OutgoingEmail{}
	if err != nil {
		rtn.Error = err
	} else {
		rtn.Body = string(buf)
		rtn.To = templateData.UserRecord.Email
	}

	return rtn
}
