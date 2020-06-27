package main

import (
	dg "github.com/bwmarrin/discordgo"
)

// Add mention to the db ".mention add <mentionID> <regex> <action> <description>"
func add(s *dg.Session, event *dg.MessageCreate, mentionID string, regex string, action string, description string) {
	member := event.Member
	if member == nil {
		reply(s, event, "Please use this command in a guild.")
		return
	}

	if !checkChannel(s) {
		reply(s, event, "I can't find my designated channel. Let a developer know.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}
}

// Change what happens on mention ".mention change_action <id> <type>"
func changeAction(s *dg.Session, event *dg.MessageCreate, mentionID string, action string) {
	member := event.Member
	if member == nil {
		reply(s, event, "Please use this command in a guild.")
		return
	}

	if !checkChannel(s) {
		reply(s, event, "I can't find my designated channel. Let a developer know.")
		return
	}

	if !checkRoles(s, member) {
		reply(s, event, "You are not allowed to use this command.")
		return
	}
}

// Update mention in db ".mention update <id> <new regex>"

// Remove mention from db ".mention remove <id>"

// View all mentions ".mention mentions"

// View one mention ".mention <id>"

// Possibly for later:
// - Set mention in db to active
// - Set mention in db to inactive
