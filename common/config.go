package common

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DbUser, DbPass, DbAddr, DbName string
}

type ConfigFile struct {
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

	var configFile ConfigFile

	err = yaml.Unmarshal(yamlFile, &configFile)
	if err != nil {
		panic(err)
	}

	var config Config

	config.DbUser = configFile.DbUser
	config.DbPass = configFile.DbPass
	config.DbAddr = configFile.DbAddr
	config.DbName = configFile.DbName

	return &config
}
