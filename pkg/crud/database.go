package crud

import (
	"database/sql"
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func SetupDatabes(url string) *sql.DB {
	if url == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("$HOME not defined")
		}

		p := path.Join(home, ".tiny-cloud")
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}

		url = path.Join(p, "tiny-cloud.db")
	}

	if err := openDatabase(url); err != nil {
		log.Fatalln(err)
	}

	if err := setupProfileTabel(); err != nil {
		log.Fatalln(err)
	}

	return db
}

func openDatabase(url string) error {

	var err error
	db, err = sql.Open("sqlite3", url)
	if err != nil {
		return err
	}

	return db.Ping()
}

func CloseDatabes() {
	if db != nil {
		db.Close()
	}
}

func setupProfileTabel() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Profiles (
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		Description TEXT NOT NULL
	);`)
	return err
}
