package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

var counter = 1
var allFilters []FilterableMention

func onMessage(s *dg.Session, event *dg.MessageCreate) {
	msg := event.Message

	counter--

	if counter == 0 {
		// reinitiate the regexes
		allFilters = initiateFilters(s, event)
		counter = 100
	}

	// Ignore messages that aren't in a guild
	if len(event.GuildID) == 0 {
		return
	}

	// Ignore messages from bots
	if event.Author.Bot {
		return
	}

	if !strings.HasPrefix(msg.Content, botConfig.Prefix) {
		for _, Filter := range allFilters {
			re, err := regexp.Compile(Filter.Regex)

			if err != nil {
				report(err)
				return
			}

			result := re.MatchString(msg.Content)

			if result {
				if Filter.Action == "remove" {
					s.ChannelMessageDelete(event.ChannelID, event.Message.ID)
				}

				if Filter.Action == "filter" {
					msg, err = newMentionEmbed(s, botConfig.ChannelID, event.Author, msg)

					if err != nil {
						report(err)
						return
					}
				}

				return
			}
		}
	}

	// Decided to leave all command checking at commands.go
	//  Mainly since this allows me to easily add custom permission handling
	//  Perhaps even add custom channel handling
	//  All at a later date

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
		} else {
			mentionid := args[2]
			regex := args[3]
			action := args[4]
			description := strings.Join(args[5:], " ")

			add(s, event, mentionid, regex, action, description)
			allFilters = initiateFilters(s, event)
			counter = 100
		}
	}

	if command == "change_action" {
		// args = [prefix, change_action, <mentionID> <action>]
		if len(args) < 4 {
			_, err := s.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s change_action <mention ID> <action (filter | remove)>`", botConfig.Prefix),
			)

			if err != nil {
				report(err)
				return
			}
		} else {
			mentionid := args[2]
			action := args[3]

			changeAction(s, event, mentionid, action)
			allFilters = initiateFilters(s, event)
			counter = 100
		}
	}

	if command == "change_regex" {
		// args = [prefix, change_regex, <mentionID> <regex>]
		if len(args) < 4 {
			_, err := s.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s change_regex <mention ID> <regex>`", botConfig.Prefix),
			)

			if err != nil {
				report(err)
				return
			}
		} else {
			mentionid := args[2]
			regex := args[3]

			changeRegex(s, event, mentionid, regex)
			allFilters = initiateFilters(s, event)
			counter = 100
		}
	}

	if command == "change_description" {
		// args = [prefix, change_description, <mentionID> <description>]
		if len(args) < 4 {
			_, err := s.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s change_description <mention ID> <description>`", botConfig.Prefix),
			)

			if err != nil {
				report(err)
				return
			}
		} else {
			mentionid := args[2]
			description := strings.Join(args[3:], " ")

			changeDescription(s, event, mentionid, description)
			allFilters = initiateFilters(s, event)
			counter = 100
		}
	}

	if command == "remove" {
		// args = [prefix, remove, <mentionID>]
		if len(args) < 3 {
			_, err := s.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s remove <mention ID>`", botConfig.Prefix),
			)

			if err != nil {
				report(err)
				return
			}
		} else {
			mentionid := args[2]

			removeMention(s, event, mentionid)
			allFilters = initiateFilters(s, event)
			counter = 100
		}
	}

	if command == "mentions" {
		// args = [prefix, mentions]
		mentions(s, event)
	}

	if command == "mention" {
		// args = [prefix, mention, <mentionID>]
		if len(args) < 3 {
			_, err := s.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s mention <mention ID>`", botConfig.Prefix),
			)

			if err != nil {
				report(err)
				return
			}
		} else {
			mentionid := args[2]

			mention(s, event, mentionid)
		}
	}
}

func onReady(s *dg.Session, _ *dg.Ready) {
	log.Printf("Ready as %s (version %s)\n", s.State.User.Username, version)
}
