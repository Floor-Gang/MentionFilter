package main

import (
	dg "github.com/bwmarrin/discordgo"
)

// Add mention to the db ".mention add <regex>"
func add(s *dg.Session, event *dg.MessageCreate, regex string) {
	member := event.Member
	if member == nil {
		// reply(s, event, "Please use this command in a guild.")
		return
	}
}

// Change what happens on mention ".mention set_type <id> <type>"

// Update mention in db ".mention update <id> <new regex>"

// Remove mention from db ".mention remove <id>"

// View all mentions ".mention mentions"

// View one mention ".mention <id>"

// Possibly for later:
// - Set mention in db to active
// - Set mention in db to inactive
