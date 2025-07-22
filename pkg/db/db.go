package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	comment TEXT,
	title VARCHAR(256),
	repeat VARCHAR(128)
	);
	
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
	`

var db *sql.DB

func Init(dbFile string) error {
	dbFile = "scheduler.db"
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil

}
