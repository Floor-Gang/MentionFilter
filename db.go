package main

import (
	sql "database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Controller structure
type Controller struct {
	db *sql.DB
}

// Mention structure
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

// PartialActionMention structure
type PartialActionMention struct {
	// ID of the mention
	MentionID string
	// Action
	Action string
}

// PartialRegexMention structure
type PartialRegexMention struct {
	// ID of the mention
	MentionID string
	// Regex
	Regex string
}

// PartialDescriptionMention structure
type PartialDescriptionMention struct {
	// ID of the mention
	MentionID string
	// Description
	Description string
}

func initDB() {
	_, err := os.Create(dbName)

	if err != nil {
		report(err)
		return
	}
}

func getController() *Controller {
	if _, err := os.Stat(dbName); err != nil {
		initDB()
	}

	db, err := sql.Open("sqlite3", dbName)

	if err != nil {
		report(err)
	}

	controller := Controller{
		db: db,
	}

	controller.init()

	return &controller
}

func (c Controller) init() {
	statement, err := c.db.Prepare(
		`CREATE TABLE IF NOT EXISTS mentions 
		(mention_id INT PRIMARY KEY NOT NULL, 
		 regex TEXT NOT NULL, 
		 action TEXT NOT NULL, 
		 description TEXT NOT NULL);`,
	)
	if err != nil {
		panic(err)
	} else {
		_, err = statement.Exec()
		if err != nil {
			panic(err)
		}
	}
}

func (c Controller) removeMention(id string) error {
	_, err := c.db.Exec(
		"DELETE FROM mentions WHERE mention_id=?",
		id,
	)

	return err
}

func (c Controller) addMention(req Mention) error {
	statement, err := c.db.Prepare(
		"INSERT INTO mentions (mention_id, regex, action, description) VALUES (?,?,?,?)",
	)

	if err != nil {
		report(err)
		return nil
	}

	_, err = statement.Exec(
		req.MentionID,
		req.Regex,
		req.Action,
		req.Description,
	)

	return err
}

func (c Controller) getMention(mentionID string) (Mention, error) {
	statement, err := c.db.Prepare(
		`SELECT * FROM mentions WHERE mention_id=?`,
	)

	if err != nil {
		return Mention{}, err
	}

	res, err := statement.Query(mentionID)

	if err != nil {
		return Mention{}, err
	}

	result := Mention{
		MentionID:   "",
		Regex:       "",
		Action:      "",
		Description: "",
	}

	for res.Next() {
		err = res.Scan(&result.MentionID, &result.Regex, &result.Action, &result.Description)
		if err != nil {
			return Mention{}, err
		}
	}

	return result, nil
}

func (c Controller) updateAction(req PartialActionMention) error {
	statement, err := c.db.Prepare(
		`UPDATE mentions
		 SET action=?
		 WHERE mention_id=?`,
	)

	if err != nil {
		report(err)
		return nil
	}

	_, err = statement.Exec(
		req.MentionID,
		req.Action,
	)

	return err
}

func (c Controller) updateRegex(req PartialRegexMention) error {
	statement, err := c.db.Prepare(
		`UPDATE mentions
		 SET regex=?
		 WHERE mention_id=?`,
	)

	if err != nil {
		report(err)
		return nil
	}

	_, err = statement.Exec(
		req.MentionID,
		req.Regex,
	)

	return err
}

func (c Controller) updateDescription(req PartialDescriptionMention) error {
	statement, err := c.db.Prepare(
		`UPDATE mentions
		 SET description=?
		 WHERE mention_id=?`,
	)

	if err != nil {
		report(err)
		return nil
	}

	_, err = statement.Exec(
		req.MentionID,
		req.Description,
	)

	return err
}

func (c Controller) hasMentionID(mentionID string) (bool, error) {
	result, err := c.db.Query(`SELECT description 
							   FROM mentions 
							   WHERE mention_id=?`, mentionID)
	description := ""

	if err != nil {
		panic(err)
	}

	for result.Next() {
		if err := result.Scan(&description); err != nil {
			return false, err
		}
	}

	return len(description) > 0, nil
}

func (c Controller) getAllMentions() (*sql.Rows, error) {
	rows, err := c.db.Query("SELECT * FROM mentions")

	if err != nil {
		return rows, err
	}

	return rows, nil
}
