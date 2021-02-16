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

var recoveryDeniedDisabledEmail = template.Must(templateFuncs.New("recoveryDeniedEmail").Parse(
	`Hi {{ FullName . }},

A request was made to recover your account. Unfortunately, this account has been disabled. Please first contact an administrator to restore functionality.

If you did not make this request, you can ignore this email.

Thanks,

The ASHIRT Team
`,
))
