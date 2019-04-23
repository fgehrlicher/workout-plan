package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Plans struct {
		Directory string `yaml:"directory"`
	}
}

func LoadConfig(configFilePath string) (*Config, error) {
	config := &Config{}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}