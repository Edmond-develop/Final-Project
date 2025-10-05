package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(256) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX scheduler_date_ind ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if install {
		_, err := DB.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
