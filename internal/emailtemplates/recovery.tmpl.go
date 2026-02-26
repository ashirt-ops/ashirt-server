package emailtemplates

import (
	_ "embed"
	"text/template"
)

//go:embed recovery.html
var recoveryTemplate string

//go:embed recovery_denied.html
var recoveryDeniedTemplate string

var recoveryEmail = template.Must(templateFuncs.New("recoveryEmail").Parse(
	recoveryTemplate,
))

var recoveryDeniedDisabledEmail = template.Must(templateFuncs.New("recoveryDeniedEmail").Parse(
	recoveryDeniedTemplate,
))
