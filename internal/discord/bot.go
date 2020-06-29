package discord

import (
	"github.com/Floor-Gang/MentionFilter/internal"
	"github.com/Floor-Gang/MentionFilter/internal/db"
	dg "github.com/bwmarrin/discordgo"
)

type Bot struct {
	version    string
	session    *dg.Session
	config     internal.Config
	controller *db.Controller
}

func Start(configPath string, dbName string) error {
	var err error
	botConfig := internal.GetConfig(configPath)

	controller := db.GetController(dbName)

	client, err := dg.New("Bot " + botConfig.Token)

	if err != nil {
		panic(err)
	}

	bot := Bot{
		session:    client,
		controller: controller,
		config:     botConfig,
	}

	client.AddHandler(bot.onReady)
	client.AddHandler(bot.onMessage)

	if err = client.Open(); err != nil {
		panic(err)
	}

	return err
}
