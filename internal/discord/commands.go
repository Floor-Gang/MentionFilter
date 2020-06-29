package discord

import (
	"fmt"

	"github.com/Floor-Gang/MentionFilter/internal"
	"github.com/Floor-Gang/MentionFilter/internal/db"
	dg "github.com/bwmarrin/discordgo"
)

// Add mention to the db ".mention add <regex> <action> <description>".
func (b *Bot) add(event *dg.MessageCreate, regex string, action string, description string) {
	member := event.Member
	if !b.checkChannel(event.Message) {
		b.reply(event, "I only work in my designated channel.")
		return
	}

	if !b.checkRoles(member) {
		b.reply(event, "You are not allowed to use this command.")
		return
	}

	if !checkAction(action) {
		b.reply(event, "Variable `<action>` can only be 'filter' or 'remove'")
		return
	}

	req := db.AddMention{
		Regex:       regex,
		Action:      action,
		Description: description,
	}

	err := b.controller.AddMention(req)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
	} else {
		b.reply(event, "Mention added.")
	}
}

// Change what happens on mention ".mention change_action <id> <type>".
func (b *Bot) changeAction(event *dg.MessageCreate, mentionID string, action string) {
	member := event.Member
	if !b.checkChannel(event.Message) {
		b.reply(event, "I only work in my designated channel.")
		return
	}

	if !b.checkRoles(member) {
		b.reply(event, "You are not allowed to use this command.")
		return
	}

	if !checkAction(action) {
		b.reply(event, "Variable `<action>` can only be 'filter' or 'remove'")
		return
	}

	mentionIDExists, err := b.controller.HasMentionID(mentionID)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		b.reply(event, "A mention with this ID does not yet exist.")
		return
	}

	req := db.PartialActionMention{
		MentionID: mentionID,
		Action:    action,
	}

	err = b.controller.UpdateAction(req)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
	} else {
		b.reply(event, "Mention action changed.")
	}
}

// Change regex of mention ".mention change_regex <id> <regex>".
func (b *Bot) changeRegex(event *dg.MessageCreate, mentionID string, regex string) {
	member := event.Member
	if !b.checkChannel(event.Message) {
		b.reply(event, "I only work in my designated channel.")
		return
	}

	if !b.checkRoles(member) {
		b.reply(event, "You are not allowed to use this command.")
		return
	}

	mentionIDExists, err := b.controller.HasMentionID(mentionID)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		b.reply(event, "A mention with this ID does not yet exist.")
		return
	}

	req := db.PartialRegexMention{
		MentionID: mentionID,
		Regex:     regex,
	}

	err = b.controller.UpdateRegex(req)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
	} else {
		b.reply(event, "Mention regex changed.")
	}
}

// Change description of mention ".mention change_description <id> <description>".
func (b *Bot) changeDescription(event *dg.MessageCreate, mentionID string, description string) {
	member := event.Member
	if !b.checkChannel(event.Message) {
		b.reply(event, "I only work in my designated channel.")
		return
	}

	if !b.checkRoles(member) {
		b.reply(event, "You are not allowed to use this command.")
		return
	}

	mentionIDExists, err := b.controller.HasMentionID(mentionID)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		b.reply(event, "A mention with this ID does not yet exist.")
		return
	}

	req := db.PartialDescriptionMention{
		MentionID:   mentionID,
		Description: description,
	}

	err = b.controller.UpdateDescription(req)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
	} else {
		b.reply(event, "Mention description changed.")
	}
}

// Remove mention from db ".mention remove <id>".
func (b *Bot) removeMention(event *dg.MessageCreate, mentionID string) {
	member := event.Member
	if !b.checkChannel(event.Message) {
		b.reply(event, "I only work in my designated channel.")
		return
	}

	if !b.checkRoles(member) {
		b.reply(event, "You are not allowed to use this command.")
		return
	}

	mentionIDExists, err := b.controller.HasMentionID(mentionID)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		b.reply(event, "A mention with this ID does not yet exist.")
		return
	}

	err = b.controller.RemoveMention(mentionID)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
	} else {
		b.reply(event, "Mention removed")
	}
}

// View all mentions ".mention mentions".
func (b *Bot) mentions(event *dg.MessageCreate) {
	member := event.Member
	if !b.checkChannel(event.Message) {
		b.reply(event, "I only work in my designated channel.")
		return
	}

	if !b.checkRoles(member) {
		b.reply(event, "You are not allowed to use this command.")
		return
	}

	rows, err := b.controller.GetAllMentions()
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
	} else {
		var id string
		var regex string
		var action string
		var description string
		var mentionsSlice []db.Mention

		for rows.Next() {
			rows.Scan(&id, &regex, &action, &description)
			mention := db.Mention{
				MentionID:   id,
				Regex:       regex,
				Action:      action,
				Description: description,
			}
			mentionsSlice = append(mentionsSlice, mention)
		}

		_, err := b.allMentionsEmbed(mentionsSlice, "All mentions")
		if err != nil {
			internal.Report(err)
			b.reply(event, "Something went wrong.")
		}
	}
}

// View one mention ".mention mention <id>".
func (b *Bot) mention(event *dg.MessageCreate, mentionID string) {
	member := event.Member
	if !b.checkChannel(event.Message) {
		b.reply(event, "I only work in my designated channel.")
		return
	}

	if !b.checkRoles(member) {
		b.reply(event, "You are not allowed to use this command.")
		return
	}

	mentionIDExists, err := b.controller.HasMentionID(mentionID)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		b.reply(event, "A mention with this ID does not yet exist.")
		return
	}

	result, err := b.controller.GetMention(mentionID)
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
	} else {
		var mentionsSlice []db.Mention

		mentionsSlice = append(mentionsSlice, result)

		_, err := b.allMentionsEmbed(mentionsSlice, fmt.Sprintf("Mention %s", mentionID))
		if err != nil {
			internal.Report(err)
			b.reply(event, "Something went wrong.")
		}
	}
}

// Help function
func (b *Bot) help(event *dg.MessageCreate) {
	_, err := b.helpEmbed()
	if err != nil {
		internal.Report(err)
		b.reply(event, "Something went wrong.")
	}
}

// Possibly for later:
// - Set mention in db to active
// - Set mention in db to inactive
