package data

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

	if err := setupAwsSettigsTable(); err != nil {
		log.Fatalln(err)
	}

	if err := createProfileView(); err != nil {
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
		Name TEXT NOT NULL UNIQUE,
		Description TEXT NOT NULL,
		Active BOOL NOT NULL DEFAULT FALSE,
		Created DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func setupAwsSettigsTable() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS AwsSettings (
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		ProfileId NOT NULL,
		Region TEXT NOT NULL,
		AccessKey TEXT NOT NULL,
		SecretAccessKey TEXT NOT NULL
	);`)
	return err
}

func createProfileView() error {
	_, _ = db.Exec("DROP VIEW IF EXISTS v_profiles;")
	_, err := db.Exec(`CREATE VIEW v_profiles AS
	SELECT 
		p.Id, p.Name, p.Description, ((ROW_NUMBER() OVER()) = 1) AS Active,
		aws.Region AS AwsRegion, aws.AccessKey AS AwsAccessKey, aws.SecretAccessKey AS AwsSecretAccessKey
	FROM (select * from profiles order by Active DESC, Id DESC) AS p 
	LEFT OUTER JOIN AwsSettings AS aws ON p.Id = aws.ProfileId
	`)
	return err
}
