package emailtemplates

import (
	"text/template"
)

const recoveryParagraph = "When you click on the link above, " +
	"you will be directed to the account management area. " +
	"If you have permanently lost the ability to log in normally, " +
	"you may need to unlink, and then re-link your authentication method. " +
	"Here are the steps to do that:\n" +
	"  1. Find your normal authentication method under the \"Authentication Methods\" header.\n" +
	"  2. Press the delete button next to your normal authentication method.\n" +
	"  3. You will be prompted to confirm the deletion by entering the name of the authentication " +
	"method. Do this, then press the \"Delete\" button.\n" +
	"  4. Once again, find your normal authentication method under the \"Authentication Methods\" " +
	"header.\n" +
	"  5. Click the \"Link\" button and, depending on the authentication mechanism, choose new values " +
	"for your authentication."

var recoveryEmail = template.Must(templateFuncs.New("recoveryEmail").Parse(
	`Hi {{ FullName . }},

A request was made to recover your account. You can recover your account by {{ AddRecoveryAuth . "Cliking Here" }}.

` + recoveryParagraph + `

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
