package db

import (
	"database/sql"
	"os"

	"github.com/Floor-Gang/MentionFilter/internal"
)

func initDB(dbName string) {
	_, err := os.Create(dbName)

	if err != nil {
		internal.Report(err)
		return
	}
}

// GetController acquires controller
func GetController(dbName string) *Controller {
	if _, err := os.Stat(dbName); err != nil {
		initDB(dbName)
	}

	db, err := sql.Open("sqlite3", dbName)

	if err != nil {
		internal.Report(err)
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
		(mention_id INTEGER PRIMARY KEY AUTOINCREMENT, 
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

// RemoveMention allows user to remove a filter from the database
func (c Controller) RemoveMention(id string) error {
	_, err := c.db.Exec(
		"DELETE FROM mentions WHERE mention_id=?",
		id,
	)

	return err
}

// AddMention allows user to add a filter to the database
func (c Controller) AddMention(req AddMention) error {
	statement, err := c.db.Prepare(
		"INSERT INTO mentions (regex, action, description) VALUES (?,?,?)",
	)

	if err != nil {
		internal.Report(err)
		return nil
	}

	_, err = statement.Exec(
		req.Regex,
		req.Action,
		req.Description,
	)

	return err
}

// GetMention retrieves a single mention based on ID
func (c Controller) GetMention(mentionID string) (Mention, error) {
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

// UpdateAction allows user to update action of a mention
func (c Controller) UpdateAction(req PartialActionMention) error {
	tx, _ := c.db.Begin()
	statement, err := tx.Prepare(
		`UPDATE mentions
		 SET action=?
		 WHERE mention_id=?`,
	)

	if err != nil {
		internal.Report(err)
		return err
	}

	_, err = statement.Exec(
		req.Action,
		req.MentionID,
	)

	if err != nil {
		internal.Report(err)
		return err
	}

	err = tx.Commit()

	return err
}

// UpdateRegex allows user to update regex of a mention
func (c Controller) UpdateRegex(req PartialRegexMention) error {
	tx, _ := c.db.Begin()
	statement, err := tx.Prepare(
		`UPDATE mentions
		 SET regex=?
		 WHERE mention_id=?`,
	)

	if err != nil {
		internal.Report(err)
		return err
	}

	_, err = statement.Exec(
		req.Regex,
		req.MentionID,
	)

	err = tx.Commit()

	return err
}

// UpdateDescription allows user to update description of a mention
func (c Controller) UpdateDescription(req PartialDescriptionMention) error {
	tx, _ := c.db.Begin()
	statement, err := tx.Prepare(
		`UPDATE mentions
		 SET description=?
		 WHERE mention_id=?`,
	)

	if err != nil {
		internal.Report(err)
		return err
	}

	_, err = statement.Exec(
		req.Description,
		req.MentionID,
	)

	if err != nil {
		internal.Report(err)
		return err
	}

	err = tx.Commit()

	return err
}

// HasMentionID checks if a filter in the database exists
func (c Controller) HasMentionID(mentionID string) (bool, error) {
	result, err := c.db.Query(`SELECT description 
							   FROM mentions 
							   WHERE mention_id=?`, mentionID)
	description := ""

	if err != nil {
		return false, err
	}

	for result.Next() {
		if err := result.Scan(&description); err != nil {
			return false, err
		}
	}

	return len(description) > 0, nil
}

// GetAllMentions acquires all mentions in the database
func (c Controller) GetAllMentions() (*sql.Rows, error) {
	rows, err := c.db.Query("SELECT * FROM mentions")

	if err != nil {
		return rows, err
	}

	return rows, nil
}
