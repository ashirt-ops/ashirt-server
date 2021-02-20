package config

import "github.com/kelseyhightower/envconfig"

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

func EmailFromAddress() string {
	return email.FromAddress
}

func EmailType() string {
	return email.Type
}

func EmailHost() string {
	return email.Host
}

func EmailUserName() string {
	return email.UserName
}

func EmailPassword() string {
	return email.Password
}

func EmailIdentity() string {
	return email.Identity
}

func EmailSMTPAuthType() string {
	return email.SMTPAuthType
}

func EmailSecret() string {
	return email.Secret
}
