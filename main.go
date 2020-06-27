package main

import (
	dg "github.com/bwmarrin/discordgo"
)

const (
	version    = "1.0.0"
	configPath = "./config.yml"
	dbName     = "mentions.db"
)

var (
	botConfig  Config
	controller *Controller
)

func main() {
	var err error

	botConfig = getConfig()
	if err != nil {
		panic(err)
	}

	controller = getController()

	client, err := dg.New("Bot " + botConfig.Token)

	if err != nil {
		panic(err)
	}

	client.AddHandler(onReady)
	client.AddHandler(onMessage)

	if err = client.Open(); err != nil {
		panic(err)
	}

	keepalive()
}
