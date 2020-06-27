package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

// Config structure
type Config struct {
	Token     string
	Prefix    string
	ChannelID string
	leadDevID string
	adminID   string
}

// Retrieve a configuration
func getConfig() Config {
	if _, err := os.Stat(configPath); err != nil {
		return genConfig()
	}
	file, err := ioutil.ReadFile(configPath)
	config := Config{}
	configMap := make(map[string]string)

	if err != nil {
		return genConfig()
	}

	err = yaml.Unmarshal(file, configMap)

	if err != nil {
		log.Println("Failed to read configuration file")
		panic(err)
	}

	config.Token = configMap["token"]
	config.Prefix = configMap["prefix"]
	config.ChannelID = configMap["channelid"]
	config.leadDevID = configMap["leaddevid"]
	config.adminID = configMap["adminid"]

	return config
}

// Generate a configuration
func genConfig() Config {
	config := Config{
		Token:     "",
		Prefix:    ".mention",
		ChannelID: "",
		leadDevID: "718816943452323880",
		adminID:   "718453895550074930",
	}

	_, err := os.Create(configPath)

	if err != nil {
		log.Println("Failed to create configuration file")
		panic(err)
	}

	_, err = os.Open(configPath)
	serialized, err := yaml.Marshal(config)
	err = ioutil.WriteFile(configPath, serialized, 600)

	if err != nil {
		log.Println("Failed to write to configuration file")
		panic(err)
	}

	return config
}

func (config Config) save() {
	file, err := os.Open(configPath)
	serialized, err := yaml.Marshal(config)
	_, err = file.Write(serialized)

	if err != nil {
		log.Println("Failed to save configuration file.")
		log.Fatalln(err)
	}
}
