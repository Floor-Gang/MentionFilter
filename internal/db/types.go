package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Controller structure.
type Controller struct {
	db *sql.DB
}

// Mention structure.
type Mention struct {
	// ID of the mention
	MentionID string
	// Regex pattern
	Regex string
	// Action
	Action string
	// Description
	Description string
}

// PartialActionMention structure.
type PartialActionMention struct {
	// ID of the mention
	MentionID string
	// Action
	Action string
}

// PartialRegexMention structure.
type PartialRegexMention struct {
	// ID of the mention
	MentionID string
	// Regex
	Regex string
}

// PartialDescriptionMention structure.
type PartialDescriptionMention struct {
	// ID of the mention
	MentionID string
	// Description
	Description string
}

// FilterableMention structure.
type FilterableMention struct {
	// Regex
	Regex string
	// Action
	Action string
}
