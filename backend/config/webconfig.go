// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// WebConfig is a namespaced app-specific configuration.
type WebConfig struct {
	ImgstoreBucketName       string        `split_words:"true"`
	ImgstoreRegion           string        `split_words:"true"`
	CsrfAuthKey              string        `split_words:"true"`
	SessionStoreKey          string        `split_words:"true"`
	RecoveryExpiry           time.Duration `split_words:"true" default:"24h"`
	DisableLocalRegistration bool          `split_words:"true"`
	Port                     int
}

// DBConfig provides configuration details on connecting to the backend database
type DBConfig struct {
	URI string `required:"true"`
}

// AuthConfig provides configuration details, namespaced to
type AuthConfig struct {
	Services    []string
	AuthConfigs map[string]AuthInstanceConfig
}

// AuthInstanceConfig provides all of the _possible_ configuration values for an auth instance.
// Note: it is expected that not all fields will be populated. It is up to the user to verify
// that these fields exist and have correct values
type AuthInstanceConfig struct {
	ClientID                 string `split_words:"true"`
	ClientSecret             string `split_words:"true"`
	Issuer                   string `split_words:"true"`
	BackendURL               string `split_words:"true"`
	SuccessRedirectURL       string `split_words:"true"`
	FailureRedirectURLPrefix string `split_words:"true"`
	ProfileToShortnameField  string `split_words:"true"`
}

var (
	app  WebConfig
	db   DBConfig
	auth AuthConfig
)

// LoadConfig loads all of the environment configuration specified in environment variables
func loadConfig(funcs []func() error) error {
	var err error

	for _, f := range funcs {
		err = f()
		if err != nil {
			break
		}
	}

	return err
}

// LoadWebConfig loads all of the environment configuration from environment variables. This version
// exists primarily to load data for the "Web" / UI server
func LoadWebConfig() error {
	return loadConfig([]func() error{
		loadAppConfig,
		loadDBConfig,
		loadAuthConfig,
	})
}

// LoadAPIConfig loads all of the environment configuration from environment variables. This version
// exists primarily to load data for the "API" / tools server
func LoadAPIConfig() error {
	return loadConfig([]func() error{
		loadAppConfig,
		loadDBConfig,
	})
}

func loadAppConfig() error {
	config := WebConfig{}
	err := envconfig.Process("app", &config)
	app = config

	return err
}

func loadDBConfig() error {
	config := DBConfig{}
	err := envconfig.Process("DB", &config)
	db = config

	return err
}

func loadAuthConfig() error {
	config := AuthConfig{}
	servicesStr := os.Getenv("AUTH_SERVICES")
	if servicesStr == "" {
		return errors.New("Auth services not defined")
	}

	servicesArr := strings.Split(servicesStr, ",")
	for i, s := range servicesArr {
		servicesArr[i] = strings.TrimSpace(s)
	}

	config.Services = servicesArr
	config.AuthConfigs = make(map[string]AuthInstanceConfig)
	for _, service := range servicesArr {
		innerConfig := AuthInstanceConfig{}
		err := envconfig.Process("auth_"+service, &innerConfig)
		if err != nil {
			return err
		}
		config.AuthConfigs[service] = innerConfig
	}
	auth = config

	return nil
}

// DBUri retrieves the environment variable DB_URI
func DBUri() string {
	return db.URI
}

// AuthConfigInstance attempts to retrieve a particular auth configuration set from the environment.
// Note this looks for environment variables prefixed with AUTH_${SERVICE_NAME}, and will only
// retrieve these values for services named in the AUTH_SERVICES environment
func AuthConfigInstance(name string) AuthInstanceConfig {
	v, _ := auth.AuthConfigs[name]
	return v
}

// SupportedAuthServices retrieves the parsed AUTH_SERVICES value from the environment
func SupportedAuthServices() []string {
	return auth.Services
}

// ImageStoreBucketName retrieves the APP_IMGSTORE_BUCKET_NAME value from the environment
func ImageStoreBucketName() string {
	return app.ImgstoreBucketName
}

// ImageStoreRegion retrieves the APP_IMGSTORE_REGION value from the environment
func ImageStoreRegion() string {
	return app.ImgstoreRegion
}

// CSRFAuthKey retrieves the APP_CSRF_AUTH_KEY value from the environment
func CSRFAuthKey() string {
	return app.CsrfAuthKey
}

// SessionStoreKey retrieves the SESSION_STORE_KEY value from the environment
func SessionStoreKey() string {
	return app.SessionStoreKey
}

// Port retrieves the APP_PORT value from the environment
func Port() string {
	return strconv.Itoa(app.Port)
}

// RecoveryExpiry retrieves the APP_RECOVERY_EXPIRY value from the environment
func RecoveryExpiry() time.Duration {
	return app.RecoveryExpiry
}

// IsRegistrationEnabled returns true if local registration is enabled, false otherwise.
func IsRegistrationEnabled() bool {
	return !app.DisableLocalRegistration
}
