package main

import (
	"fmt"
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
	dirPath, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	err = discord.Start(fmt.Sprintf("%s\\config.yml", dirPath), dbName)

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
