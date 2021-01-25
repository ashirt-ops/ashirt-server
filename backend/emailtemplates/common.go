package emailtemplates

import (
	"text/template"

	"github.com/theparanoids/ashirt-server/backend/models"

	"github.com/theparanoids/ashirt-server/backend/authschemes/recoveryauth"
	"github.com/theparanoids/ashirt-server/backend/database"
)

const appFrontendRoot = "http://localhost:8080" // TODO: we should get this fron envvars

type EmailTemplateData struct {
	SendToAddress string
	UserRecord    *models.User
	DB            *database.Connection
}

var templateFuncs = template.New("base").Funcs(template.FuncMap{
	"AddRecoveryAuth": func(data EmailTemplateData, label string) (string, error) {
		// create recovery URL
		recoveryCode, err := recoveryauth.GenerateRecoveryCodeForUser(data.DB, data.UserRecord.ID)
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
