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

// PartialMention structure
type PartialMention struct {
	// ID of the mention
	MentionID string
	// Action
	Action string
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

func (c Controller) updateMention(req PartialMention) error {
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
