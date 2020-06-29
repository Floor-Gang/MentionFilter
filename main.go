package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Floor-Gang/MentionFilter/internal/discord"
)

const (
	configPath = "./config.yml"
	dbName     = "mentions.db"
)

func main() {
	err := discord.Start(configPath, dbName)

	if err != nil {
		panic(err)
	}

	keepalive()
}

func keepalive() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
