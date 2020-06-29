package internal

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
)

// Config structure.
type Config struct {
	Token     string `yaml:"token"`
	Prefix    string `yaml:"prefix"`
	ChannelID string `yaml:"channel"`
	LeadDevID string `yaml:"leadev"`
	AdminID   string `yaml:"admin"`
}

// Retrieve a configuration.
func GetConfig(configPath string) Config {
	if _, err := os.Stat(configPath); err != nil {
		genConfig(configPath)
		panic("Please populate the new config file.")
	}

	var file, err = ioutil.ReadFile(configPath)

	if err != nil {
		genConfig(configPath)
		panic("Please populate the new config file.")
	}

	config := Config{}
	err = yaml.Unmarshal(file, config)

	if err != nil {
		log.Println("Failed to read configuration file")
		panic(err)
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

	_, err := os.Create(configPath)

	if err != nil {
		log.Println("Failed to create configuration file")
		panic(err)
	}

	_, _ = os.Open(configPath)
	serialized, err := yaml.Marshal(config)

	if err != nil {
		log.Printf("Failed to serialize config\n%s\n", err)
		panic(err)
	}

	err = ioutil.WriteFile(configPath, serialized, 0660)

	if err != nil {
		log.Println("Failed to write to configuration file")
		panic(err)
	}
}
