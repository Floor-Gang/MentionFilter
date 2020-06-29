package discord

import (
	"fmt"
	"github.com/Floor-Gang/MentionFilter/internal"
	"github.com/Floor-Gang/MentionFilter/internal/db"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strings"
)

var counter = 1
var allFilters []db.FilterableMention

func (b *Bot) onMessage(_ *dg.Session, event *dg.MessageCreate) {
	msg := event.Message

	counter--

	if counter == 0 {
		// re initiate the regex's
		allFilters = b.initiateFilters(event)
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

	if !strings.HasPrefix(msg.Content, b.config.Prefix) {
		for _, Filter := range allFilters {
			re, err := regexp.Compile(Filter.Regex)

			if err != nil {
				internal.Report(err)
				return
			}

			result := re.MatchString(msg.Content)

			if result {
				if Filter.Action == "remove" {
					// TODO: Handle this error
					b.session.ChannelMessageDelete(event.ChannelID, event.Message.ID)
				}

				if Filter.Action == "filter" {
					_, err = b.newMentionEmbed(msg)

					if err != nil {
						internal.Report(err)
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
			_, err := b.session.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s add <mention ID> <regex> <action (filter | remove)> <description>`", b.config.Prefix),
			)

			if err != nil {
				internal.Report(err)
				return
			}
		} else {
			mentionID := args[2]
			regex := args[3]
			action := args[4]
			description := strings.Join(args[5:], " ")

			b.add(event, mentionID, regex, action, description)
			allFilters = b.initiateFilters(event)
			counter = 100
		}
	}

	if command == "change_action" {
		// args = [prefix, change_action, <mentionID> <action>]
		if len(args) < 4 {
			_, err := b.session.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s change_action <mention ID> <action (filter | remove)>`", b.config.Prefix),
			)

			if err != nil {
				internal.Report(err)
				return
			}
		} else {
			mentionID := args[2]
			action := args[3]

			b.changeAction(event, mentionID, action)
			allFilters = b.initiateFilters(event)
			counter = 100
		}
	}

	if command == "change_regex" {
		// args = [prefix, change_regex, <mentionID> <regex>]
		if len(args) < 4 {
			_, err := b.session.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s change_regex <mention ID> <regex>`", b.config.Prefix),
			)

			if err != nil {
				internal.Report(err)
				return
			}
		} else {
			mentionID := args[2]
			regex := args[3]

			b.changeRegex(event, mentionID, regex)
			allFilters = b.initiateFilters(event)
			counter = 100
		}
	}

	if command == "change_description" {
		// args = [prefix, change_description, <mentionID> <description>]
		if len(args) < 4 {
			_, err := b.session.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s change_description <mention ID> <description>`", b.config.Prefix),
			)

			if err != nil {
				internal.Report(err)
				return
			}
		} else {
			mentionID := args[2]
			description := strings.Join(args[3:], " ")

			b.changeDescription(event, mentionID, description)
			allFilters = b.initiateFilters(event)
			counter = 100
		}
	}

	if command == "remove" {
		// args = [prefix, remove, <mentionID>]
		if len(args) < 3 {
			_, err := b.session.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s remove <mention ID>`", b.config.Prefix),
			)

			if err != nil {
				internal.Report(err)
				return
			}
		} else {
			mentionID := args[2]

			b.removeMention(event, mentionID)
			allFilters = b.initiateFilters(event)
			counter = 100
		}
	}

	if command == "mentions" {
		// args = [prefix, mentions]
		b.mentions(event)
	}

	if command == "mention" {
		// args = [prefix, mention, <mentionID>]
		if len(args) < 3 {
			_, err := b.session.ChannelMessageSend(
				msg.ChannelID,
				fmt.Sprintf("`%s mention <mention ID>`", b.config.Prefix),
			)

			if err != nil {
				internal.Report(err)
				return
			}
		} else {
			mentionID := args[2]

			b.mention(event, mentionID)
		}
	}
}

func (b *Bot) onReady(s *dg.Session, _ *dg.Ready) {
	log.Printf("Ready as %s (version %s)\n", s.State.User.Username, b.version)
}
