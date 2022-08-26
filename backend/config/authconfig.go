// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package config

import (
	"errors"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/theparanoids/ashirt-server/backend/helpers"
)

// AuthConfig provides configuration details, namespaced to
type AuthConfig struct {
	Services    []string
	AuthConfigs map[string]AuthInstanceConfig
}

// AuthInstanceConfig provides all of the _possible_ configuration values for an auth instance.
// Note: it is expected that not all fields will be populated. It is up to the user to verify
// that these fields exist and have correct values
type AuthInstanceConfig struct {
	Type                string
	Name                string
	RegistrationEnabled bool `ignored:"true"`

	//generic oidc
	FriendlyName          string `split_words:"true"`
	ProviderURL           string `split_words:"true"`
	ClientID              string `split_words:"true"`
	ClientSecret          string `split_words:"true"`
	Scopes                string
	ProfileFirstNameField string `split_words:"true"`
	ProfileLastNameField  string `split_words:"true"`
	ProfileEmailField     string `split_words:"true"`
	ProfileSlugField      string `split_words:"true"`

	//webauthn
	DisplayName string `split_words:"true"`
	Timeout     string
	Conveyance  string
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

	serviceRegistrationStr := os.Getenv("AUTH_SERVICES_ALLOW_REGISTRATION")
	serviceRegistrationArr := strings.Split(serviceRegistrationStr, ",")

	config.Services = servicesArr
	config.AuthConfigs = make(map[string]AuthInstanceConfig)
	for _, service := range servicesArr {
		innerConfig := AuthInstanceConfig{}

		if helpers.ContainsMatch(serviceRegistrationArr, service) {
			innerConfig.RegistrationEnabled = true
		}
		err := envconfig.Process("auth_"+service, &innerConfig)
		if err != nil {
			return err
		}
		config.AuthConfigs[service] = innerConfig
	}
	auth = config

	return nil
}

// AuthConfigInstance attempts to retrieve a particular auth configuration set from the environment.
// Note this looks for environment variables prefixed with AUTH_${SERVICE_NAME}, and will only
// retrieve these values for services named in the AUTH_SERVICES environment
func AuthConfigInstance(name string) AuthInstanceConfig {
	if name == "ashirt" { // special case -- local auth doesn't have any normal environment variables
		return AuthInstanceConfig{
			Name:                "ashirt",
			Type:                "local",
			RegistrationEnabled: auth.AuthConfigs[name].RegistrationEnabled,
		}
	}
	v := auth.AuthConfigs[name]
	return v
}

// SupportedAuthServices retrieves the parsed AUTH_SERVICES value from the environment
func SupportedAuthServices() []string {
	return auth.Services
}
