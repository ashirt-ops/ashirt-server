// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package config

import (
	"strconv"
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
	FrontendIndexURL         string        `split_words:"true"`
	BackendURL               string        `split_words:"true"`
	SuccessRedirectURL       string        `split_words:"true"`
	FailureRedirectURLPrefix string        `split_words:"true"`
	UseLambdaRIE             bool          `split_words:"true"`
	Flags                    string
	Port                     int
}

// DBConfig provides configuration details on connecting to the backend database
type DBConfig struct {
	URI string `required:"true"`
}

type ContentStoreConfig struct {
	Type   string `split_words:"true"`
	Bucket string `split_words:"true"`
	Region string `split_words:"true"`
}

var (
	app   WebConfig
	db    DBConfig
	auth  AuthConfig
	email EmailConfig
	store ContentStoreConfig
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
		loadEmailConfig,
		loadStoreConfig,
	})
}

// LoadAPIConfig loads all of the environment configuration from environment variables. This version
// exists primarily to load data for the "API" / tools server
func LoadAPIConfig() error {
	return loadConfig([]func() error{
		loadAppConfig,
		loadDBConfig,
		loadStoreConfig,
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

func loadStoreConfig() error {
	config := ContentStoreConfig{}
	err := envconfig.Process("store", &config)
	store = config

	return err
}

// DBUri retrieves the environment variable DB_URI
func DBUri() string {
	return db.URI
}

// ImageStoreBucketName retrieves the APP_IMGSTORE_BUCKET_NAME value from the environment
func ImageStoreBucketName() string {
	return app.ImgstoreBucketName
}

// AWSRegion retrieves the APP_IMGSTORE_REGION value from the environment
func AWSRegion() string {
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

// FrontendIndexURL retrieves the APP_FRONTEND_INDEX_URL value from the environment
func FrontendIndexURL() string {
	return app.FrontendIndexURL
}

func AllAppConfig() WebConfig {
	return app
}

func AllStoreConfig() ContentStoreConfig {
	return store
}

func StoreType() string {
	return store.Type
}

func StoreBucket() string {
	return store.Bucket
}

func StoreRegion() string {
	return store.Region
}

func UseLambdaRIE() bool {
	return app.UseLambdaRIE
}
