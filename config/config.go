package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const defaultConfigFilePath = "./config.yml"

type Config struct {
	Plans struct {
		Directory          string `yaml:"directory"`
		ExerciseDefinition string `yaml:"exercise-definition"`
	} `yaml:"plans"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Timeout  struct {
			Request int `yaml:"request"`
		} `yaml:"timeout"`
	} `yaml:"database"`
	Server struct {
		Ip      string `yaml:"ip"`
		Port    string `yaml:"port"`
		Timeout struct {
			Write int `yaml:"write"`
			Read  int `yaml:"read"`
			Idle  int `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
}

func LoadConfig(configFile ...string) (*Config, error) {
	var configFilePath string
	if len(configFile) != 0 {
		configFilePath = configFile[0]
	} else {
		configFilePath = defaultConfigFilePath
	}
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
