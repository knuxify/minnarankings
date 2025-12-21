package common

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DbUser string `yaml:"db_user"`
	DbPass string `yaml:"db_pass"`
	DbAddr string `yaml:"db_addr"`
	DbName string `yaml:"db_name"`
}

func ParseConfigFile(filename string) *Config {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return &config
}
