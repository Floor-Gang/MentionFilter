package internal

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

// Config structure.
type Config struct {
	Token     string `yaml:"token"`
	Prefix    string `yaml:"prefix"`
	ChannelID string `yaml:"channel"`
	LeadDevID string `yaml:"leadev"`
	AdminID   string `yaml:"admin"`
}

// GetConfig retrieves a configuration.
func GetConfig(configPath string) Config {
	if _, err := os.Stat(configPath); err != nil {
		genConfig(configPath)
		panic("Please populate the new config file.")
	}

	file, err := ioutil.ReadFile(configPath)

	if err != nil {
		genConfig(configPath)
		log.Fatalln("Failed to read configuration file. " + err.Error())
	}

	config := Config{}

	if err = yaml.Unmarshal(file, &config); err != nil {
		log.Fatalln("Failed to parse configuration file. " + err.Error())
	}

	return config
}

// Generate a configuration.
func genConfig(configPath string) {
	config := Config{
		Token:     "",
		Prefix:    ".mention",
		ChannelID: "",
		LeadDevID: "",
		AdminID:   "",
	}

	if _, err := os.Create(configPath); err != nil {
		log.Fatalln("Failed to create configuration file. " + err.Error())
	}

	serialized, err := yaml.Marshal(config)

	if err != nil {
		log.Fatalln("Failed to serialize config. " + err.Error())
	}

	if err = ioutil.WriteFile(configPath, serialized, 0660); err != nil {
		log.Fatalln("Failed to write to configuration file. " + err.Error())
	}
}
