package config

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"

	"github.com/DevopsArtFactory/act/pkg/constants"
	"github.com/DevopsArtFactory/act/pkg/schema"
	"github.com/DevopsArtFactory/act/pkg/tools"
)

// GetConfig parse a configuration file and return Config struct
func GetConfig() (*schema.Config, error) {
	if err := tools.ClearOsEnv(); err != nil {
		return nil, err
	}
	return selectConfigWithProfile()
}

// ReadConfigOnly parse a configuration file and return Config struct without cleaning OS environment
func ReadConfigOnly() (*schema.Config, error) {
	return selectConfigWithProfile()
}

// selectConfigWithProfile choose config with specific profile
func selectConfigWithProfile() (*schema.Config, error) {
	targetProfile := viper.GetString("profile")
	configs, err := getLocalConfig()
	if err != nil {
		return nil, err
	}

	for _, c := range configs {
		if c.Profile == targetProfile {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("profile does not exist: %s", targetProfile)
}

// getLocalConfig read local configuration file
func getLocalConfig() ([]schema.Config, error) {
	var config []schema.Config

	if !tools.FileExists(constants.BaseFilePath) {
		return config, fmt.Errorf("no configuration file exists in $HOME/%s", constants.BaseFilePath)
	}

	yamlFile, err := ioutil.ReadFile(constants.BaseFilePath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return config, err
	}

	return setDefault(config), nil
}

func setDefault(configs []schema.Config) []schema.Config {
	for i, c := range configs {
		if c.Duration == 0 {
			configs[i].Duration = 7200
		}
	}

	return configs
}

// GetInitConfig creates new Config struct for initialization
func GetInitConfig(name string) []schema.Config {
	config := schema.Config{
		Profile:  constants.DefaultProfile,
		Name:     name,
		Duration: 3600,
		AssumeRoles: map[string]string{
			"preprod": constants.EmptyString,
			"prod":    constants.EmptyString,
		},
		Databases: map[string][]string{
			"preprod": {constants.EmptyString},
			"prod":    {constants.EmptyString},
		},
	}

	return []schema.Config{config}
}

// GetCurrentAWSConfig parse an aws configuration with profile
func GetCurrentAWSConfig(profile string) (*schema.AWSConfig, error) {
	cfg, err := ReadAWSConfig()
	if err != nil {
		return nil, err
	}

	section, err := cfg.GetSection(profile)
	if err != nil {
		return nil, err
	}

	return &schema.AWSConfig{
		AccessKeyID:     section.Key("aws_access_key_id").String(),
		SecretAccessKey: section.Key("aws_secret_access_key").String(),
	}, nil
}

// ReadAWSConfig parse an aws configuration
func ReadAWSConfig() (*ini.File, error) {
	if !tools.FileExists(constants.AWSCredentialsPath) {
		return nil, fmt.Errorf("no aws configuration file exists in $HOME/%s", constants.AWSCredentialsPath)
	}

	cfg, err := ini.Load(constants.AWSCredentialsPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
