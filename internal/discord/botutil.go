package discord

import (
	"fmt"
	. "github.com/Floor-Gang/MentionFilter/internal"
	"github.com/Floor-Gang/MentionFilter/internal/db"
	dg "github.com/bwmarrin/discordgo"
	"time"
)

// Util struct methods

func (b *Bot) reply(event *dg.MessageCreate, message string) {
	_, err := b.session.ChannelMessageSend(event.ChannelID, fmt.Sprintf("<@%s> %s", event.Author.ID, message))
	if err != nil {
		Report(err)
	}
}

func (b *Bot) checkChannel(commandMessage *dg.Message) bool {
	return commandMessage.ChannelID == b.config.ChannelID
}

func (b *Bot) initiateFilters(event *dg.MessageCreate) []db.FilterableMention {
	var allFilters []db.FilterableMention

	rows, err := b.controller.GetAllMentions()
	if err != nil {
		Report(err)
		b.reply(event, "Something went wrong.")
	} else {
		var id string
		var regex string
		var action string
		var description string

		for rows.Next() {
			_ = rows.Scan(&id, &regex, &action, &description)
			filterableMention := db.FilterableMention{
				Regex:  regex,
				Action: action,
			}
			allFilters = append(allFilters, filterableMention)
		}
	}

	return allFilters
}

// NewMentionEmbed makes an embed with the mentionMessage
func (b *Bot) newMentionEmbed(mentionMessage *dg.Message) (*dg.Message, error) {
	messageURL := fmt.Sprintf("https://discordapp.com/channels/%s/%s/%s", mentionMessage.GuildID, mentionMessage.ChannelID, mentionMessage.ID)
	// TODO: Is there a better way to convert?
	timeStamp := fmt.Sprintf("%s", mentionMessage.Timestamp)

	guild, err := b.session.State.Guild(mentionMessage.GuildID)

	if err != nil {
		// TODO: Handle this error
		guild, err = b.session.Guild(mentionMessage.GuildID)
	}

	embed := dg.MessageEmbed{
		Author: &dg.MessageEmbedAuthor{
			Name:    mentionMessage.Author.Username,
			IconURL: mentionMessage.Author.AvatarURL(""),
		},
		Color: 0xff0000,
		Fields: []*dg.MessageEmbedField{
			{
				Name: "Server:",
				// guild is possible nil if error above (L61) isn't handled properly
				Value:  guild.Name,
				Inline: true,
			},
			{
				Name:   "Channel:",
				Value:  fmt.Sprintf("<#%s>", mentionMessage.ChannelID),
				Inline: true,
			},
			{
				Name:   "Author:",
				Value:  mentionMessage.Author.Mention(),
				Inline: true,
			},
			{
				Name:   "Time (UTC):",
				Value:  timeStamp,
				Inline: true,
			},
			{
				Name:   "Message Link:",
				Value:  messageURL,
				Inline: false,
			},
			{
				Name:   "Message:",
				Value:  mentionMessage.Content,
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     "New mention",
	}

	msg, err := b.session.ChannelMessageSendEmbed(b.config.ChannelID, &embed)

	if err != nil {
		Report(err)
		return nil, err
	}

	return msg, nil
}

func (b *Bot) checkRoles(member *dg.Member) bool {
	return StringInSlice(b.config.LeadDevID, member.Roles) ||
		StringInSlice(b.config.AdminID, member.Roles)
}

func checkAction(action string) bool {
	if action == "remove" || action == "filter" {
		return true
	}
	return false
}

// AllMentionsEmbed makes an embed with all mentions
func (b *Bot) allMentionsEmbed(mentionsSlice []db.Mention, title string) (*dg.Message, error) {
	var EmbedFields []*dg.MessageEmbedField
	for _, mention := range mentionsSlice {
		NewField := &dg.MessageEmbedField{
			Name:   fmt.Sprintf("Mention ID: %s", mention.MentionID),
			Value:  fmt.Sprintf("Regex: %s\nAction: %s\nDescription: %s", mention.Regex, mention.Action, mention.Description),
			Inline: false,
		}

		EmbedFields = append(EmbedFields, NewField)
	}

	embed := dg.MessageEmbed{
		Author:    &dg.MessageEmbedAuthor{},
		Color:     0xff0000,
		Fields:    EmbedFields,
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     title,
	}

	msg, err := b.session.ChannelMessageSendEmbed(b.config.ChannelID, &embed)

	if err != nil {
		Report(err)
		return nil, err
	}

	return msg, nil
}
