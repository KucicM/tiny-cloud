package tinycloud

import (
	"database/sql"
)

func SetupProfileTabel() error {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Profiles (
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		Description TEXT NOT NULL
	);`)
	return err
}
