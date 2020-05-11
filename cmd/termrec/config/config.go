package config

import (
	"encoding/base64"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var configFilePath = filepath.Join(os.Getenv("HOME"), ".config", "ashirt", "term-recorder.yaml")

// TermRecorderConfig is the root configuration file for this application. Config data is loaded from
// 3 sources: Configuration file, then environment variables, then command line flags.
type TermRecorderConfig struct {
	ConfigVersion   int64  `yaml:"configVersion"`
	OutputDir       string `yaml:"outputDir"      split_words:"true"`
	OutputFileName  string `yaml:"-"              split_words:"true"`
	OperationID     int64  `yaml:"operationID"    split_words:"true"`
	RecordingShell  string `yaml:"recordingShell" split_words:"true"`
	APIURL          string `yaml:"apiURL"         split_words:"true" envconfig:"api_url"`
	AccessKey       string `yaml:"accessKey"      split_words:"true"`
	SecretKeyBase64 string `yaml:"secretKey"      split_words:"true" envconfig:"secret_key"`
	SecretKey       []byte `yaml:"-"`
	StartupScript   string `yaml:"startupScript"  split_words:"true"`
	// OperationSlug   int     `yaml:"operationSlug" split_words:"true"`

	issues []string
}

// Issues returns a list of encountered issues while parsing the config datasources
func (t *TermRecorderConfig) Issues() []string {
	return t.issues
}

func generateStaticConfig() error {
	dummyConfig := TermRecorderConfig{
		ConfigVersion:  1,
		RecordingShell: os.Getenv("SHELL"),
	}

	outFile, err := os.Create(configFilePath)
	if err != nil {
		return errors.Wrap(err, "Unable to create config file")
	}
	defer outFile.Close()

	encoder := yaml.NewEncoder(outFile)
	err = encoder.Encode(dummyConfig)
	if err != nil {
		return errors.Wrap(err, "Unable to write config file")
	}
	return errors.Wrap(outFile.Close(), "Could not close config file")
}

// Parse actually reads the various data sources.
func Parse() TermRecorderConfig {
	config := TermRecorderConfig{
		// defaults go here
		issues:         make([]string, 0, 4),
		RecordingShell: os.Getenv("SHELL"),
	}
	f, err := os.Open(configFilePath)
	defer f.Close()
	if err != nil {
		if os.IsNotExist(err) { //is of type files does not exist {
			if err = generateStaticConfig(); err != nil {
				config.issues = append(config.issues, "Unable to generate new config: "+err.Error())
			} else {
				config.issues = append(config.issues, "Unable to find confile file. One has been created at: "+configFilePath)
			}
		} else {
			config.issues = append(config.issues, "Unable to open up default config file path: "+err.Error())
		}
	} else {
		config.parseFile(f)
	}
	config.parseEnv()
	config.parseCli()

	config.SecretKey, err = base64.StdEncoding.DecodeString(config.SecretKeyBase64)
	if err != nil {
		config.issues = append(config.issues, "Unable to decode SecretKey: "+err.Error())
	}

	return config
}

func (t *TermRecorderConfig) parseEnv() {
	err := envconfig.Process("ASHIRT_TERM_RECORDER", t)
	if err != nil {
		t.issues = append(t.issues, "Error reading env config: "+err.Error())
	}
}

func (t *TermRecorderConfig) parseFile(reader io.Reader) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.issues = append(t.issues, "Unable to read config file from datasource: "+err.Error())
		return
	}
	err = yaml.Unmarshal(bytes, &t)
	if err != nil {
		t.issues = append(t.issues, "Unable to interpret datasource as a yaml file: "+err.Error())
	}
}

func (t *TermRecorderConfig) parseCli() {
	attachStringFlag("output-dir", "", "Location where output file will be stored", t.OutputDir, &t.OutputDir)
	attachStringFlag("output-file", "o", "Name of the output file", t.OutputFileName, &t.OutputFileName)

	attachIntFlag("operation", "", "Operation ID to use when saving evidence", t.OperationID, &t.OperationID)

	attachStringFlag("shell", "s", "Name of the output file", t.RecordingShell, &t.RecordingShell)
	attachStringFlag("svc", "", "Name of the output file", t.APIURL, &t.APIURL)

	attachStringFlag("startup-script", "ss", "Script to run after beginning a recording", t.StartupScript, &t.StartupScript)
	flag.Parse()
}

// the below are small helpers to provide both short and long form flags -- not ideal, as it messes up
// the -h flag.

func attachStringFlag(longName, shortName, description, defaultValue string, variable *string) {
	if shortName != "" {
		flag.StringVar(variable, shortName, defaultValue, description)
	}
	flag.StringVar(variable, longName, defaultValue, description)
}

func attachIntFlag(longName, shortName, description string, defaultValue int64, variable *int64) {
	if shortName != "" {
		flag.Int64Var(variable, shortName, defaultValue, description)
	}
	flag.Int64Var(variable, longName, defaultValue, description)
}

func attachBoolFlag(longName, shortName, description string, defaultValue bool, variable *bool) {
	if shortName != "" {
		flag.BoolVar(variable, shortName, defaultValue, description)
	}
	flag.BoolVar(variable, longName, defaultValue, description)
}

func attachFloatFlag(longName, shortName, description string, defaultValue float64, variable *float64) {
	if shortName != "" {
		flag.Float64Var(variable, shortName, defaultValue, description)
	}
	flag.Float64Var(variable, longName, defaultValue, description)
}
