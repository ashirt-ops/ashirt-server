package emailtemplates

import (
	"text/template"
)

var recoveryEmail = template.Must(templateFuncs.New("recoveryEmail").Parse(
	`Hi {{ FullName . }},

A request was made to recover your account. You can recover your account by {{ AddRecoveryAuth . "Cliking Here" }}.

If you did not make this request, you can ignore this email.

Thanks,

The ASHIRT Team
`,
))
