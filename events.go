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

	// args = [prefix, command] // Splits on whitespace
	args := strings.Fields(msg.Content)

	if len(args) < 2 {
		return
	}

	command := args[1]

	if command == "add" {
		// args = [prefix, add, <mentionID> <regex> <action> <description>]
		if len(args) < 6 {
			_, err := s.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s add <mention ID> <regex> <action (filter | remove)> <description>`", botConfig.Prefix),
			)

			if err != nil {
				report(err)
				return
			}
		}

		mentionid := args[2]
		regex := args[3]
		action := args[4]
		description := strings.Join(args[5:], " ")

		add(s, event, mentionid, regex, action, description)
	}
}

func onReady(s *dg.Session, _ *dg.Ready) {
	log.Printf("Ready as %s (version %s)\n", s.State.User.Username, version)
}
