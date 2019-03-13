// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package config

import (
	"encoding/base64"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const configFile = "conf.yaml"

// Config holds the configuration data
type Config struct {
	WorkingFolder   string `yaml:"workingFolder"`
	Campaign        string `yaml:"campaign"`
	CampaignID      uint64 `yaml:"campaignID"`
	APIURL          string `yaml:"apiURL"`
	AccessKey       string `yaml:"accessKey"`
	SecretKeyBase64 string `yaml:"secretKey"`
	SecretKey       []byte `yaml:"-"`
}

// Read reads the configuration file and loads it into a Config struct
func Read() (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	config.SecretKey, err = base64.StdEncoding.DecodeString(config.SecretKeyBase64)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// Save saves the config to the config file
func (c *Config) Save() error {
	marshaled, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, marshaled, 0600)
}
