package main

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

// Add mention to the db ".mention add <mentionID> <regex> <action> <description>"
func add(s *dg.Session, event *dg.MessageCreate, mentionID string, regex string, action string, description string) {
	member := event.Member
	if !checkChannel(s, event.Message) {
		reply(s, event, "I only work in my designated channel.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}

	if !checkAction(action) {
		reply(s, event, "Variable `<action>` can only be 'filter' or 'remove'")
		return
	}

	mentionIDExists, err := controller.hasMentionID(mentionID)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
		return
	}
	if mentionIDExists {
		reply(s, event, "A mention with this ID already exists.")
		return
	}

	req := Mention{
		MentionID:   mentionID,
		Regex:       regex,
		Action:      action,
		Description: description,
	}

	err = controller.addMention(req)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
	} else {
		reply(s, event, "Mention added.")
	}
}

// Change what happens on mention ".mention change_action <id> <type>"
func changeAction(s *dg.Session, event *dg.MessageCreate, mentionID string, action string) {
	member := event.Member
	if !checkChannel(s, event.Message) {
		reply(s, event, "I only work in my designated channel.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}

	if !checkAction(action) {
		reply(s, event, "Variable `<action>` can only be 'filter' or 'remove'")
		return
	}

	mentionIDExists, err := controller.hasMentionID(mentionID)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		reply(s, event, "A mention with this ID does not yet exist.")
		return
	}

	req := PartialActionMention{
		MentionID: mentionID,
		Action:    action,
	}

	err = controller.updateAction(req)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
	} else {
		reply(s, event, "Mention action changed.")
	}
}

// Change regex of mention ".mention change_regex <id> <regex>"
func changeRegex(s *dg.Session, event *dg.MessageCreate, mentionID string, regex string) {
	member := event.Member
	if !checkChannel(s, event.Message) {
		reply(s, event, "I only work in my designated channel.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}

	mentionIDExists, err := controller.hasMentionID(mentionID)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		reply(s, event, "A mention with this ID does not yet exist.")
		return
	}

	req := PartialRegexMention{
		MentionID: mentionID,
		Regex:     regex,
	}

	err = controller.updateRegex(req)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
	} else {
		reply(s, event, "Mention regex changed.")
	}
}

// Change description of mention ".mention change_description <id> <description>"
func changeDescription(s *dg.Session, event *dg.MessageCreate, mentionID string, description string) {
	member := event.Member
	if !checkChannel(s, event.Message) {
		reply(s, event, "I only work in my designated channel.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}

	mentionIDExists, err := controller.hasMentionID(mentionID)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		reply(s, event, "A mention with this ID does not yet exist.")
		return
	}

	req := PartialDescriptionMention{
		MentionID:   mentionID,
		Description: description,
	}

	err = controller.updateDescription(req)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
	} else {
		reply(s, event, "Mention description changed.")
	}
}

// Remove mention from db ".mention remove <id>"
func removeMention(s *dg.Session, event *dg.MessageCreate, mentionID string) {
	member := event.Member
	if !checkChannel(s, event.Message) {
		reply(s, event, "I only work in my designated channel.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}

	mentionIDExists, err := controller.hasMentionID(mentionID)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
		return
	}
	if !mentionIDExists {
		reply(s, event, "A mention with this ID does not yet exist.")
		return
	}

	err = controller.removeMention(mentionID)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
	} else {
		reply(s, event, "Mention removed")
	}
}

// View all mentions ".mention mentions"
func mentions(s *dg.Session, event *dg.MessageCreate) {
	member := event.Member
	if !checkChannel(s, event.Message) {
		reply(s, event, "I only work in my designated channel.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}

	rows, err := controller.getAllMentions()
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
	} else {
		var id string
		var regex string
		var action string
		var description string
		var mentionsSlice []Mention

		for rows.Next() {
			rows.Scan(&id, &regex, &action, &description)
			mention := Mention{
				MentionID:   id,
				Regex:       regex,
				Action:      action,
				Description: description,
			}
			mentionsSlice = append(mentionsSlice, mention)
		}

		_, err := AllMentionsEmbed(s, botConfig.ChannelID, mentionsSlice, "All mentions")
		if err != nil {
			report(err)
			reply(s, event, "Something went wrong.")
		}
	}
}

// View one mention ".mention <id>"
func mention(s *dg.Session, event *dg.MessageCreate, mentionID string) {
	member := event.Member
	if !checkChannel(s, event.Message) {
		reply(s, event, "I only work in my designated channel.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}

	result, err := controller.getMention(mentionID)
	if err != nil {
		report(err)
		reply(s, event, "Something went wrong.")
	} else {
		var mentionsSlice []Mention

		mentionsSlice = append(mentionsSlice, result)

		_, err := AllMentionsEmbed(s, botConfig.ChannelID, mentionsSlice, fmt.Sprintf("Mention %s", mentionID))
		if err != nil {
			report(err)
			reply(s, event, "Something went wrong.")
		}
	}
}

// Possibly for later:
// - Set mention in db to active
// - Set mention in db to inactive
