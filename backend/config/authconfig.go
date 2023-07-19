// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package config

import (
	"errors"
	"os"
	"strings"

	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/kelseyhightower/envconfig"
)

// AuthConfig provides configuration details for all authentication services
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
	OIDCConfig
	WebauthnConfig
}

type OIDCConfig struct {
	FriendlyName          string `split_words:"true"`
	ProviderURL           string `split_words:"true"`
	ClientID              string `split_words:"true"`
	ClientSecret          string `split_words:"true"`
	Scopes                string
	ProfileFirstNameField string `split_words:"true"`
	ProfileLastNameField  string `split_words:"true"`
	ProfileEmailField     string `split_words:"true"`
	ProfileSlugField      string `split_words:"true"`
}

type WebauthnConfig struct {
	DisplayName string `split_words:"true"`
	// All of the below have innate defaults, and so are effectively optional
	Timeout                         int
	RPOrigin                        string `envconfig:"RP_ORIGIN"`
	AttestationPreference           string `split_words:"true"`
	Debug                           bool
	AuthenticatorAttachment         string `split_words:"true"`
	AuthenticatorResidentKey        string `split_words:"true"`
	AuthenticatorRequireResidentKey *bool  `split_words:"true"`
	AuthenticatorUserVerification   string `split_words:"true"`
}

func (w WebauthnConfig) Conveyance() protocol.ConveyancePreference {
	val := strings.TrimSpace(strings.ToLower(w.AttestationPreference))
	if val == "indirect" {
		return protocol.PreferIndirectAttestation
	}
	if val == "direct" {
		return protocol.PreferDirectAttestation
	}
	// standard default
	return protocol.PreferNoAttestation
}

func (w WebauthnConfig) AuthenticatorAttachmentPreference() protocol.AuthenticatorAttachment {
	val := strings.TrimSpace((strings.ToLower((w.AuthenticatorAttachment))))
	if val == "platform" {
		return protocol.Platform
	}
	if val == "cross-platform" {
		return protocol.Platform
	}
	if val == "" {
		return "" // letting the library decide what the default is
	}
	// Warn here?
	return ""
}

func (w WebauthnConfig) AuthenticatorResidentKeyPreference() protocol.ResidentKeyRequirement {
	val := strings.TrimSpace(strings.ToLower(w.AuthenticatorResidentKey))
	if val == "preferred" {
		return protocol.ResidentKeyRequirementPreferred
	}
	if val == "required" {
		return protocol.ResidentKeyRequirementRequired
	}
	return protocol.ResidentKeyRequirementDiscouraged
}

func (w WebauthnConfig) AuthenticatorUserVerificationPreference() protocol.UserVerificationRequirement {
	val := strings.TrimSpace(strings.ToLower(w.AttestationPreference))
	if val == "required" {
		return protocol.VerificationRequired
	}
	if val == "discouraged" {
		return protocol.VerificationDiscouraged
	}
	return protocol.VerificationPreferred
}

func (w WebauthnConfig) BuildAuthenticatorSelection() protocol.AuthenticatorSelection {
	return protocol.AuthenticatorSelection{
		RequireResidentKey:      w.AuthenticatorRequireResidentKey,
		UserVerification:        w.AuthenticatorUserVerificationPreference(),
		ResidentKey:             w.AuthenticatorResidentKeyPreference(),
		AuthenticatorAttachment: w.AuthenticatorAttachmentPreference(),
	}
}

func splitNoSpaces(s, delimiter string) []string {
	arr := strings.Split(s, delimiter)
	for i, v := range arr {
		arr[i] = strings.TrimSpace(v)
	}
	return arr
}

func loadAuthConfig() error {
	config := AuthConfig{}
	servicesStr := os.Getenv("AUTH_SERVICES")
	if servicesStr == "" {
		return errors.New("Auth services not defined")
	}

	servicesArr := splitNoSpaces(servicesStr, ",")

	serviceRegistrationStr := os.Getenv("AUTH_SERVICES_ALLOW_REGISTRATION")
	serviceRegistrationArr := splitNoSpaces(serviceRegistrationStr, ",")

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
