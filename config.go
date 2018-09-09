package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Config This is the top level configuration
type Config struct {
	Credentials AzureCredentials       `yaml:"credentials"`
	XXX         map[string]interface{} `yaml:",inline"`
}

// AzureCredentials This captures the credentials to access an Azure subscription
type AzureCredentials struct {
	Environment    string `yaml:"environment"`
	SubscriptionID string `yaml:"subscription_id"`
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
	TenantID       string `yaml:"tenant_id"`

	XXX map[string]interface{} `yaml:",inline"`
}

func (config *Config) loadConfig(configFile string) (err error) {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("Error reading config file: %s", err)
	}

	if err := yaml.Unmarshal(yamlFile, config); err != nil {
		return fmt.Errorf("Error parsing config file: %s", err)
	}

	return nil
}

func checkOverflow(m map[string]interface{}, ctx string) error {
	if len(m) > 0 {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		return fmt.Errorf("unknown fields in %s: %s", ctx, strings.Join(keys, ", "))
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (config *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config
	if err := unmarshal((*plain)(config)); err != nil {
		return err
	}
	if err := checkOverflow(config.XXX, "config"); err != nil {
		return err
	}
	return nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (credentials *AzureCredentials) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain AzureCredentials
	if err := unmarshal((*plain)(credentials)); err != nil {
		return err
	}
	if err := checkOverflow(credentials.XXX, "config"); err != nil {
		return err
	}
	return nil
}
