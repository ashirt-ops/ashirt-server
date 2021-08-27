package config

import "github.com/kelseyhightower/envconfig"

// EmailConfig is a struct that houses the configuration details related specifically to the (optional)
// email services
type EmailConfig struct {
	FromAddress  string `split_words:"true"`
	Type         string `split_words:"true"`
	Host         string `split_words:"true"`
	UserName     string `split_words:"true"`
	Password     string `split_words:"true"`
	Identity     string `split_words:"true"`
	Secret       string `split_words:"true"`
	SMTPAuthType string `split_words:"true"`
}

func loadEmailConfig() error {
	config := EmailConfig{}
	err := envconfig.Process("email", &config)
	email = config

	return err
}

// EmailFromAddress contains the from address for all outgoing emails
func EmailFromAddress() string {
	return email.FromAddress
}

// EmailType contains the type of email servicer to use (e.g. STMP vs memory-based solution)
func EmailType() string {
	return email.Type
}

// EmailHost contains the location of the (presumably smtp) host that the email servicer will connect to
func EmailHost() string {
	return email.Host
}

// EmailUserName contains the "username" part of the information needed to authenticate with the email host
func EmailUserName() string {
	return email.UserName
}

// EmailPassword contains the "password" part of the information needed to authenticate with the email host
func EmailPassword() string {
	return email.Password
}

// EmailIdentity contains the identity feature when using an Plain SMTP authentication
func EmailIdentity() string {
	return email.Identity
}

// EmailSMTPAuthType contains the option for how to authenticate with the smtp service
func EmailSMTPAuthType() string {
	return email.SMTPAuthType
}

// EmailSecret contains the secret needed when using CRAMMD5 SMTP authentication
func EmailSecret() string {
	return email.Secret
}
