package config

import (
	"gitlab.com/gps2.0/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

var appConfigInstance *AppConfig
var once sync.Once

// GetAppConfig returns an AppConfig instance
func GetAppConfig() *AppConfig {
	once.Do(func() {
		loadConfig()
	})

	return appConfigInstance
}

func loadConfig() {
	yamlFile, err := ioutil.ReadFile("config.yml")

	server.CheckError(err)

	err = yaml.Unmarshal(yamlFile, &appConfigInstance)
	server.CheckError(err)
}
