package main

import (
	"fmt"
	"log"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func onMessage(s *dg.Session, event *dg.MessageCreate) {
	msg := event.Message

	// Ignore messages that don't start with the prefix & aren't in a guild
	if !strings.HasPrefix(msg.Content, botConfig.Prefix) || len(event.GuildID) == 0 {
		return
	}

	// Ignore messages from bots
	if event.Author.Bot {
		return
	}

	// args = [prefix, command]
	args := strings.Fields(msg.Content)

	if len(args) < 2 {
		return
	}

	command := args[1]

	if command == "add" {
		// args = [prefix, request, <name>]
		if len(args) < 3 {
			_, err := s.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s request <the nickname you want>`", botConfig.Prefix),
			)

			if err != nil {
				report(err)
				return
			}
		}

		regexString := strings.Join(args[2:], " ")

		add(s, event, regexString)
	}
}

func onReady(s *dg.Session, _ *dg.Ready) {
	log.Printf("Ready as %s (version %s)\n", s.State.User.Username, version)
}
