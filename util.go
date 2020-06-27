package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	dg "github.com/bwmarrin/discordgo"
)

func keepalive() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func report(err error) {
	log.Printf("An error occurred %s\n", err)
}

func checkChannel(s *dg.Session) bool {
	channel, err := s.Channel(botConfig.ChannelID)

	if err != nil {
		report(err)
		return false
	}

	return channel != nil
}

func reply(s *dg.Session, event *dg.MessageCreate, message string) {
	_, err := s.ChannelMessageSend(event.ChannelID, fmt.Sprintf("<@%s> %s", event.Author.ID, message))
	if err != nil {
		report(err)
	}
}

// This makes an embed with the mention
func newMentionEmbed(s *dg.Session, channelID string, user *dg.User, mention string) (*dg.Message, error) {
	embed := dg.MessageEmbed{}
	name := fmt.Sprintf("%s#%s", user.Username, user.Discriminator)

	// Make the embed
	embed.Author = &dg.MessageEmbedAuthor{
		Name:    name,
		IconURL: user.AvatarURL(""),
	}
	embed.Color = 0xff0000
	embed.Description = fmt.Sprintf("")
	embed.Title = ""

	msg, err := s.ChannelMessageSendEmbed(channelID, &embed)

	if err != nil {
		report(err)
		return nil, err
	}

	return msg, nil
}
