package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

const defaultConfigFilePath = "./config.yml"

type Config struct {
	Plans struct {
		Directory       string `yaml:"directory"`
		DefinitionsFile string `yaml:"definitions-file"`
	} `yaml:"plans"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Timeout  struct {
			Request int `yaml:"request"`
			Startup int `yaml:"startup"`
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
	Auth struct {
		Token struct {
			Issuer  string `yaml:"issuer"`
			Service string `yaml:"service"`
		}
	}
}

var configSingleton Config
var configOnce sync.Once

func GetConfig(configFile ...string) (*Config, error) {
	var err error
	configOnce.Do(func() {
		var configFilePath string
		if len(configFile) != 0 {
			configFilePath = configFile[0]
		} else {
			configFilePath = defaultConfigFilePath
		}
		configSingleton = Config{}
		err = configSingleton.loadFromFile(configFilePath)
	})

	return &configSingleton, err
}

func (config *Config) loadFromFile(configFilePath string) error {
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return errors.New(
			fmt.Sprintf(
				"can´t load config file: %v",
				err,
			),
		)
	}
	return yaml.Unmarshal(data, config)
}
