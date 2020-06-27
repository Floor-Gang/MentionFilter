package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// This makes an embed with the mentionMessage
func newMentionEmbed(s *dg.Session, channelID string, user *dg.User, mentionMessage *dg.Message) (*dg.Message, error) {
	messageURL := fmt.Sprintf("https://discordapp.com/channels/%s/%s/%s", mentionMessage.GuildID, mentionMessage.ChannelID, mentionMessage.ID)
	timeStamp := fmt.Sprintf("%s", mentionMessage.Timestamp)

	guild, err := s.State.Guild(mentionMessage.GuildID)
	if err != nil {
		guild, err = s.Guild(mentionMessage.GuildID)
	}

	embed := dg.MessageEmbed{
		Author: &dg.MessageEmbedAuthor{
			Name:    mentionMessage.Author.Username,
			IconURL: mentionMessage.Author.AvatarURL(""),
		},
		Color: 0xff0000,
		Fields: []*dg.MessageEmbedField{
			&dg.MessageEmbedField{
				Name:   "Server:",
				Value:  guild.Name,
				Inline: true,
			},
			&dg.MessageEmbedField{
				Name:   "Channel:",
				Value:  fmt.Sprintf("<#%s>", mentionMessage.ChannelID),
				Inline: true,
			},
			&dg.MessageEmbedField{
				Name:   "Author:",
				Value:  mentionMessage.Author.Mention(),
				Inline: true,
			},
			&dg.MessageEmbedField{
				Name:   "Time (UTC):",
				Value:  timeStamp,
				Inline: true,
			},
			&dg.MessageEmbedField{
				Name:   "Message Link:",
				Value:  messageURL,
				Inline: false,
			},
			&dg.MessageEmbedField{
				Name:   "Message:",
				Value:  mentionMessage.Content,
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     "New mention",
	}

	msg, err := s.ChannelMessageSendEmbed(channelID, &embed)

	if err != nil {
		report(err)
		return nil, err
	}

	return msg, nil
}
